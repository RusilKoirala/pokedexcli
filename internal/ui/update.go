package ui

import (
	"fmt"
	"image"
	"math/rand"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rusilkoirala/pokedexcli/internal/locations"
	"github.com/rusilkoirala/pokedexcli/internal/pokeapi"
)

type pokemonListMsg struct {
	list []string
}

type pokemonDetailMsg struct {
	pokemon *pokeapi.Pokemon
	sprite  image.Image
}

type encounterMsg struct {
	pokemon *pokeapi.Pokemon
	sprite  image.Image
}

type errorMsg struct {
	err error
}

type tickMsg time.Time

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.currentView == encounterView {

				return m, nil
			}
			m.pokedex.Save()
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			switch m.currentView {
			case menuView:
				if m.cursor < 3 {
					m.cursor++
				}
			case listView:
				if m.cursor < len(m.pokemonList)-1 {
					m.cursor++
				}
			case exploreView:
				if m.cursor < locations.GetLocationCount()-1 {
					m.cursor++
				}
			}

		case "enter":
			return m.handleEnter()

		case "b", "esc":
			return m.handleBack()

		case "c":
			if m.currentView == detailView && m.selectedPoke != nil {
				return m.handleCatch()
			} else if m.currentView == encounterView && m.encounterState == choosing {
				return m.handleCatchWild()
			}

		case "e":
			if m.currentView == exploreView {
				return m.handleExplore()
			}

		case "r":
			if m.currentView == encounterView && m.encounterState == choosing {
				return m.handleRun()
			}

		case "i":
			if m.currentView == encounterView && m.encounterState == choosing {

				m.message = "Inspecting Pokemon..."
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

	case encounterMsg:
		m.encounterPokemon = msg.pokemon
		m.encounterSprite = msg.sprite
		m.encounterState = appearing
		m.loading = false
		m.totalEncounters++
		return m, tick()

	case tickMsg:
		return m.handleTick()

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
			m.currentView = exploreView
			m.cursor = 0
			m.message = ""
		case 3:
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

	case exploreView:
		m.currentLocation = m.cursor
	}
	// Stay in explore view
	return m, nil
}

func (m Model) handleBack() (tea.Model, tea.Cmd) {
	m.message = ""
	switch m.currentView {
	case listView, myPokedexView, exploreView:
		m.currentView = menuView
		m.cursor = 0

	case detailView:
		m.currentView = listView
		m.cursor = 0

	case encounterView:
		if m.encounterState == caught || m.encounterState == escaped {
			m.currentView = exploreView
			m.encounterState = appearing
		}
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

func tick() tea.Cmd {
	return tea.Tick(800*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) loadEncounter(pokemonID int) tea.Cmd {
	return func() tea.Msg {
		pokemon, err := m.api.GetPokemon(strconv.Itoa(pokemonID))
		if err != nil {
			return errorMsg{err}
		}

		var sprite image.Image
		if pokemon.Sprites.FrontDefault != "" {
			sprite, _ = m.api.DownloadSprite(pokemon.Sprites.FrontDefault)
		}

		return encounterMsg{
			pokemon: pokemon,
			sprite:  sprite,
		}
	}
}

func (m Model) calculateCatchRate() float64 {
	baseCatchRate := 0.4

	dexBonus := float64(m.pokedex.Count()) * 0.01

	total := baseCatchRate + dexBonus

	if total > 0.9 {
		total = 0.9
	}

	return total
}

func (m Model) handleRun() (tea.Model, tea.Cmd) {
	m.message = "You ran away safely!"
	m.encounterState = escaped
	return m, tick()
}

func (m Model) handleCatchWild() (tea.Model, tea.Cmd) {
	m.encounterState = throwing
	m.shakeCount = 0
	return m, tick()
}

func (m Model) handleTick() (tea.Model, tea.Cmd) {
	switch m.encounterState {
	case appearing:
		m.encounterState = choosing
		return m, nil

	case throwing:
		m.encounterState = shaking
		m.shakeCount = 0
		return m, tick()

	case shaking:
		m.shakeCount++
		if m.shakeCount < 3 {
			return m, tick()
		}

		catchRate := m.calculateCatchRate()
		if rand.Float64() < catchRate {
			m.encounterState = caught
			m.pokedex.Catch(m.encounterPokemon.Name)
			m.pokedex.Save()
			m.message = fmt.Sprintf("Gotcha! %s was caught!", m.encounterPokemon.Name)
		} else {
			m.encounterState = escaped
			m.message = fmt.Sprintf("%s broke free and escaped!", m.encounterPokemon.Name)
		}

		return m, nil

	case escaped:
		time.Sleep(2 * time.Second)
		m.currentView = exploreView
		m.encounterState = appearing
		return m, nil
	}
	return m, nil
}

func (m Model) handleExplore() (tea.Model, tea.Cmd) {
	location := locations.GetLocation((m.currentLocation))
	pokemonID := location.GetRandomPokemonID()

	m.currentView = encounterView
	m.loading = true
	m.message = ""

	return m, m.loadEncounter(pokemonID)
}
