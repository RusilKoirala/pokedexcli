package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rusilkoirala/pokedexcli/internal/player"
)

var (
	panelStyle      = lipgloss.NewStyle().BorderForeground(lipgloss.Color("#4ECDC4")).Padding(1, 2)
	panelTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFE66D")).Bold(true).Align(lipgloss.Center)
	panelTextStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F7F7F7"))
	xpBarStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#95E1D3"))
)

func RenderPlayerPanel(p *player.Player, width int) string {
	var s strings.Builder

	s.WriteString(panelTitleStyle.Render("TRAINER INFO"))
	s.WriteString("\n\n")

	s.WriteString(panelTextStyle.Render(fmt.Sprintf("Trainer: %s", p.Name)))
	s.WriteString(panelTextStyle.Render(fmt.Sprintf("Level: %d", p.Level)))
	s.WriteString("\n\n")

	xpText := fmt.Sprint("XP: %d/%d", p.XP, p.XPToNextLevel)
	s.WriteString(panelTextStyle.Render(xpText))
	s.WriteString(xpBarStyle.Render(p.GetXPBar(width - 6)))
	s.WriteString("\n\n")

	if p.HasStarter {
		s.WriteString(panelTitleStyle.Render("PARTNER"))
		s.WriteString("\n")
		s.WriteString(panelTextStyle.Render(fmt.Sprintf("⚡ %s", p.StarterPokemon)))
		s.WriteString("\n\n")
	}
	return panelStyle.Width(width).Render(s.String())
}
