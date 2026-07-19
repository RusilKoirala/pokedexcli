package battle

import (
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
}

// create Battle
func NewBattle(playerPokemon, wildPokemon *pokeapi.Pokemon) *Battle {

	playerHP := calculateHP(playerPokemon, 10) // player Pokemon at level 10 :D
	wildHP := calculateHP(wildPokemon, rand.Intn(6)+5)

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
			Level:     rand.Intn(6) + 5,
		},
		Turn:         1,
		IsPlayerTurn: true,
		BattleLog:    []string{},
		IsOver:       false,
		PlayerWon:    false,
	}
}

// calculate HP based on base stat and stat
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

// player's attack
func (b *Battle) PlayerAttack() string {
	if !b.IsPlayerTurn || b.IsOver {
		return ""
	}
	damage := b.calculateDamage(b.PlayerPokemon, b.WildPokemon)
	b.WildPokemon.CurrentHP -= damage

	message := ""
	if b.WildPokemon.CurrentHP <= 0 {
		b.WildPokemon.CurrentHP = 0
		b.IsOver = true
		b.PlayerWon = true
		message = "Wild " + b.WildPokemon.Pokemon.Name + " fainted! You won!"
	} else {
		message = "Dealt " + string(rune(damage)) + " damage!"
	}

	b.IsPlayerTurn = false
	b.Turn++
	return message
}

// enemy attack
func (b *Battle) EnemyAttack() string {
	if b.IsPlayerTurn || b.IsOver {
		return ""
	}
	damage := b.calculateDamage(b.WildPokemon, b.PlayerPokemon)
	b.PlayerPokemon.CurrentHP -= damage

	message := ""

	if b.PlayerPokemon.CurrentHP <= 0 {
		b.PlayerPokemon.CurrentHP = 0
		b.IsOver = true
		b.PlayerWon = false
		message = b.PlayerPokemon.Pokemon.Name + " fainted! You lost!"
	} else {
		message = "Enemy dealt " + string(rune(damage)) + " damage!"
	}

	b.IsPlayerTurn = true
	b.Turn++
	return message
}

// calculate dmage from attacker to defender
func (b *Battle) calculateDamage(attacker, defender *BattlePokemon) int {
	attack := b.getStat(attacker.Pokemon, "attack")
	defense := b.getStat(defender.Pokemon, "defender")

	baseDamage := (attack * 2) - (defense / 2)

	randomFactor := 0.85 + (rand.Float64() * 0.15)
	damage := int(float64(baseDamage) * randomFactor)

	if damage < 5 {
		damage = 5
	}
	return damage
}

// get stat valuee
func (b *Battle) getStat(p *pokeapi.Pokemon, statName string) int {
	for _, stat := range p.Stats {
		if stat.Stat.Name == statName {
			return stat.BaseStat
		}
	}
	return 50
}

// get catch ratee by hp
func (b *Battle) GetCatchRate() float64 {
	if b.WildPokemon.CurrentHP == 0 {
		return 0.95
	}

	hpPercent := float64(b.WildPokemon.CurrentHP) / float64(b.WildPokemon.MaxHP)

	catchRate := 0.4 + (0.5 * (1 - hpPercent))

	return catchRate
}
