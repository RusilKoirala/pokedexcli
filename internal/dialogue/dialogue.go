package dialogue

import (
	"strings"
	"time"
)

type DialogueState int

const (
	DialogueHidden DialogueState = iota
	DialogueTyping
	DialogueWaiting
	DialogueComplete
)

type DialogueBox struct {
	CurrentLine   int
	Lines         []string
	DisplayedText string
	FullText      string
	CharIndex     int
	State         DialogueState
	TypingSpeed   time.Duration
	LastUpdate    time.Time
	SpeakerName   string
}

func NewDialogueBox(speakerName string, lines []string) *DialogueBox {
	return &DialogueBox{
		CurrentLine:   0,
		Lines:         lines,
		State:         DialogueTyping,
		TypingSpeed:   30 * time.Millisecond,
		LastUpdate:    time.Now(),
		SpeakerName:   speakerName,
		DisplayedText: "",
		FullText:      lines[0],
		CharIndex:     0,
	}
}

// Update advances the typing animation
func (d *DialogueBox) Update() bool {
	if d.State != DialogueTyping {
		return false
	}

	now := time.Now()
	if now.Sub(d.LastUpdate) < d.TypingSpeed {
		return false
	}

	d.LastUpdate = now

	if d.CharIndex < len(d.FullText) {
		d.CharIndex++
		d.DisplayedText = d.FullText[:d.CharIndex]
		return true
	}

	// Finished typing current line
	d.State = DialogueWaiting
	return true
}

// NextLine advances to the next dialogue line
func (d *DialogueBox) NextLine() bool {
	if d.State == DialogueTyping {
		// Skip typing animation, show full text
		d.DisplayedText = d.FullText
		d.CharIndex = len(d.FullText)
		d.State = DialogueWaiting
		return true
	}

	if d.CurrentLine < len(d.Lines)-1 {
		d.CurrentLine++
		d.FullText = d.Lines[d.CurrentLine]
		d.DisplayedText = ""
		d.CharIndex = 0
		d.State = DialogueTyping
		return true
	}

	// Dialogue complete
	d.State = DialogueComplete
	return false
}

// IsComplete returns true if all dialogue is shown
func (d *DialogueBox) IsComplete() bool {
	return d.State == DialogueComplete
}

// GetDisplayText returns the currently displayed text with word wrap
func (d *DialogueBox) GetDisplayText(maxWidth int) []string {
	return WrapText(d.DisplayedText, maxWidth)
}

// WrapText wraps text to fit within maxWidth
func WrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	lines := []string{}
	currentLine := ""

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) <= maxWidth {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// GetProgress returns a visual indicator of dialogue progress
func (d *DialogueBox) GetProgress() string {
	if d.State == DialogueWaiting {
		return "▼"
	}
	return ""
}
