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
		X:          6,
		Y:          5,
		LocationID: 0,
		Dialogue: []string{
			"Ah, hello there! Welcome to the world of Pokemon!",
			"My name is Oak. People call me the Pokemon Professor!",
			"Head north to Viridian Forest to catch Pokemon!",
			"Walk the path — but watch out for trainers!",
		},
	})

	// MOM
	manager.AddNPC(&NPC{
		ID:          "mom",
		Name:        "Mom",
		Type:        NPCTownsFolk,
		X:           13,
		Y:           6,
		PokemonID:   10,
		LineOfSight: 5,
		LocationID:  0,
		Dialogue: []string{
			"Be careful out there sweetie!",
			"Come back home if you get hurt!",
		},
	})

	manager.AddNPC(&NPC{
		ID: "bug_catcher", Name: "Bug Catcher Jimmy",
		Type: NPCTrainer, X: 13, Y: 6, LocationID: 1,
		IsTrainer: true, IsDefeated: false,
		PokemonID: 10, LineOfSight: 5,
		Dialogue: []string{
			"Hey! Stop right there!",
			"Nobody passes without battling me first!",
		},
	})

	// (viridian Forest)  Lass Lisa

	manager.AddNPC(&NPC{
		ID: "lass_lisa", Name: "Lass Lisa",
		Type: NPCTrainer, X: 23, Y: 13, LocationID: 1,
		IsTrainer: true, IsDefeated: false,
		PokemonID: 35, LineOfSight: 5, // Clefairy
		Dialogue: []string{
			"La la la — oh! A challenger!",
			"My Clefairy will defeat you!",
		},
	})

	// Mt. Moon — Hiker Dan

	manager.AddNPC(&NPC{
		ID: "hiker_dan", Name: "Hiker Dan",
		Type: NPCTrainer, X: 7, Y: 9, LocationID: 2,
		IsTrainer: true, IsDefeated: false,
		PokemonID: 74, LineOfSight: 5, // Geodude
		Dialogue: []string{
			"You dare enter MY cave?!",
			"My Geodude has trained in these tunnels for years!",
		},
	})

	// Route 1 — Youngster Joey

	manager.AddNPC(&NPC{
		ID: "youngster_joey", Name: "Youngster Joey",
		Type: NPCTrainer, X: 6, Y: 4, LocationID: 3,
		IsTrainer: true, IsDefeated: false,
		PokemonID: 19, LineOfSight: 6, // Rattata
		Dialogue: []string{
			"Yo! I've been waiting for a challenger!",
			"My Rattata is in the top percentage of all Rattata!",
		},
	})
	// Route 1 — Lass Iris

	manager.AddNPC(&NPC{
		ID: "lass_iris", Name: "Lass Iris",
		Type: NPCTrainer, X: 18, Y: 14, LocationID: 3,
		IsTrainer: true, IsDefeated: false,
		PokemonID: 39, LineOfSight: 6, // Jigglypuff
		Dialogue: []string{
			"Oh! You surprised me!",
			"Now I'm going to have to battle you!",
		},
	})
	// Safari Zone — Ranger Kim

	manager.AddNPC(&NPC{
		ID: "ranger_kim", Name: "Ranger Kim",
		Type: NPCTrainer, X: 10, Y: 7, LocationID: 4,
		IsTrainer: true, IsDefeated: false,
		PokemonID: 128, LineOfSight: 6, // Tauros
		Dialogue: []string{
			"HALT! This is a protected safari zone!",
			"Prove yourself worthy before exploring further!",
		},
	})
	return manager
}
