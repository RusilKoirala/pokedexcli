package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rusilkoirala/pokedexcli/internal/dialogue"
)

var (
	dialogueBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FFE66D")).
				Padding(1, 2).
				Width(70)

	dialogueSpeakerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFE66D")).
				Bold(true)

	dialogueTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F7F7F7"))

	dialogueProgressStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFE66D")).
				Align(lipgloss.Right)
)

// render the dialogue box at the bottom
func RenderDialogueBox(dialogueBox *dialogue.DialogueBox, width int) string {
	if dialogueBox == nil {
		return ""
	}

	var content strings.Builder

	// speaker name
	speakerLine := dialogueSpeakerStyle.Render(dialogueBox.SpeakerName + ":")
	content.WriteString(speakerLine + "\n\n")

	// progress indicator
	progress := dialogueBox.GetProgress()
	if progress != "" {
		content.WriteString("\n")
		content.WriteString(dialogueProgressStyle.Render(progress + " Press Space to continue"))
	}

	boxWidth := width - 1
	if boxWidth > 80 {
		boxWidth = 80
	}
	return dialogueBoxStyle.Width(boxWidth).Render(content.String())
}
