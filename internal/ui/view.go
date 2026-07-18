package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700")).
			Background(lipgloss.Color("#1A1A1A")).
			Padding(0, 1)

	menuItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true).
			Background(lipgloss.Color("#333333"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	typeColors = map[string]string{
		"normal":   "#A8A878",
		"fire":     "#F08030",
		"water":    "#6890F0",
		"electric": "#F8D030",
		"grass":    "#78C850",
		"ice":      "#98D8D8",
		"fighting": "#C03028",
		"poison":   "#A040A0",
		"ground":   "#E0C068",
		"flying":   "#A890F0",
		"psychic":  "#F85888",
		"bug":      "#A8B820",
		"rock":     "#B8A038",
		"ghost":    "#705898",
		"dragon":   "#7038F8",
		"dark":     "#705848",
		"steel":    "#B8B8D0",
		"fairy":    "#EE99AC",
	}
)

func (m Model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("POKEDEX CLI") + "\n\n")

	if m.loading {
		s.WriteString("Loading... \n")
		return s.String()
	}

	switch m.currentView {
	case menuView:
		s.WriteString(m.renderMenu())
	case listView:
		s.WriteString(m.renderList())
	case detailView:
		s.WriteString(m.renderDetail())
	case myPokedexView:
		s.WriteString(m.renderMyPokedex())
	}

	if m.message != "" {
		s.WriteString("\n" + m.message + "\n")
	}
	return s.String()
}

func (m Model) renderMenu() string {
	var s strings.Builder

	options := []string{"Browse Pokemon", "My Pokedex", "Exit"}

	for i, option := range options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			s.WriteString(cursor + " " + selectedStyle.Render(option) + "\n")
		} else {
			s.WriteString(cursor + " " + menuItemStyle.Render(option) + "\n")
		}
	}

	s.WriteString("\n" + helpStyle.Render("↑/↓: navigate • enter: select • q: quit"))
	return s.String()
}

func (m Model) renderList() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("Pokemon List (Page %d)\n\n", m.page+1))

	for i, name := range m.pokemonList {
		cursor := " "
		caught := ""
		if m.pokedex.Has(name) {
			caught = " ✓"
		}

		if m.cursor == i {
			cursor = ">"
			s.WriteString(cursor + " " + selectedStyle.Render(name+caught) + "\n")
		} else {
			s.WriteString(cursor + " " + menuItemStyle.Render(name+caught) + "\n")
		}
	}
	s.WriteString("\n" + helpStyle.Render("↑/↓: navigate • enter: view • n: next page • p: prev page • b: back"))
	return s.String()
}

func (m Model) renderDetail() string {
	if m.selectedPoke == nil {
		return "No Pokemon Selected"
	}

	p := m.selectedPoke
	var s strings.Builder

	s.WriteString(fmt.Sprintf("%s (#%d)\n\n", strings.ToUpper(p.Name), p.ID))

	s.WriteString("Types: ")
	for i, t := range p.Types {
		color := typeColors[t.Type.Name]
		typeStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true)

		s.WriteString(typeStyle.Render(t.Type.Name))
		if i < len(p.Types)-1 {
			s.WriteString(", ")
		}
	}
	s.WriteString("\n\n")

	s.WriteString(fmt.Sprintf("Height: %.1fm Weight: %.1fkg XP: %d\n\n",
		float64(p.Height)/10, float64(p.Weight)/10, p.BaseExperience))

	s.WriteString("Stats:\n")
	for _, stat := range p.Stats {
		bar := strings.Repeat("█", stat.BaseStat/10)
		s.WriteString(fmt.Sprintf("  %-15s %3d %s\n", stat.Stat.Name+":", stat.BaseStat, bar))
	}

	s.WriteString("\n")
	if m.pokedex.Has(p.Name) {
		s.WriteString("Already Caught\n")
	} else {
		s.WriteString(helpStyle.Render("Press 'c' to catch this pokemon! \n"))
	}
	s.WriteString("\n" + helpStyle.Render("b: back"))
	return s.String()
}

func (m Model) renderMyPokedex() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("🎒 My Pokedex (%d caught)\n\n", m.pokedex.Count()))

	if len(m.pokemonList) == 0 {
		s.WriteString("Your Pokedex is empty! Go catch SOME POKEMON!\n")
	} else {
		for i, name := range m.pokemonList {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				s.WriteString(cursor + " " + selectedStyle.Render(name) + "\n")
			} else {
				s.WriteString(cursor + " " + menuItemStyle.Render(name) + "\n")
			}
		}
	}
	s.WriteString("\n" + helpStyle.Render("↑/↓: navigate • enter: view • b: back"))
	return s.String()
}
