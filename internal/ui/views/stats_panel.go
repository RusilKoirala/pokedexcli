package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rusilkoirala/pokedexcli/internal/player"
	"github.com/rusilkoirala/pokedexcli/internal/quest"
	"github.com/rusilkoirala/pokedexcli/internal/ui/layout"
)

var (
	statsPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4ECDC4")).
			Padding(1, 1)

	statsTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFE66D")).
			Bold(true).
			Align(lipgloss.Center)

	statsTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F7F7F7"))

	statsXPBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#95E1D3"))
)

// renders the fixed size stats panel
func RenderStatsPanel(p *player.Player, qm *quest.QuestManager) string {
	var s strings.Builder

	// Title
	s.WriteString(statsTitleStyle.Render("═ TRAINER ═") + "\n\n")

	// Player info
	s.WriteString(statsTextStyle.Render(fmt.Sprintf("Name: %s", p.Name)) + "\n")
	s.WriteString(statsTextStyle.Render(fmt.Sprintf("Level: %d", p.Level)) + "\n\n")

	// XP bar
	xpText := fmt.Sprintf("XP: %d/%d", p.XP, p.XPToNextLevel)
	s.WriteString(statsTextStyle.Render(xpText) + "\n")
	s.WriteString(statsXPBarStyle.Render(p.GetXPBar(18)) + "\n\n")

	// Starter
	if p.HasStarter {
		s.WriteString(statsTitleStyle.Render("═ PARTNER ═") + "\n")
		s.WriteString(statsTextStyle.Render(fmt.Sprintf("⚡ %s", p.StarterPokemon)) + "\n\n")
	}

	questInfo := RenderQuestInfo(qm)
	if questInfo != "" {
		s.WriteString("\n")
		s.WriteString(questInfo)
	}

	// Spacing
	s.WriteString("\n\n")

	// Controls
	s.WriteString(statsTitleStyle.Render("═ CONTROLS ═") + "\n\n")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#8B9798")).Render("WASD: Move") + "\n")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#8B9798")).Render("E: Talk") + "\n")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#8B9798")).Render("B: Back") + "\n")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#8B9798")).Render("Q: Quit") + "\n")

	return statsPanelStyle.
		Width(layout.StatsWidth).
		Height(layout.StatsHeight).
		Render(s.String())
}
