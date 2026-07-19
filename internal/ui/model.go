package ui

import (
	"image"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rusilkoirala/pokedexcli/internal/battle"
	"github.com/rusilkoirala/pokedexcli/internal/pokeapi"
	"github.com/rusilkoirala/pokedexcli/internal/pokedex"
)

type view int

const (
	startView view = iota
	menuView
	creditsView
	listView
	detailView
	myPokedexView
	exploreView
	encounterView
	pokemonSelectView
	battleView
)

type encounterState int

const (
	appearing encounterState = iota
	choosing
	throwing
	shaking
	caught
	escaped
)

type battleAction int

const (
	actionAttack battleAction = iota
	actionRun
	actionCatch
	actionItems
)

type Model struct {
	api          *pokeapi.Client
	pokedex      *pokedex.Pokedex
	currentView  view
	cursor       int
	pokemonList  []string
	selectedPoke *pokeapi.Pokemon
	spriteImage  image.Image
	message      string
	loading      bool
	page         int

	currentLocation  int
	encounterPokemon *pokeapi.Pokemon
	encounterSprite  image.Image
	encounterState   encounterState
	shakeCount       int
	totalEncounters  int

	currentBattle        *battle.Battle
	selectedMoveIndex    int
	selectedPokemonIndex int
	battleLog            string
	playerBattleSprite   image.Image
	enemyBattleSprite    image.Image
}

func NewModel() Model {
	dex, _ := pokedex.Load()
	return Model{
		api:               pokeapi.NewClient(),
		pokedex:           dex,
		currentView:       startView,
		cursor:            0,
		page:              0,
		currentLocation:   0,
		shakeCount:        0,
		totalEncounters:   0,
		selectedMoveIndex: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
