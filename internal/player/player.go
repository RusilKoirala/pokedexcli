package player

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Player struct {
	Name          string `json:"name"`
	Level         int    `json:"level"`
	XP            int    `json:"xp"`
	XPToNextLevel int    `json:"xp_to_next_level"`
	StarterPokemon string `json:"starter_pokemon"`
	HasStarter    bool   `json:"has_starter"`
}

// NewPlayer creates a new player
func NewPlayer(name string) *Player {
	return &Player{
		Name:          name,
		Level:         1,
		XP:            0,
		XPToNextLevel: 100,
		StarterPokemon: "",
		HasStarter:    false,
	}
}

// SetStarter sets the player's starter Pokemon
func (p *Player) SetStarter(starter string) {
	p.StarterPokemon = starter
	p.HasStarter = true
}

// GainXP adds XP and handles level ups
func (p *Player) GainXP(amount int) bool {
	p.XP += amount
	
	// Check if leveled up
	if p.XP >= p.XPToNextLevel {
		return p.LevelUp()
	}
	
	return false
}

// LevelUp increases player level
func (p *Player) LevelUp() bool {
	if p.XP < p.XPToNextLevel {
		return false
	}
	
	p.XP -= p.XPToNextLevel
	p.Level++
	
	// Calculate next level XP requirement (increases by 20% each level)
	p.XPToNextLevel = int(float64(p.XPToNextLevel) * 1.2)
	
	return true
}

// GetXPBar returns a visual XP bar string
func (p *Player) GetXPBar(width int) string {
	percentage := float64(p.XP) / float64(p.XPToNextLevel)
	filled := int(percentage * float64(width))
	
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	
	return bar
}

// Save saves player data to file
func (p *Player) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	configDir := filepath.Join(homeDir, ".pokedexcli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	
	filePath := filepath.Join(configDir, "player.json")
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filePath, data, 0644)
}

// Load loads player data from file
func Load() (*Player, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	
	filePath := filepath.Join(homeDir, ".pokedexcli", "player.json")
	
	// If file doesn't exist, return new player
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return NewPlayer("Red"), nil
	}
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	
	var player Player
	if err := json.Unmarshal(data, &player); err != nil {
		return nil, err
	}
	
	return &player, nil
}

// GetXPForAction returns XP rewards for different actions
func GetXPForAction(action string) int {
	rewards := map[string]int{
		"catch":       50,
		"battle_win":  75,
		"battle_lose": 20,
		"explore":     10,
	}
	
	if xp, ok := rewards[action]; ok {
		return xp
	}
	
	return 0
}

// CanAccessArea checks if player level is high enough for an area
func (p *Player) CanAccessArea(requiredLevel int) bool {
	return p.Level >= requiredLevel
}

// GetAreaMessage returns a message about area access
func (p *Player) GetAreaMessage(areaName string, requiredLevel int) string {
	if p.CanAccessArea(requiredLevel) {
		return fmt.Sprintf("You can now explore %s!", areaName)
	}
	return fmt.Sprintf("🔒 %s requires Level %d (You are Level %d)", areaName, requiredLevel, p.Level)
}
