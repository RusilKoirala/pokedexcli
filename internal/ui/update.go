package ui

import (
	"fmt"
	"image"
	"math/rand"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rusilkoirala/pokedexcli/internal/battle"
	"github.com/rusilkoirala/pokedexcli/internal/dialogue"
	"github.com/rusilkoirala/pokedexcli/internal/locations"
	"github.com/rusilkoirala/pokedexcli/internal/player"
	"github.com/rusilkoirala/pokedexcli/internal/pokeapi"
	"github.com/rusilkoirala/pokedexcli/internal/town"
)

type pokemonListMsg struct {
	list []string
}

type pokemonDetailMsg struct {
	pokemon *pokeapi.Pokemon
	sprite  image.Image
}

type encounterMsg struct {
	pokemon *pokeapi.Pokemon
	sprite  image.Image
}

type dialogueTickMsg time.Time

type battleStartMsg struct {
	battle       *battle.Battle
	playerSprite image.Image
	enemySprite  image.Image
}

type battleActionMsg struct {
	message string
}

type errorMsg struct {
	err error
}

type tickMsg time.Time

func dialogueTick() tea.Cmd {
	return tea.Tick(30*time.Millisecond, func(t time.Time) tea.Msg {
		return dialogueTickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case dialogueTickMsg:
		if m.dialogueActive && m.currentDialogue != nil {
			if m.currentDialogue.Update() {
				return m, dialogueTick()
			}
		}
		return m, dialogueTick()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "w", "up":
			if m.currentView == starterSelectionView {
				if m.cursor > 0 {
					m.cursor--
				}
				return m, nil
			}
			if m.currentView == overworldView {
				return m.handleMove(0, -1)
			}
			// existing up movement for menus
			switch m.currentView {
			case pokemonSelectView:
				if m.cursor > 0 {
					m.cursor--
				}
			default:
				if m.cursor > 0 {
					m.cursor--
				}
			}

		case "s", "down":
			if m.currentView == starterSelectionView {
				if m.cursor < 2 { // 3 starters (0,1,2)
					m.cursor++
				}
				return m, nil
			}
			if m.currentView == overworldView {
				return m.handleMove(0, 1)
			}
			// existing down movement for menus
			switch m.currentView {
			case startView:
				if m.cursor < 2 {
					m.cursor++
				}
			case menuView:
				if m.cursor < 3 {
					m.cursor++
				}
			case listView:
				if m.cursor < len(m.pokemonList)-1 {
					m.cursor++
				}
			case exploreView:
				if m.cursor < locations.GetLocationCount()-1 {
					m.cursor++
				}
			case pokemonSelectView:
				if m.cursor < len(m.pokemonList)-1 {
					m.cursor++
				}
			}

		case "a", "left":
			if m.currentView == overworldView {
				return m.handleMove(-1, 0)
			}
			// existing left movement for battle
			if m.currentView == battleView {
				if m.selectedMoveIndex > 0 {
					m.selectedMoveIndex--
				}
			}

		case "d", "right":
			if m.currentView == overworldView {
				return m.handleMove(1, 0)
			}
			// existing right movement for battle
			if m.currentView == battleView && m.currentBattle != nil {
				maxMoves := len(m.currentBattle.PlayerPokemon.Moves) - 1
				if m.selectedMoveIndex < maxMoves {
					m.selectedMoveIndex++
				}
			}

		case "enter":
			if m.currentView == starterSelectionView {
				return m.handleStarterSelection()
			}
			if m.currentView == battleView {
				return m.handleBattleAction()
			}

			if m.currentView == exploreView {
				m.currentLocation = m.cursor
				m.currentMap = town.GetMap(m.currentLocation)
				m.playerX = m.currentMap.StartX
				m.playerY = m.currentMap.StartY
				m.stepCount = 0
				m.encounterSteps = 0
				m.currentView = overworldView
				return m, nil
			}
			return m.handleEnter()

		case "b", "esc":
			if m.currentView == encounterView && m.encounterState == choosing {
				m.currentView = pokemonSelectView
				m.cursor = 0
				m.pokemonList = m.pokedex.List()
				if len(m.pokemonList) == 0 {
					m.message = "You don't have any Pokemon to battle with!"
					return m, nil
				}
				return m, nil
			} else if m.currentView == battleView {
				return m.handleBack()
			}
			return m.handleBack()

		case "c":
			if m.currentView == detailView && m.selectedPoke != nil {
				return m.handleCatch()
			} else if m.currentView == encounterView && m.encounterState == choosing {
				return m.handleCatchWild()
			} else if m.currentView == battleView {
				return m, m.executeBattleCatch()
			}
		case " ", "space":
			if m.dialogueActive && m.currentDialogue != nil {
				hasMore := m.currentDialogue.NextLine()
				if !hasMore {
					m.dialogueActive = false
					m.currentDialogue = nil
				}
				return m, nil
			}

		case "e":
			if m.currentView == exploreView {
				return m.handleExplore()
			}

			if m.currentView == overworldView && !m.dialogueActive {
				return m.handleTalkToNPC()
			}

		case "r":
			if m.currentView == encounterView && m.encounterState == choosing {
				return m.handleRun()
			} else if m.currentView == battleView {
				m.message = "You ran away!"
				m.currentView = overworldView // Return to overworld instead of exploreView
				m.currentBattle = nil
			}

		case "i":
			if m.currentView == encounterView && m.encounterState == choosing {
				m.message = "Inspecting Pokemon..."
			}

		case "n":
			if m.currentView == listView {
				m.page++
				m.cursor = 0
				return m, m.loadPokemonList()
			}

		case "p":
			if m.currentView == listView && m.page > 0 {
				m.page--
				m.cursor = 0
				return m, m.loadPokemonList()
			}
		}

	case pokemonListMsg:
		m.pokemonList = msg.list
		m.loading = false

	case pokemonDetailMsg:
		m.selectedPoke = msg.pokemon
		m.spriteImage = msg.sprite
		m.loading = false

	case encounterMsg:
		m.encounterPokemon = msg.pokemon
		m.encounterSprite = msg.sprite
		m.encounterState = appearing
		m.loading = false
		m.totalEncounters++
		return m, tick()

	case battleStartMsg:
		m.currentBattle = msg.battle
		m.playerBattleSprite = msg.playerSprite
		m.enemyBattleSprite = msg.enemySprite
		m.currentView = battleView
		m.selectedMoveIndex = 0
		m.loading = false
		m.battleLog = "Battle started!"

	case battleActionMsg:
		m.battleLog = msg.message
		m.loading = false

		if m.currentBattle != nil && m.currentBattle.IsOver {
			if m.currentBattle.PlayerWon {
				m.message = "You won the battle!"
			} else {
				m.message = "You lost the battle!"
			}
		}

	case tickMsg:
		return m.handleTick()

	case errorMsg:
		m.message = fmt.Sprintf("Error: %v", msg.err)
		m.loading = false
	}

	return m, nil
}

