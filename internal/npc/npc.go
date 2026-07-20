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
	ID         string
	Name       string
	Type       NPCType
	X          int
	Y          int
	LocationID int
	Dialogue   []string
	IsTrainer  bool
	IsDefeated bool
	PokemonID  int
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
		X:          12,
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
		X:          10,
		Y:          2,
		LocationID: 0,
		Dialogue: []string{
			"Hi sweetie! I'm so proud of you for starting",
			"your Pokemon journey!",
			"Remember to call home once in a while!",
		},
	})

	// bug catcher
	manager.AddNPC(&NPC{
		ID:         "bug_catcher",
		Name:       "Bug Catcher Jimmy",
		Type:       NPCTrainer,
		X:          10,
		Y:          7,
		LocationID: 1,
		IsTrainer:  true,
		IsDefeated: false,
		PokemonID:  10,
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
		X:          8,
		Y:          8,
		LocationID: 2,
		Dialogue: []string{
			"These caves are full of Geodude and Zubat!",
			"Watch your step, it's dark in here!",
		},
	})

	return manager
}
