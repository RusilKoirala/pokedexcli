package ui

import (
	"fmt"
	"image"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rusilkoirala/pokedexcli/internal/pokeapi"
)

type pokemonListMsg struct {
	list []string
}

type pokemonDetailMsg struct {
	pokemon *pokeapi.Pokemon
	sprite  image.Image
}

type errorMsg struct {
	err error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.pokedex.Save()
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.currentView == menuView && m.cursor < 2 {
				m.cursor++
			} else if m.currentView == listView && m.cursor < len(m.pokemonList)-1 {
				m.cursor++
			}

		case "enter":
			return m.handleEnter()

		case "b", "esc":
			return m.handleBack()

		case "c":
			if m.currentView == detailView && m.selectedPoke != nil {
				return m.handleCatch()
			}

		case "n":
			if m.currentView == listView {
				m.page++
				m.cursor = 0
				return m, m.loadPokemonList()
			}

		case "p":
			if m.currentView == listView && m.page > 0 {
				m.page--
				m.cursor = 0
				return m, m.loadPokemonList()
			}
		}

	case pokemonListMsg:
		m.pokemonList = msg.list
		m.loading = false

	case pokemonDetailMsg:
		m.selectedPoke = msg.pokemon
		m.spriteImage = msg.sprite
		m.loading = false

	case errorMsg:
		m.message = fmt.Sprintf("Error: %v", msg.err)
		m.loading = false
	}

	return m, nil
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case menuView:
		switch m.cursor {
		case 0:
			m.currentView = listView
			m.cursor = 0
			m.loading = true
			return m, m.loadPokemonList()
		case 1:
			m.currentView = myPokedexView
			m.cursor = 0
			m.pokemonList = m.pokedex.List()
		case 2:
			m.pokedex.Save()
			return m, tea.Quit
		}
	case listView:
		if len(m.pokemonList) > 0 {
			m.currentView = detailView
			m.loading = true
			return m, m.loadPokemonDetail(m.pokemonList[m.cursor])
		}
	case myPokedexView:
		if len(m.pokemonList) > 0 {
			m.currentView = detailView
			m.loading = true
			return m, m.loadPokemonDetail(m.pokemonList[m.cursor])
		}
	}
	return m, nil
}

func (m Model) handleBack() (tea.Model, tea.Cmd) {
	m.message = ""
	switch m.currentView {
	case listView, myPokedexView:
		m.currentView = menuView
		m.cursor = 0

	case detailView:
		m.currentView = listView
		m.cursor = 0
	}
	return m, nil
}

func (m Model) handleCatch() (tea.Model, tea.Cmd) {
	if m.pokedex.Has(m.selectedPoke.Name) {
		m.message = fmt.Sprint("%s is already in your Pokedex!!", m.selectedPoke.Name)
	} else if m.pokedex.Catch(m.selectedPoke.Name) {
		m.message = fmt.Sprintf("YAY you caught %s!", m.selectedPoke.Name)
		m.pokedex.Save()
	} else {
		m.message = fmt.Sprintf("%s escaped! Nice tryyy", m.selectedPoke.Name)
	}
	return m, nil
}

func (m Model) loadPokemonList() tea.Cmd {
	return func() tea.Msg {
		list, err := m.api.GetPokemonList(20, m.page*20)
		if err != nil {
			return errorMsg{err}
		}
		names := make([]string, len(list.Results))
		for i, p := range list.Results {
			names[i] = p.Name
		}
		return pokemonListMsg{names}
	}
}

func (m Model) loadPokemonDetail(name string) tea.Cmd {
	return func() tea.Msg {
		pokemon, err := m.api.GetPokemon(name)
		if err != nil {
			return errorMsg{err}
		}

		var sprite image.Image
		if pokemon.Sprites.FrontDefault != "" {
			sprite, _ = m.api.DownloadSprite(pokemon.Sprites.FrontDefault)
		}

		return pokemonDetailMsg{
			pokemon: pokemon,
			sprite:  sprite,
		}
	}
}