func (m Model) handleTalkToNPC() (tea.Model, tea.Cmd) {
	if m.currentMap == nil {
		return m, nil
	}

	positions := []struct{ dx, dy int }{
		{0, -1},
		{0, 1},
		{-1, 0},
		{1, 0},
	}

	for _, pos := range positions {
		checkX := m.playerX + pos.dx
		checkY := m.playerY + pos.dy

		npcFound := m.npcManager.GetNPCAt(checkX, checkY, m.currentLocation)
		if npcFound != nil {
			m.currentDialogue = dialogue.NewDialogueBox(npcFound.Name, npcFound.Dialogue)
			m.dialogueActive = true
			return m, dialogueTick()
		}

	}
	m.message = "No one here to talk to!"
	return m, nil
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case startView:
		switch m.cursor {
		case 0: // Play
			// Check if player has starter
			if !m.player.HasStarter {
				m.currentView = starterSelectionView
				m.cursor = 0
			} else {
				m.currentView = menuView
				m.cursor = 0
			}
		case 1:
			m.currentView = creditsView
			m.cursor = 0
		case 2:
			m.pokedex.Save()
			m.player.Save()
			return m, tea.Quit
		}
	case menuView:
		switch m.cursor {
		case 0:
			m.currentView = listView
			m.cursor = 0
			m.loading = true
			return m, m.loadPokemonList()
		case 1:
			m.currentView = myPokedexView
			m.cursor = 0
			m.pokemonList = m.pokedex.List()
		case 2:
			m.currentView = exploreView
			m.cursor = 0
			m.message = ""
		case 3:
			m.pokedex.Save()
			return m, tea.Quit
		}
	case listView:
		if len(m.pokemonList) > 0 {
			m.currentView = detailView
			m.loading = true
			return m, m.loadPokemonDetail(m.pokemonList[m.cursor])
		}
	case myPokedexView:
		if len(m.pokemonList) > 0 {
			m.currentView = detailView
			m.loading = true
			return m, m.loadPokemonDetail(m.pokemonList[m.cursor])
		}

	case pokemonSelectView:
		if len(m.pokemonList) > 0 {
			playerPokemonName := m.pokemonList[m.cursor]
			m.loading = true
			return m, m.loadPlayerPokemonForBattle(playerPokemonName)
		}

	case exploreView:
		m.currentLocation = m.cursor
	}
	return m, nil
}

