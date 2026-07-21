package pokedex

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type CaughtPokemon struct {
	Name     string    `json:"name"`
	CaughtAt time.Time `json:"caught_at"`
}

type Pokedex struct {
	Pokemon map[string]CaughtPokemon `json:"pokemon"`
}

func New() *Pokedex {
	return &Pokedex{
		Pokemon: make(map[string]CaughtPokemon),
	}
}

// catch a pokemon - always succeeds (game logic handles catch rate)
func (p *Pokedex) Catch(name string) bool {
	if _, exists := p.Pokemon[name]; exists {
		return false
	}

	p.Pokemon[name] = CaughtPokemon{
		Name:     name,
		CaughtAt: time.Now(),
	}
	return true
}

// if pokemon in pokedex
func (p *Pokedex) Has(name string) bool {
	_, exists := p.Pokemon[name]
	return exists
}

// returns number of count of caught pokemon
func (p *Pokedex) Count() int {
	return len(p.Pokemon)
}

func (p *Pokedex) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".pokedex")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	filepath := filepath.Join(configDir, "pokedex.json")
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

// loads the pokedex from a file
func Load() (*Pokedex, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return New(), nil
	}

	filepath := filepath.Join(homeDir, ".pokedex", "pokedex.json")
	data, err := os.ReadFile(filepath)
	if err != nil {
		return New(), nil
	}

	var p Pokedex
	if err := json.Unmarshal(data, &p); err != nil {
		return New(), nil
	}

	return &p, nil
}

// list all caught pokemons name
func (p *Pokedex) List() []string {
	names := make([]string, 0, len(p.Pokemon))
	for name := range p.Pokemon {
		names = append(names, name)
	}
	return names
}
