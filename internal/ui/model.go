package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rusilkoirala/pokedexcli/internal/pokeapi"
	"github.com/rusilkoirala/pokedexcli/internal/pokedex"
)

type view int

type Model struct {
	api          *pokeapi.Client
	pokedex      *pokedex.Pokedex
	currentView  view
	pokemonList  []string
	selectedPoke *pokeapi.Pokemon
	message      string
	loading      bool
	cursor       int
	page         int
}

const (
	menuView view = iota
	listView
	detailView
	myPokedexView
)

func NewModel() Model {
	dex, _ := pokedex.Load()
	return Model{
		api:         pokeapi.NewClient(),
		pokedex:     dex,
		currentView: menuView,
		cursor:      0,
		page:        0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