func (m Model) handleBack() (tea.Model, tea.Cmd) {
	m.message = ""
	switch m.currentView {
	case creditsView:
		m.currentView = startView
		m.cursor = 0
	case listView, myPokedexView, exploreView:
		m.currentView = menuView
		m.cursor = 0

	case detailView:
		m.currentView = listView
		m.cursor = 0

	case pokemonSelectView:
		m.currentView = encounterView
		m.encounterState = choosing
		m.cursor = 0

	case overworldView:
		m.currentView = exploreView
		m.cursor = 0
		m.currentMap = nil

	case battleView:
		m.currentView = overworldView // Return to overworld instead of exploreView
		m.currentBattle = nil

	case encounterView:
		if m.encounterState == caught || m.encounterState == escaped {
			m.currentView = overworldView // Return to overworld instead of exploreView
			m.encounterState = appearing
		}
	}

	return m, nil
}

func (m Model) handleCatch() (tea.Model, tea.Cmd) {
	if m.pokedex.Has(m.selectedPoke.Name) {
		m.message = fmt.Sprintf("%s is already in your Pokedex!!", m.selectedPoke.Name)
	} else if m.pokedex.Catch(m.selectedPoke.Name) {
		m.message = fmt.Sprintf("YAY you caught %s!", m.selectedPoke.Name)
		m.pokedex.Save()
	} else {
		m.message = fmt.Sprintf("%s escaped! Nice tryyy", m.selectedPoke.Name)
	}
	return m, nil
}

func (m Model) loadPokemonList() tea.Cmd {
	return func() tea.Msg {
		list, err := m.api.GetPokemonList(20, m.page*20)
		if err != nil {
			return errorMsg{err}
		}
		names := make([]string, len(list.Results))
		for i, p := range list.Results {
			names[i] = p.Name
		}
		return pokemonListMsg{names}
	}
}

func (m Model) loadPokemonDetail(name string) tea.Cmd {
	return func() tea.Msg {
		pokemon, err := m.api.GetPokemon(name)
		if err != nil {
			return errorMsg{err}
		}

		var sprite image.Image
		if pokemon.Sprites.FrontDefault != "" {
			sprite, _ = m.api.DownloadSprite(pokemon.Sprites.FrontDefault)
		}

		return pokemonDetailMsg{
			pokemon: pokemon,
			sprite:  sprite,
		}
	}
}

