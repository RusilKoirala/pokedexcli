package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rusilkoirala/pokedexcli/internal/quest"
)

var (
	questPanelStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#F7F7F7"))
	questTitleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFE66D"))
	questProgressStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#95E1D3"))
)

func RenderQuestInfo(qm *quest.QuestManager) string {
	var s strings.Builder

	activeQuests := qm.GetActiveQuests()
	completedQuests := qm.GetCompletedQuests()

	if len(activeQuests) == 0 && len(completedQuests) == 0 {
		return ""
	}

	s.WriteString(questTitleStyle.Render("= QUESTS =") + "\n\n")

	shown := 0
	for _, q := range activeQuests {
		if shown >= 1 {
			break
		}
		s.WriteString(questPanelStyle.Render(q.Title) + "\n")
		progressBar := fmt.Sprintf("[%d/%d]", q.Progress, q.Target)
		s.WriteString(questProgressStyle.Render(progressBar) + "\n\n")
		shown++
	}

	if len(completedQuests) > 0 {
		s.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#95E1D3")).
			Render(fmt.Sprintf("✓ %d Quest(s) Ready!", len(completedQuests))) + "\n")
	}

	return s.String()
}
