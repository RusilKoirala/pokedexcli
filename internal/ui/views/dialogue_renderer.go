package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rusilkoirala/pokedexcli/internal/dialogue"
	"github.com/rusilkoirala/pokedexcli/internal/ui/layout"
)

var (
	dialogueBoxBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#FFE66D")).
		Padding(1, 2)

	dialogueSpeakerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFE66D")).
		Bold(true)

	dialogueTextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F7F7F7"))

	dialogueIndicatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFE66D")).
		Bold(true)
)

// RenderDialogue renders the fixed-size dialogue box at bottom
func RenderDialogue(dialogueBox *dialogue.DialogueBox) string {
	if dialogueBox == nil {
		// Empty dialogue space when not active
		emptyBox := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8B9798")).
			Italic(true).
			Align(lipgloss.Center).
			Render("Press E near an NPC to talk")
		
		return dialogueBoxBorderStyle.
			Width(layout.DialogueWidth).
			Height(layout.DialogueHeight - 4).
			Render(emptyBox)
	}

	var content strings.Builder

	// Speaker name
	content.WriteString(dialogueSpeakerStyle.Render(dialogueBox.SpeakerName + ":") + "\n\n")

	// Dialogue text (wrapped to fit)
	wrappedText := dialogueBox.GetDisplayText(layout.DialogueWidth - 10)
	
	// Show max 3 lines of dialogue
	lineCount := 0
	for _, line := range wrappedText {
		if lineCount < 3 {
			content.WriteString(dialogueTextStyle.Render(line) + "\n")
			lineCount++
		}
	}
	
	// Fill remaining lines if less than 3
	for lineCount < 3 {
		content.WriteString("\n")
		lineCount++
	}

	// Progress indicator
	if dialogueBox.GetProgress() != "" {
		indicator := dialogueIndicatorStyle.Render("▼ Press SPACE to continue")
		content.WriteString("\n" + lipgloss.NewStyle().Align(lipgloss.Right).Width(layout.DialogueWidth - 10).Render(indicator))
	}

	return dialogueBoxBorderStyle.
		Width(layout.DialogueWidth).
		Height(layout.DialogueHeight - 4).
		Render(content.String())
}
