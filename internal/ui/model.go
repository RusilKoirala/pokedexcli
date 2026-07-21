package ui

import (
	"image"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rusilkoirala/pokedexcli/internal/battle"
	"github.com/rusilkoirala/pokedexcli/internal/dialogue"
	"github.com/rusilkoirala/pokedexcli/internal/npc"
	"github.com/rusilkoirala/pokedexcli/internal/player"
	"github.com/rusilkoirala/pokedexcli/internal/pokeapi"
	"github.com/rusilkoirala/pokedexcli/internal/pokedex"
	"github.com/rusilkoirala/pokedexcli/internal/quest"
	"github.com/rusilkoirala/pokedexcli/internal/town"
)

type view int

const (
	startView view = iota
	starterSelectionView
	menuView
	creditsView
	listView
	detailView
	myPokedexView
	exploreView
	encounterView
	overworldView
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

	currentMap     *town.WorldMap
	playerX        int
	playerY        int
	stepCount      int
	encounterSteps int

	npcManager      *npc.NPCManager
	currentDialogue *dialogue.DialogueBox
	dialogueActive  bool
	activeTrainerID string

	questManager *quest.QuestManager

	player *player.Player

	width  int
	height int
}

func NewModel() Model {
	dex, _ := pokedex.Load()
	playerData, _ := player.Load()

	npcMgr := npc.InitializeNPCs()
	questMgr, _ := quest.Load()

	// always start with the start screen cuz why not
	initialView := startView

	return Model{
		api:               pokeapi.NewClient(),
		pokedex:           dex,
		npcManager:        npcMgr,
		player:            playerData,
		currentView:       initialView,
		questManager:      questMgr,
		cursor:            0,
		page:              0,
		currentLocation:   0,
		shakeCount:        0,
		totalEncounters:   0,
		selectedMoveIndex: 0,
		width:             100,
		height:            30,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
