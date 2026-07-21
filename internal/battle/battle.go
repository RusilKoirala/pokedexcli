package battle

import (
	"fmt"
	"math/rand"

	"github.com/rusilkoirala/pokedexcli/internal/pokeapi"
)

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
	Moves     []*Move
}

func NewBattle(playerPokemon, wildPokemon *pokeapi.Pokemon, playerMoves, wildMoves []*Move) *Battle {
	playerHP := calculateHP(playerPokemon, 10)
	wildLevel := rand.Intn(6) + 5
	wildHP := calculateHP(wildPokemon, wildLevel)

	return &Battle{
		PlayerPokemon: &BattlePokemon{
			Pokemon:   playerPokemon,
			CurrentHP: playerHP,
			MaxHP:     playerHP,
			Level:     10,
			Moves:     playerMoves,
		},
		WildPokemon: &BattlePokemon{
			Pokemon:   wildPokemon,
			CurrentHP: wildHP,
			MaxHP:     wildHP,
			Level:     wildLevel,
			Moves:     wildMoves,
		},
		Turn:         1,
		IsPlayerTurn: true,
		BattleLog:    []string{},
		IsOver:       false,
		PlayerWon:    false,
	}
}

func calculateHP(p *pokeapi.Pokemon, level int) int {
	hpStat := 50
	for _, stat := range p.Stats {
		if stat.Stat.Name == "hp" {
			hpStat = stat.BaseStat
			break
		}
	}
	return hpStat + (level * 2)
}

func (b *Battle) PlayerAttack(moveIndex int) string {
	if !b.IsPlayerTurn || b.IsOver {
		return ""
	}

	if moveIndex >= len(b.PlayerPokemon.Moves) || moveIndex < 0 {
		return "Invalid move!"
	}

	move := b.PlayerPokemon.Moves[moveIndex]

	if !move.CanUse() {
		return fmt.Sprintf("%s has no PP left!", move.Name)
	}

	move.Use()

	// Get defender types
	defenderTypes := []string{}
	for _, t := range b.WildPokemon.Pokemon.Types {
		defenderTypes = append(defenderTypes, t.Type.Name)
	}

	damage, effectiveness := b.calculateDamageWithType(b.PlayerPokemon, b.WildPokemon, move, defenderTypes)
	b.WildPokemon.CurrentHP -= damage

	message := fmt.Sprintf("%s used %s! Dealt %d damage!", b.PlayerPokemon.Pokemon.Name, move.Name, damage)

	if effectiveness != 1.0 {
		message += "\n" + GetEffectivenessMessage(effectiveness)
	}

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

func (b *Battle) EnemyAttack() string {
	if b.IsPlayerTurn || b.IsOver {
		return ""
	}

	// Enemy picks random move with PP
	validMoves := []*Move{}
	for _, move := range b.WildPokemon.Moves {
		if move.CanUse() {
			validMoves = append(validMoves, move)
		}
	}

	if len(validMoves) == 0 {
		return "Wild Pokemon has no moves left!"
	}

	move := validMoves[rand.Intn(len(validMoves))]
	move.Use()

	// Get defender types
	defenderTypes := []string{}
	for _, t := range b.PlayerPokemon.Pokemon.Types {
		defenderTypes = append(defenderTypes, t.Type.Name)
	}

	damage, effectiveness := b.calculateDamageWithType(b.WildPokemon, b.PlayerPokemon, move, defenderTypes)
	b.PlayerPokemon.CurrentHP -= damage

	message := fmt.Sprintf("Wild %s used %s! Dealt %d damage!", b.WildPokemon.Pokemon.Name, move.Name, damage)

	if effectiveness != 1.0 {
		message += "\n" + GetEffectivenessMessage(effectiveness)
	}

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

func (b *Battle) calculateDamageWithType(attacker, defender *BattlePokemon, move *Move, defenderTypes []string) (int, float64) {
	attack := b.getStat(attacker.Pokemon, "attack")
	defense := b.getStat(defender.Pokemon, "defense")

	// Base damage from move power
	baseDamage := (attack*move.Power)/50 - (defense / 4)

	// Type effectiveness
	totalEffectiveness := 1.0
	for _, defType := range defenderTypes {
		totalEffectiveness *= GetTypeEffectiveness(move.Type, defType)
	}

	// aapply effectiveness
	damage := int(float64(baseDamage) * totalEffectiveness)

	// randomness (85%-100%)
	randomFactor := 0.85 + (rand.Float64() * 0.15)
	damage = int(float64(damage) * randomFactor)

	if damage < 1 {
		damage = 1
	}

	return damage, totalEffectiveness
}

func (b *Battle) getStat(p *pokeapi.Pokemon, statName string) int {
	for _, stat := range p.Stats {
		if stat.Stat.Name == statName {
			return stat.BaseStat
		}
	}
	return 50
}

func (b *Battle) GetCatchRate() float64 {
	if b.WildPokemon.CurrentHP == 0 {
		return 0.95
	}

	hpPercent := float64(b.WildPokemon.CurrentHP) / float64(b.WildPokemon.MaxHP)
	catchRate := 0.4 + (0.5 * (1 - hpPercent))

	return catchRate
}
