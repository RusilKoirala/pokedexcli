package battle

import (
	"fmt"
	"math/rand"

	"github.com/rusilkoirala/pokedexcli/internal/pokeapi"
)

type Move struct {
	Name     string
	Type     string
	Power    int
	Accuracy int
	PP       int
	MaxPP    int
}

var CommonMoves = map[string]*Move{
	"tackle":        NewMove("Tackle", "normal", 40, 100, 35),
	"scratch":       NewMove("Scratch", "normal", 40, 100, 35),
	"quick-attack":  NewMove("Quick Attack", "normal", 40, 100, 30),
	"thunder-shock": NewMove("Thunder Shock", "electric", 40, 100, 30),
	"thunderbolt":   NewMove("Thunderbolt", "electric", 90, 100, 15),
	"water-gun":     NewMove("Water Gun", "water", 40, 100, 25),
	"ember":         NewMove("Ember", "fire", 40, 100, 25),
	"vine-whip":     NewMove("Vine Whip", "grass", 45, 100, 25),
	"bite":          NewMove("Bite", "dark", 60, 100, 25),
	"iron-tail":     NewMove("Iron Tail", "steel", 100, 75, 15),
}

func GetDefaultMoves() []*Move {
	return []*Move{
		NewMove("Tackle", "normal", 40, 100, 35),
		NewMove("Quick Attack", "normal", 40, 100, 30),
		NewMove("Scratch", "normal", 40, 100, 35),
		NewMove("Pound", "normal", 40, 100, 35),
	}
}

func NewMove(name, moveType string, power, accuracy, maxPP int) *Move {
	return &Move{
		Name:     name,
		Type:     moveType,
		Power:    power,
		Accuracy: accuracy,
		PP:       maxPP,
		MaxPP:    maxPP,
	}
}

func (m *Move) CanUse() bool {
	return m.PP > 0
}

func (m *Move) Use() {
	if m.PP > 0 {
		m.PP--
	}
}

type Battle struct {
	PlayerPokemon *BattlePokemon
	WildPokemon   *BattlePokemon
	Turn          int
	IsPlayerTurn  bool
	BattleLog     []string
	IsOver        bool
	PlayerWon     bool
}

type BattlePokemon struct {
	Pokemon   *pokeapi.Pokemon
	CurrentHP int
	MaxHP     int
	Level     int
}

// newBattle creates a new battle instance
func NewBattle(playerPokemon, wildPokemon *pokeapi.Pokemon) *Battle {
	// calculate HP based on base stats
	playerHP := calculateHP(playerPokemon, 10) // player Pokemon at level 10 for now
	wildLevel := rand.Intn(6) + 5              // wild Pokemon level 5-10
	wildHP := calculateHP(wildPokemon, wildLevel)

	return &Battle{
		PlayerPokemon: &BattlePokemon{
			Pokemon:   playerPokemon,
			CurrentHP: playerHP,
			MaxHP:     playerHP,
			Level:     10,
		},
		WildPokemon: &BattlePokemon{
			Pokemon:   wildPokemon,
			CurrentHP: wildHP,
			MaxHP:     wildHP,
			Level:     wildLevel,
		},
		Turn:         1,
		IsPlayerTurn: true,
		BattleLog:    []string{},
		IsOver:       false,
		PlayerWon:    false,
	}
}

// calculateHP calculates HP based on base stat and level
func calculateHP(p *pokeapi.Pokemon, level int) int {
	// Find HP stat
	hpStat := 50 // default
	for _, stat := range p.Stats {
		if stat.Stat.Name == "hp" {
			hpStat = stat.BaseStat
			break
		}
	}

	// Simple formula: HP = base HP + (level * 2)
	return hpStat + (level * 2)
}

// PlayerAttack processes player's attack
func (b *Battle) PlayerAttack() string {
	if !b.IsPlayerTurn || b.IsOver {
		return ""
	}

	damage := b.calculateDamage(b.PlayerPokemon, b.WildPokemon)
	b.WildPokemon.CurrentHP -= damage

	message := fmt.Sprintf("%s attacks! Dealt %d damage!", b.PlayerPokemon.Pokemon.Name, damage)

	if b.WildPokemon.CurrentHP <= 0 {
		b.WildPokemon.CurrentHP = 0
		b.IsOver = true
		b.PlayerWon = true
		message = fmt.Sprintf("Wild %s fainted! You won!", b.WildPokemon.Pokemon.Name)
	}

	b.IsPlayerTurn = false
	b.Turn++
	return message
}

// EnemyAttack processes enemy's attack
func (b *Battle) EnemyAttack() string {
	if b.IsPlayerTurn || b.IsOver {
		return ""
	}

	damage := b.calculateDamage(b.WildPokemon, b.PlayerPokemon)
	b.PlayerPokemon.CurrentHP -= damage

	message := fmt.Sprintf("Wild %s attacks! Dealt %d damage!", b.WildPokemon.Pokemon.Name, damage)

	if b.PlayerPokemon.CurrentHP <= 0 {
		b.PlayerPokemon.CurrentHP = 0
		b.IsOver = true
		b.PlayerWon = false
		message = fmt.Sprintf("%s fainted! You lost!", b.PlayerPokemon.Pokemon.Name)
	}

	b.IsPlayerTurn = true
	b.Turn++
	return message
}

// calculateDamage calculates damage from attacker to defender
func (b *Battle) calculateDamage(attacker, defender *BattlePokemon) int {
	// Get attack and defense stats
	attack := b.getStat(attacker.Pokemon, "attack")
	defense := b.getStat(defender.Pokemon, "defense")

	// Simple damage formula
	baseDamage := (attack * 2) - (defense / 2)

	// Add randomness (85% - 100%)
	randomFactor := 0.85 + (rand.Float64() * 0.15)
	damage := int(float64(baseDamage) * randomFactor)

	// Minimum damage of 5
	if damage < 5 {
		damage = 5
	}

	return damage
}

// getStat gets a specific stat value
func (b *Battle) getStat(p *pokeapi.Pokemon, statName string) int {
	for _, stat := range p.Stats {
		if stat.Stat.Name == statName {
			return stat.BaseStat
		}
	}
	return 50 // default
}

// GetCatchRate returns catch rate based on HP remaining
func (b *Battle) GetCatchRate() float64 {
	if b.WildPokemon.CurrentHP == 0 {
		return 0.95 // Very high if fainted
	}

	hpPercent := float64(b.WildPokemon.CurrentHP) / float64(b.WildPokemon.MaxHP)

	// Lower HP = higher catch rate
	// 100% HP = 40% catch
	// 50% HP = 60% catch
	// 25% HP = 75% catch
	// 10% HP = 90% catch
	catchRate := 0.4 + (0.5 * (1 - hpPercent))

	return catchRate
}
