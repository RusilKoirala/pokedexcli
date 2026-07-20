package npc

type NPCType string

const (
	NPCProffessor NPCType = "professor"
	NPCTrainer    NPCType = "trainer"
	NPCTownsFolk  NPCType = "townsfolk"
	NPCRival      NPCType = "rival"
)

// npc
type NPC struct {
	ID          string
	Name        string
	Type        NPCType
	X           int
	Y           int
	LocationID  int
	Dialogue    []string
	IsTrainer   bool
	IsDefeated  bool
	PokemonID   int
	LineOfSight int
}

// it handles all npcs
type NPCManager struct {
	NPCs map[string]*NPC
}

func NewNPCManager() *NPCManager {
	return &NPCManager{
		NPCs: make(map[string]*NPC),
	}
}

// adds npc to the manager
func (m *NPCManager) AddNPC(npc *NPC) {
	m.NPCs[npc.ID] = npc
}

// gives position of npc
func (m *NPCManager) GetNPCAt(x, y int, locationID int) *NPC {
	for _, npc := range m.NPCs {
		if npc.X == x && npc.Y == y && npc.LocationID == locationID {
			return npc
		}
	}
	return nil
}

// returns all npcs location
func (m *NPCManager) GetNPCsForLocation(locationID int) []*NPC {
	npcs := []*NPC{}
	for _, npc := range m.NPCs {
		if npc.LocationID == locationID {
			npcs = append(npcs, npc)
		}
	}
	return npcs
}

// check if there is npc at given location
func (m *NPCManager) IsNPCPosition(x, y int, locationID int) bool {
	return m.GetNPCAt(x, y, locationID) != nil
}

/*
this functions check in all direction if there is player or not , if player is found :D we send our trainer to fight otherwise we leave
*/
func (m *NPCManager) GetTrainerInSight(playerX, playerY, locationID int) *NPC {
	for _, npc := range m.NPCs {
		if npc.LocationID != locationID || !npc.IsTrainer || npc.IsDefeated || npc.LineOfSight == 0 {
			continue
		}

		directions := []struct{ dx, dy int }{
			{0, -1},
			{0, 1},
			{-1, 0},
			{1, 0},
		}

		for _, d := range directions {
			for dist := 1; dist <= npc.LineOfSight; dist++ {
				checkX := npc.X + d.dx*dist
				checkY := npc.Y + d.dy*dist

				if checkX == playerX && checkY == playerY {
					return npc
				}
			}
		}
	}
	return nil
}

// create all npc for game
func InitializeNPCs() *NPCManager {
	manager := NewNPCManager()

	/*
		I could have made a new file where i could put all npc data and loop into it to initialize but i will have less characters so i got lazy

		I will do it if i get more npcs
	*/

	// Professor Oak
	manager.AddNPC(&NPC{
		ID:         "prof_oak",
		Name:       "Professor Oak",
		Type:       NPCProffessor,
		X:          5,
		Y:          6,
		LocationID: 0,
		Dialogue: []string{
			"Ah, hello there! Welcome to the world of Pokemon!",
			"My name is Oak. People call me the Pokemon Professor!",
			"This world is inhabited by creatures called Pokemon!",
			"Your adventure is just beginning. Good luck!",
		},
	})

	// MOM
	manager.AddNPC(&NPC{
		ID:         "mom",
		Name:       "Mom",
		Type:       NPCTownsFolk,
		X:          4,
		Y:          17,
		LocationID: 0,
		Dialogue: []string{
			"Hi sweetie! I'm so proud of you for starting",
			"your Pokemon journey!",
			"Remember to call home once in a while!",
		},
	})

	// bug catcher
	manager.AddNPC(&NPC{
		ID:          "bug_catcher",
		Name:        "Bug Catcher Jimmy",
		Type:        NPCTrainer,
		X:           12,
		Y:           10,
		LocationID:  1,
		IsTrainer:   true,
		IsDefeated:  false,
		PokemonID:   10,
		LineOfSight: 3,
		Dialogue: []string{
			"Hey! I love bug Pokemon!",
			"Have you caught any Caterpie or Weedle yet?",
			"They're the best!",
		},
	})

	// hiker in mt. moon
	manager.AddNPC(&NPC{
		ID:         "hiker",
		Name:       "Hiker Dan",
		Type:       NPCTownsFolk,
		X:          14,
		Y:          9,
		LocationID: 2,
		Dialogue: []string{
			"These caves are full of Geodude and Zubat!",
			"Watch your step, it's dark in here!",
		},
	})

	return manager
}
