package battle

type Move struct {
	Name     string
	Type     string
	Power    int
	Accuracy int
	PP       int
	MaxPP    int
}

func NewMove(name, moveType string, power, accuracy, maxPP int) *Move {
	return &Move{
		Name:     name,
		Type:     moveType,
		Power:    power,
		Accuracy: accuracy,
		PP:       maxPP,
		MaxPP:    maxPP,
	}
}

func (m *Move) CanUse() bool {
	return m.PP > 0
}

func (m *Move) Use() {
	if m.PP > 0 {
		m.PP--
	}
}

func GetDefaultMoves() []*Move {
	return []*Move{
		NewMove("Tackle", "normal", 40, 100, 35),
		NewMove("Quick Attack", "normal", 40, 100, 30),
		NewMove("Scratch", "normal", 40, 100, 35),
		NewMove("Pound", "normal", 40, 100, 35),
	}
}
