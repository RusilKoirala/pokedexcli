package ui

import (
	"fmt"
	"strings"
	"image"

	"github.com/charmbracelet/lipgloss"
	"github.com/qeesung/image2ascii/convert"
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

// convertImageToASCII converts an image to colored ASCII art
func convertImageToASCII(img image.Image) string {
	if img == nil {
		return "No sprite available"
	}
	
	// Create converter with options
	convertOptions := convert.DefaultOptions
	convertOptions.FixedWidth = 40
	convertOptions.FixedHeight = 20
	convertOptions.Colored = true
	convertOptions.Reversed = false
	
	converter := convert.NewImageConverter()
	return converter.Image2ASCIIString(img, &convertOptions)
}

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
	
	// Title
	title := fmt.Sprintf("%s (#%d)", strings.ToUpper(p.Name), p.ID)
	
	// Left column - Pokemon info
	var leftCol strings.Builder
	
	// Types
	leftCol.WriteString("Types: ")
	for i, t := range p.Types {
		color := typeColors[t.Type.Name]
		typeStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true)
		leftCol.WriteString(typeStyle.Render(t.Type.Name))
		if i < len(p.Types)-1 {
			leftCol.WriteString(", ")
		}
	}
	leftCol.WriteString("\n\n")

	// Physical info
	leftCol.WriteString(fmt.Sprintf("Height: %.1fm\n", float64(p.Height)/10))
	leftCol.WriteString(fmt.Sprintf("Weight: %.1fkg\n", float64(p.Weight)/10))
	leftCol.WriteString(fmt.Sprintf("Base XP: %d\n\n", p.BaseExperience))

	// Stats
	leftCol.WriteString("Stats:\n")
	for _, stat := range p.Stats {
		bar := strings.Repeat("█", stat.BaseStat/10)
		leftCol.WriteString(fmt.Sprintf("  %-12s %3d %s\n", stat.Stat.Name+":", stat.BaseStat, bar))
	}

	// Right column - ASCII sprite
	var rightCol string
	if m.spriteImage != nil {
		rightCol = convertImageToASCII(m.spriteImage)
	} else {
		rightCol = "\n\n   No sprite\n   available\n"
	}

	// Style the columns
	leftStyle := lipgloss.NewStyle().
		Width(50).
		Align(lipgloss.Left).
		PaddingRight(2)

	rightStyle := lipgloss.NewStyle().
		Width(40).
		Align(lipgloss.Center)

	// Join columns side by side
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(leftCol.String()),
		rightStyle.Render(rightCol),
	)

	// Build final output
	var output strings.Builder
	output.WriteString(title + "\n\n")
	output.WriteString(content + "\n\n")

	// Catch status
	if m.pokedex.Has(p.Name) {
		output.WriteString("✓ Already caught!\n")
	} else {
		output.WriteString(helpStyle.Render("Press 'c' to catch this Pokemon!\n"))
	}

	output.WriteString("\n" + helpStyle.Render("b: back"))
	return output.String()
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