func tick() tea.Cmd {
	return tea.Tick(800*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) loadEncounter(pokemonID int) tea.Cmd {
	return func() tea.Msg {
		pokemon, err := m.api.GetPokemon(strconv.Itoa(pokemonID))
		if err != nil {
			return errorMsg{err}
		}

		var sprite image.Image
		if pokemon.Sprites.FrontDefault != "" {
			sprite, _ = m.api.DownloadSprite(pokemon.Sprites.FrontDefault)
		}

		return encounterMsg{
			pokemon: pokemon,
			sprite:  sprite,
		}
	}
}

func (m Model) calculateCatchRate() float64 {
	baseCatchRate := 0.4
	dexBonus := float64(m.pokedex.Count()) * 0.01
	total := baseCatchRate + dexBonus

	if total > 0.9 {
		total = 0.9
	}

	return total
}

func (m Model) handleRun() (tea.Model, tea.Cmd) {
	// Immediately return to overworld, no message screen
	m.currentView = overworldView
	m.encounterState = appearing
	return m, nil
}

func (m Model) handleCatchWild() (tea.Model, tea.Cmd) {
	m.encounterState = throwing
	m.shakeCount = 0
	return m, tick()
}

func (m Model) handleTick() (tea.Model, tea.Cmd) {
	switch m.encounterState {
	case appearing:
		m.encounterState = choosing
		return m, nil

	case throwing:
		m.encounterState = shaking
		m.shakeCount = 0
		return m, tick()

	case shaking:
		m.shakeCount++
		if m.shakeCount < 3 {
			return m, tick()
		}

		catchRate := m.calculateCatchRate()
		if rand.Float64() < catchRate {
			m.encounterState = caught
			m.pokedex.Catch(m.encounterPokemon.Name)
			m.pokedex.Save()

			// Give XP for catching
			xp := player.GetXPForAction("catch")
			leveledUp := m.player.GainXP(xp)
			m.player.Save()

			if leveledUp {
				m.message = fmt.Sprintf("Gotcha! %s was caught! +%d XP! Level up! Now Level %d!", m.encounterPokemon.Name, xp, m.player.Level)
			} else {
				m.message = fmt.Sprintf("Gotcha! %s was caught! +%d XP!", m.encounterPokemon.Name, xp)
			}
		} else {
			m.encounterState = escaped
			m.message = fmt.Sprintf("%s broke free and escaped!", m.encounterPokemon.Name)
		}

		return m, nil

	case escaped:
		time.Sleep(2 * time.Second)
		m.currentView = overworldView // Return to overworld instead of exploreView
		m.encounterState = appearing
		return m, nil
	}
	return m, nil
}

func (m Model) handleExplore() (tea.Model, tea.Cmd) {
	location := locations.GetLocation((m.currentLocation))
	pokemonID := location.GetRandomPokemonID()

	m.currentView = encounterView
	m.loading = true
	m.message = ""

	return m, m.loadEncounter(pokemonID)
}

// Fetch moves from API
func (m Model) loadMovesForPokemon(pokemon *pokeapi.Pokemon) []*battle.Move {
	moves := []*battle.Move{}

	// Take first 4 moves from pokemon
	maxMoves := 4
	if len(pokemon.Moves) < maxMoves {
		maxMoves = len(pokemon.Moves)
	}

	for i := 0; i < maxMoves; i++ {
		moveName := pokemon.Moves[i].Move.Name

		// Fetch move details from API
		move, err := m.api.GetMove(moveName)
		if err == nil && move.Power > 0 {
			battleMove := battle.NewMove(
				move.Name,
				move.Type.Name,
				move.Power,
				move.Accuracy,
				move.PP,
			)
			moves = append(moves, battleMove)
		}
	}

	// If no moves found, use defaults
	if len(moves) == 0 {
		moves = battle.GetDefaultMoves()
	}

	return moves
}

func (m Model) loadPlayerPokemonForBattle(name string) tea.Cmd {
	return func() tea.Msg {
		playerPokemon, err := m.api.GetPokemon(name)
		if err != nil {
			return errorMsg{err}
		}

		var playerSprite image.Image
		if playerPokemon.Sprites.FrontDefault != "" {
			playerSprite, _ = m.api.DownloadSprite(playerPokemon.Sprites.FrontDefault)
		}

		var enemySprite image.Image
		if m.encounterPokemon.Sprites.FrontDefault != "" {
			enemySprite, _ = m.api.DownloadSprite(m.encounterPokemon.Sprites.FrontDefault)
		}

		// Load moves for both Pokemon
		playerMoves := m.loadMovesForPokemon(playerPokemon)
		enemyMoves := m.loadMovesForPokemon(m.encounterPokemon)

		newBattle := battle.NewBattle(playerPokemon, m.encounterPokemon, playerMoves, enemyMoves)

		return battleStartMsg{
			battle:       newBattle,
			playerSprite: playerSprite,
			enemySprite:  enemySprite,
		}
	}
}

func (m Model) executeBattleAttack() tea.Cmd {
	return func() tea.Msg {
		if m.currentBattle == nil {
			return battleActionMsg{
				message: "No battle active",
			}
		}

		// Player attack with selected move
		playerMsg := m.currentBattle.PlayerAttack(m.selectedMoveIndex)

		if m.currentBattle.IsOver {
			// Give XP based on win/lose
			var xp int
			if m.currentBattle.PlayerWon {
				xp = player.GetXPForAction("battle_win")
			} else {
				xp = player.GetXPForAction("battle_lose")
			}

			leveledUp := m.player.GainXP(xp)
			m.player.Save()

			if leveledUp {
				return battleActionMsg{
					message: playerMsg + fmt.Sprintf("\n+%d XP! Level up! Now Level %d!", xp, m.player.Level),
				}
			}

			return battleActionMsg{
				message: playerMsg + fmt.Sprintf("\n+%d XP!", xp),
			}
		}

		// Enemy attack
		enemyMsg := m.currentBattle.EnemyAttack()

		return battleActionMsg{
			message: playerMsg + "\n" + enemyMsg,
		}
	}
}

func (m Model) executeBattleCatch() tea.Cmd {
	return func() tea.Msg {
		if m.currentBattle == nil {
			return battleActionMsg{
				message: "No battle active",
			}
		}

		catchRate := m.currentBattle.GetCatchRate()

		if rand.Float64() < catchRate {
			m.pokedex.Catch(m.encounterPokemon.Name)
			m.pokedex.Save()

			// Give XP for catching in battle
			xp := player.GetXPForAction("catch")
			leveledUp := m.player.GainXP(xp)
			m.player.Save()

			m.currentBattle.IsOver = true
			m.currentBattle.PlayerWon = true

			if leveledUp {
				return battleActionMsg{
					message: fmt.Sprintf("You caught %s! +%d XP! Level up! Now Level %d!", m.encounterPokemon.Name, xp, m.player.Level),
				}
			}

			return battleActionMsg{
				message: fmt.Sprintf("You caught %s! +%d XP!", m.encounterPokemon.Name, xp),
			}
		}
		return battleActionMsg{
			message: m.encounterPokemon.Name + " broke free!",
		}
	}
}

func (m Model) handleBattleAction() (tea.Model, tea.Cmd) {
	if m.currentBattle == nil || m.currentBattle.IsOver {
		return m, nil
	}

	return m, m.executeBattleAttack()
}

// move playerr
func (m Model) handleMove(dx, dy int) (tea.Model, tea.Cmd) {
	if m.currentMap == nil {
		return m, nil
	}

	newX := m.playerX + dx
	newY := m.playerY + dy

	// check if npc there
	if m.npcManager.IsNPCPosition(newX, newY, m.currentLocation) {
		return m, nil
	}

	// check if movement valid
	if m.currentMap.IsWalkable(newX, newY) {
		m.playerX = newX
		m.playerY = newY
		m.stepCount++
		m.encounterSteps++

		if m.currentMap.IsGrass(newX, newY) {
			if m.encounterSteps >= 5+rand.Intn(5) {
				m.encounterSteps = 0
				return m.triggerWildEncounter()
			}
		}
	}
	return m, nil
}

// trigger encounterr
func (m Model) triggerWildEncounter() (tea.Model, tea.Cmd) {

	if m.currentMap == nil {
		return m, nil
	}

	//get random pokemon on basis of map
	pokemonID := rand.Intn(m.currentMap.MaxPokemonID-m.currentMap.MinPokemonID+1) + m.currentMap.MinPokemonID

	m.currentView = encounterView
	m.loading = true
	m.message = ""

	return m, m.loadEncounter(pokemonID)
}

// handleStarterSelection processes starter choice
func (m Model) handleStarterSelection() (tea.Model, tea.Cmd) {
	// Import views package to access Starters
	starter := []struct {
		Name string
	}{
		{"Charmander"},
		{"Bulbasaur"},
		{"Squirtle"},
	}[m.cursor]

	// Set the starter
	m.player.SetStarter(starter.Name)
	m.player.Save()

	// Add starter to Pokedex!
	m.pokedex.Catch(starter.Name)
	m.pokedex.Save()

	m.message = fmt.Sprintf("You chose %s! Your journey begins!", starter.Name)
	m.currentView = menuView
	m.cursor = 0

	return m, nil
}
