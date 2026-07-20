package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	starterBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFE66D")).
		Padding(1, 2).
		Width(40)

	selectedStarterStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1A1A2E")).
		Background(lipgloss.Color("#FFE66D")).
		Bold(true).
		Padding(0, 2)

	normalStarterStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F7F7F7")).
		Padding(0, 2)
)

type StarterChoice struct {
	Name        string
	Type        string
	Description string
	ID          int
}

var Starters = []StarterChoice{
	{
		Name:        "Charmander",
		Type:        "Fire",
		Description: "A fire-type with a flame on its tail",
		ID:          4,
	},
	{
		Name:        "Bulbasaur",
		Type:        "Grass",
		Description: "A grass-type with a bulb on its back",
		ID:          1,
	},
	{
		Name:        "Squirtle",
		Type:        "Water",
		Description: "A water-type turtle Pokémon",
		ID:          7,
	},
}

func RenderStarterSelection(cursor int, width, height int) string {
	var s strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFE66D")).
		Align(lipgloss.Center).
		Width(width).
		Render("🎓 PROFESSOR OAK'S LABORATORY")

	s.WriteString(title + "\n\n")

	// Dialogue
	dialogue := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F7F7F7")).
		Italic(true).
		Align(lipgloss.Center).
		Width(width).
		Render("Welcome, young trainer! Choose your first Pokémon to begin your journey.")

	s.WriteString(dialogue + "\n\n")

	// Starter choices
	for i, starter := range Starters {
		var line string
		
		starterText := starter.Name + " (" + starter.Type + ")"
		
		if i == cursor {
			line = selectedStarterStyle.Render("▸ " + starterText)
		} else {
			line = normalStarterStyle.Render("  " + starterText)
		}
		
		centered := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(width).
			Render(line)
		
		s.WriteString(centered + "\n")
		
		// Show description for selected
		if i == cursor {
			desc := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#8B9798")).
				Italic(true).
				Align(lipgloss.Center).
				Width(width).
				Render(starter.Description)
			s.WriteString(desc + "\n")
		}
		
		s.WriteString("\n")
	}

	// Controls
	controls := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B9798")).
		Italic(true).
		Align(lipgloss.Center).
		Width(width).
		Render("↑/↓: navigate  •  enter: select  •  q: quit")

	s.WriteString("\n" + controls)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, s.String())
}
