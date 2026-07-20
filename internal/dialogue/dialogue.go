package dialogue

import "time"

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

// typing animation
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

	d.State = DialogueWaiting
	return true
}
