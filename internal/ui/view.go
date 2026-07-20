package ui

import (
	"fmt"
	"image"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/qeesung/image2ascii/convert"
	"github.com/rusilkoirala/pokedexcli/internal/locations"
	"github.com/rusilkoirala/pokedexcli/internal/town"
	"github.com/rusilkoirala/pokedexcli/internal/ui/views"
)

var (
	// Color palette - Pokemon inspired but professional
	primaryColor   = lipgloss.Color("#FF6B6B") // Pokemon Red
	secondaryColor = lipgloss.Color("#4ECDC4") // Cyan/Blue
	accentColor    = lipgloss.Color("#FFE66D") // Yellow
	successColor   = lipgloss.Color("#95E1D3") // Mint green
	textColor      = lipgloss.Color("#F7F7F7") // Off-white
	mutedColor     = lipgloss.Color("#8B9798") // Gray
	bgDark         = lipgloss.Color("#1A1A2E") // Dark bg
	bgLight        = lipgloss.Color("#16213E") // Light bg

	// Title styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			Border(lipgloss.DoubleBorder(), true).
			BorderForeground(primaryColor).
			Padding(0, 2).
			MarginBottom(1)

	// Menu item styles
	menuItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			PaddingLeft(4).
			PaddingRight(4).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(bgDark).
			Background(accentColor).
			Bold(true).
			PaddingLeft(2).
			PaddingRight(2).
			MarginBottom(1)

	// Help text
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			MarginTop(1)

	// Box styles
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1, 2).
			MarginTop(1)

	infoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Align(lipgloss.Left)

	// Status styles
	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Type colors
	typeColors = map[string]string{
		"normal":   "#A8A878",
		"fire":     "#F08030",
		"water":    "#6890F0",
		"electric": "#F8D030",
		"grass":    "#78C850",
		"ice":      "#98D8D8",
		"fighting": "#C03028",
		"poison":   "#A040A0",
		"ground":   "#E0C068",
		"flying":   "#A890F0",
		"psychic":  "#F85888",
		"bug":      "#A8B820",
		"rock":     "#B8A038",
		"ghost":    "#705898",
		"dragon":   "#7038F8",
		"dark":     "#705848",
		"steel":    "#B8B8D0",
		"fairy":    "#EE99AC",
	}
)

func convertImageToASCII(img image.Image) string {
	if img == nil {
		return "No sprite available"
	}

	convertOptions := convert.DefaultOptions
	convertOptions.FixedWidth = 40
	convertOptions.FixedHeight = 20
	convertOptions.Colored = true
	convertOptions.Reversed = false

	converter := convert.NewImageConverter()
	return converter.Image2ASCIIString(img, &convertOptions)
}

func convertImageToASCIISmall(img image.Image) string {
	if img == nil {
		return "No sprite"
	}

	convertOptions := convert.DefaultOptions
	convertOptions.FixedWidth = 20
	convertOptions.FixedHeight = 10
	convertOptions.Colored = true
	convertOptions.Reversed = false

	converter := convert.NewImageConverter()
	return converter.Image2ASCIIString(img, &convertOptions)
}

func (m Model) View() string {
	if m.loading {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(accentColor).Render("Loading..."))
	}

	var content string
	switch m.currentView {
	case startView:
		content = m.renderStartScreen()
	case starterSelectionView:
		content = views.RenderStarterSelection(m.cursor, m.width, m.height)
	case menuView:
		content = m.renderMenu()
	case creditsView:
		content = m.renderCredits()
	case listView:
		content = m.renderList()
	case overworldView:
		content = m.renderOverworldWithPanel()
	case detailView:
		content = m.renderDetail()
	case myPokedexView:
		content = m.renderMyPokedex()
	case exploreView:
		content = m.renderExplore()
	case encounterView:
		content = m.renderEncounter()
	case pokemonSelectView:
		content = m.renderPokemonSelect()
	case battleView:
		content = m.renderBattle()
	}

	if m.message != "" {
		content += "\n" + lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Render("» "+m.message)
	}

	// Center all content
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) renderStartScreen() string {
	var s strings.Builder

	// Clean ASCII Art Title
	asciiTitle := `
╔═══════════════════════════════════════════════════╗
║                                                   ║
║   ██████   ██████  ██  ██ ███████ ██████  ███████║
║   ██   ██ ██    ██ ██ ██  ██      ██   ██ ██     ║
║   ██████  ██    ██ █████   █████  ██   ██ █████  ║
║   ██      ██    ██ ██  ██  ██      ██   ██ ██     ║
║   ██       ██████  ██   ██ ███████ ██████  ███████║
║                                                   ║
║              ██████ ██      ██                    ║
║             ██      ██      ██                    ║
║             ██      ██      ██                    ║
║             ██      ██      ██                    ║
║              ██████ ███████ ██                    ║
║                                                   ║
╚═══════════════════════════════════════════════════╝
`

	titleStyled := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Render(asciiTitle)

	s.WriteString(titleStyled + "\n")

	subtitle := lipgloss.NewStyle().
		Foreground(textColor).
		Italic(true).
		Render("A Terminal Adventure")

	s.WriteString(subtitle + "\n\n")

	// Clean menu options
	options := []string{
		"Play",
		"Credits",
		"Exit",
	}

	for i, option := range options {
		if m.cursor == i {
			s.WriteString(selectedStyle.Render("▸ "+option) + "\n")
		} else {
			s.WriteString(menuItemStyle.Render("  "+option) + "\n")
		}
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Render("↑/↓: navigate  •  enter: select  •  q: quit"))
	s.WriteString("\n\n")

	version := lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("v1.0.0")

	s.WriteString(version)

	return s.String()
}

func (m Model) renderCredits() string {
	var s strings.Builder

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		Border(lipgloss.ThickBorder(), true).
		BorderForeground(primaryColor).
		Padding(0, 3).
		Render("CREDITS")

	s.WriteString(title + "\n\n")

	credits := lipgloss.NewStyle().
		Foreground(textColor).
		Render(`
Created By: Your Name

Built With:
  • Go Programming Language
  • Bubble Tea TUI Framework
  • Lipgloss Styling
  • PokeAPI

Special Thanks:
  • Pokemon Company
  • Open source community
  • You for playing!

GitHub: github.com/yourusername/pokedexcli
`)

	creditsBox := boxStyle.Render(credits)

	s.WriteString(creditsBox + "\n\n")
	s.WriteString(helpStyle.Render("press 'b' to go back"))

	return s.String()
}

func (m Model) renderMenu() string {
	var s strings.Builder

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		Border(lipgloss.RoundedBorder(), false, false, true, false).
		BorderForeground(primaryColor).
		Padding(0, 1).
		MarginBottom(2).
		Render("POKEDEX CLI")

	s.WriteString(header + "\n")

	options := []string{
		"Browse Pokemon",
		"My Pokedex",
		"Go Exploring",
		"Exit",
	}

	for i, option := range options {
		if m.cursor == i {
			s.WriteString(selectedStyle.Render("▸ "+option) + "\n")
		} else {
			s.WriteString(menuItemStyle.Render("  "+option) + "\n")
		}
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Render("↑/↓: navigate  •  enter: select  •  q: quit"))
	return s.String()
}

func (m Model) renderList() string {
	var s strings.Builder

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		Render(fmt.Sprintf("Pokemon List  │  Page %d", m.page+1))

	s.WriteString(header + "\n")
	s.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render(strings.Repeat("─", 50)) + "\n\n")

	for i, name := range m.pokemonList {
		caught := ""
		if m.pokedex.Has(name) {
			caught = lipgloss.NewStyle().Foreground(successColor).Render(" ✓")
		}

		if m.cursor == i {
			s.WriteString(selectedStyle.Render("▸ "+name) + caught + "\n")
		} else {
			s.WriteString(menuItemStyle.Render("  "+name) + caught + "\n")
		}
	}

	s.WriteString("\n" + helpStyle.Render("↑/↓: navigate  •  enter: view  •  n: next  •  p: prev  •  b: back"))
	return s.String()
}

func (m Model) renderDetail() string {
	if m.selectedPoke == nil {
		return "No Pokemon Selected"
	}

	p := m.selectedPoke

	// Title
	title := fmt.Sprintf("%s (#%d)", strings.ToUpper(p.Name), p.ID)

	// Left column - Pokemon info
	var leftCol strings.Builder

	// Types
	leftCol.WriteString("Types: ")
	for i, t := range p.Types {
		color := typeColors[t.Type.Name]
		typeStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true)
		leftCol.WriteString(typeStyle.Render(t.Type.Name))
		if i < len(p.Types)-1 {
			leftCol.WriteString(", ")
		}
	}
	leftCol.WriteString("\n\n")

	// Physical info
	leftCol.WriteString(fmt.Sprintf("Height: %.1fm\n", float64(p.Height)/10))
	leftCol.WriteString(fmt.Sprintf("Weight: %.1fkg\n", float64(p.Weight)/10))
	leftCol.WriteString(fmt.Sprintf("Base XP: %d\n\n", p.BaseExperience))

	// Stats
	leftCol.WriteString("Stats:\n")
	for _, stat := range p.Stats {
		bar := strings.Repeat("█", stat.BaseStat/10)
		leftCol.WriteString(fmt.Sprintf("  %-12s %3d %s\n", stat.Stat.Name+":", stat.BaseStat, bar))
	}

	// Right column - ASCII sprite
	var rightCol string
	if m.spriteImage != nil {
		rightCol = convertImageToASCII(m.spriteImage)
	} else {
		rightCol = "\n\n   No sprite\n   available\n"
	}

	// Style the columns
	leftStyle := lipgloss.NewStyle().
		Width(50).
		Align(lipgloss.Left).
		PaddingRight(2)

	rightStyle := lipgloss.NewStyle().
		Width(40).
		Align(lipgloss.Center)

	// Join columns side by side
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(leftCol.String()),
		rightStyle.Render(rightCol),
	)

	// Build final output
	var output strings.Builder
	output.WriteString(title + "\n\n")
	output.WriteString(content + "\n\n")

	// Catch status
	if m.pokedex.Has(p.Name) {
		output.WriteString("✓ Already caught!\n")
	} else {
		output.WriteString(helpStyle.Render("Press 'c' to catch this Pokemon!\n"))
	}

	output.WriteString("\n" + helpStyle.Render("b: back"))
	return output.String()
}

func (m Model) renderMyPokedex() string {
	var s strings.Builder

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		Render(fmt.Sprintf("My Pokedex  │  %d Caught", m.pokedex.Count()))

	s.WriteString(header + "\n")
	s.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render(strings.Repeat("─", 50)) + "\n\n")

	if len(m.pokemonList) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("Your Pokedex is empty!\nGo exploring to catch some Pokemon.")
		s.WriteString(emptyMsg + "\n")
	} else {
		for i, name := range m.pokemonList {
			if m.cursor == i {
				s.WriteString(selectedStyle.Render("▸ "+name) + "\n")
			} else {
				s.WriteString(menuItemStyle.Render("  "+name) + "\n")
			}
		}
	}

	s.WriteString("\n" + helpStyle.Render("↑/↓: navigate  •  enter: view  •  b: back"))
	return s.String()
}

func (m Model) renderExplore() string {
	var s strings.Builder

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		Render("Exploration Mode")

	s.WriteString(header + "\n")
	s.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render(strings.Repeat("─", 50)) + "\n")
	s.WriteString(helpStyle.Render("Choose a location and press 'e' to explore") + "\n\n")

	for i := 0; i < locations.GetLocationCount(); i++ {
		location := locations.GetLocation(i)

		if m.cursor == i {
			s.WriteString(selectedStyle.Render("▸ "+location.Name) + "\n")
			s.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true).
				PaddingLeft(4).
				Render(location.Description) + "\n\n")
		} else {
			s.WriteString(menuItemStyle.Render("  "+location.Name) + "\n\n")
		}
	}

	// Stats
	statsText := fmt.Sprintf("Encounters: %d  │  Caught: %d", m.totalEncounters, m.pokedex.Count())
	statsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Foreground(textColor).
		Padding(0, 2).
		Render(statsText)

	s.WriteString(statsBox + "\n\n")
	s.WriteString(helpStyle.Render("↑/↓: navigate  •  e: explore  •  b: back"))
	return s.String()
}

func (m Model) renderEncounter() string {
	if m.encounterPokemon == nil {
		return "Loading encounter..."
	}

	var s strings.Builder
	location := locations.GetLocation(m.currentLocation)

	locationHeader := lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true).
		Render(fmt.Sprintf("Location: %s", location.Name))

	s.WriteString(locationHeader + "\n\n")

	switch m.encounterState {
	case appearing, choosing:
		wildText := fmt.Sprintf("A wild %s appeared!", strings.ToUpper(m.encounterPokemon.Name))

		encounterHeader := lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Render(wildText)

		s.WriteString(encounterHeader + "\n\n")

		if m.encounterSprite != nil {
			s.WriteString(convertImageToASCII(m.encounterSprite) + "\n")
		}

		// Pokemon info line
		infoLine := fmt.Sprintf("#%03d", m.encounterPokemon.ID)
		for i, t := range m.encounterPokemon.Types {
			color := typeColors[t.Type.Name]
			typeStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(color)).
				Bold(true).
				PaddingLeft(1).
				PaddingRight(1)
			infoLine += "  " + typeStyle.Render(strings.ToUpper(t.Type.Name))
			if i < len(m.encounterPokemon.Types)-1 {
				infoLine += ""
			}
		}
		s.WriteString("\n" + infoLine + "\n\n")

		if m.pokedex.Has(m.encounterPokemon.Name) {
			s.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Render("Already in your Pokedex") + "\n\n")
		}

		catchRate := m.calculateCatchRate() * 100
		catchRateText := lipgloss.NewStyle().
			Foreground(accentColor).
			Render(fmt.Sprintf("Catch Rate: %.0f%%", catchRate))
		s.WriteString(catchRateText + "\n\n")

		actionsBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(0, 2).
			Render("[b] Battle  │  [c] Catch  │  [r] Run")
		s.WriteString(actionsBox)

	case throwing:
		animation := "You threw a Pokéball!\n\n     ●\n       →\n         ⚫"
		s.WriteString("\n" + animation + "\n")

	case shaking:
		shakes := strings.Repeat("... ", m.shakeCount+1)
		animation := "The Pokéball is shaking...\n\n" + shakes
		s.WriteString("\n" + animation + "\n")

	case caught:
		successMsg := successStyle.Render(fmt.Sprintf("Gotcha! %s was caught!", strings.ToUpper(m.encounterPokemon.Name)))

		successBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(successColor).
			Padding(1, 2).
			Render(successMsg + "\n\n" + fmt.Sprintf("%s was added to your Pokedex!", m.encounterPokemon.Name))

		s.WriteString("\n" + successBox + "\n\n")
		s.WriteString(helpStyle.Render("press 'b' to continue exploring"))

	case escaped:
		escapeMsg := errorStyle.Render(fmt.Sprintf("Oh no! %s broke free!", strings.ToUpper(m.encounterPokemon.Name)))

		escapeBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Render(escapeMsg)

		s.WriteString("\n" + escapeBox + "\n\n")
		s.WriteString(helpStyle.Render("press 'b' to continue"))
	}

	return s.String()
}

func (m Model) renderPokemonSelect() string {
	var s strings.Builder

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		Render("Choose Your Pokemon")

	s.WriteString(header + "\n")
	s.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render(strings.Repeat("─", 50)) + "\n\n")

	if len(m.pokemonList) == 0 {
		s.WriteString(helpStyle.Render("You don't have any Pokemon!\n"))
	} else {
		for i, name := range m.pokemonList {
			if m.cursor == i {
				s.WriteString(selectedStyle.Render("▸ "+name) + "\n")
			} else {
				s.WriteString(menuItemStyle.Render("  "+name) + "\n")
			}
		}
	}

	s.WriteString("\n" + helpStyle.Render("↑/↓: navigate  •  enter: select  •  b: back"))
	return s.String()
}

func (m Model) renderBattle() string {
	if m.currentBattle == nil {
		return "No battle active"
	}

	var s strings.Builder

	// Top row: Enemy Pokemon
	enemyName := strings.ToUpper(m.currentBattle.WildPokemon.Pokemon.Name)
	enemyLevel := fmt.Sprintf("Lv %d", m.currentBattle.WildPokemon.Level)

	enemyInfo := fmt.Sprintf("%s %s\nHP: %d/%d\n%s",
		enemyName,
		enemyLevel,
		m.currentBattle.WildPokemon.CurrentHP,
		m.currentBattle.WildPokemon.MaxHP,
		renderHPBar(m.currentBattle.WildPokemon.CurrentHP, m.currentBattle.WildPokemon.MaxHP, 20),
	)

	enemyInfoStyled := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Render(enemyInfo)

	enemySprite := ""
	if m.enemyBattleSprite != nil {
		enemySprite = convertImageToASCIISmall(m.enemyBattleSprite)
	}

	enemyRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		enemySprite,
		"  ",
		enemyInfoStyled,
	)

	s.WriteString(enemyRow + "\n\n")

	// Middle spacer
	s.WriteString("                    ⚔️\n\n")

	// Bottom row: Your Pokemon
	playerName := strings.ToUpper(m.currentBattle.PlayerPokemon.Pokemon.Name)
	playerLevel := fmt.Sprintf("Lv %d", m.currentBattle.PlayerPokemon.Level)

	playerInfo := fmt.Sprintf("%s %s\nHP: %d/%d\n%s",
		playerName,
		playerLevel,
		m.currentBattle.PlayerPokemon.CurrentHP,
		m.currentBattle.PlayerPokemon.MaxHP,
		renderHPBar(m.currentBattle.PlayerPokemon.CurrentHP, m.currentBattle.PlayerPokemon.MaxHP, 20),
	)

	playerInfoStyled := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(0, 1).
		Render(playerInfo)

	playerSprite := ""
	if m.playerBattleSprite != nil {
		playerSprite = convertImageToASCIISmall(m.playerBattleSprite)
	}

	playerRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		playerInfoStyled,
		"  ",
		playerSprite,
	)

	s.WriteString(playerRow + "\n\n")

	// Battle log
	if m.battleLog != "" {
		logBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			Foreground(textColor).
			Padding(0, 1).
			Width(60).
			Render(m.battleLog)
		s.WriteString(logBox + "\n\n")
	}

	// Move boxes at bottom
	if !m.currentBattle.IsOver {
		s.WriteString(m.renderMoveBoxes())
	} else {
		s.WriteString(helpStyle.Render("press 'b' to return"))
	}

	return s.String()
}

func (m Model) renderMoveBoxes() string {
	if m.currentBattle == nil || len(m.currentBattle.PlayerPokemon.Moves) == 0 {
		return ""
	}

	var boxes []string

	for i, move := range m.currentBattle.PlayerPokemon.Moves {
		var boxStyle lipgloss.Style

		// Get type color
		typeColor := typeColors[move.Type]
		if typeColor == "" {
			typeColor = "#CCCCCC"
		}

		moveContent := fmt.Sprintf("%s\n%s  PP:%d/%d",
			strings.ToUpper(move.Name),
			strings.ToUpper(move.Type),
			move.PP,
			move.MaxPP)

		if i == m.selectedMoveIndex {
			// Selected move - highlighted
			boxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(accentColor).
				Background(bgLight).
				Foreground(accentColor).
				Bold(true).
				Padding(0, 1).
				Width(18).
				Align(lipgloss.Center)
		} else {
			// Unselected move
			boxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(typeColor)).
				Foreground(textColor).
				Padding(0, 1).
				Width(18).
				Align(lipgloss.Center)
		}

		boxes = append(boxes, boxStyle.Render(moveContent))
	}

	// Join boxes horizontally with spacing
	movesRow := lipgloss.JoinHorizontal(lipgloss.Top, boxes...)

	help := helpStyle.Render("←/→: select  •  enter: attack  •  c: catch  •  r: run  •  b: back")

	return "\n" + movesRow + "\n\n" + help
}

func (m Model) renderOverworldWithPanel() string {
	if m.currentMap == nil {
		return "No map loaded"
	}

	// Calculate layout - 75% for map, 25% for right panel
	mapWidth := int(float64(m.width) * 0.75)
	rightPanelWidth := m.width - mapWidth - 4

	var s strings.Builder

	// Top: Location name (full width)
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		Align(lipgloss.Center).
		Width(m.width).
		Render("📍 " + m.currentMap.Name)

	s.WriteString(header + "\n\n")

	// Build the map content with SMART SPACING
	var mapContent strings.Builder

	for y := 0; y < m.currentMap.Height; y++ {
		rowWidth := len(m.currentMap.Tiles[y])

		// First pass: render actual row with horizontal spacing
		for x := 0; x < rowWidth; x++ {
			if x == m.playerX && y == m.playerY {
				// Player sprite - bright and bold
				playerStyle := lipgloss.NewStyle().
					Foreground(accentColor).
					Bold(true)
				mapContent.WriteString(playerStyle.Render(string(town.TilePlayer)))
			} else {
				tile := m.currentMap.Tiles[y][x]
				tileColor := getTileColor(tile)
				tileStyle := lipgloss.NewStyle().Foreground(tileColor)
				mapContent.WriteString(tileStyle.Render(string(tile)))
			}

			// Add horizontal spacing with smart fill
			if x < rowWidth-1 {
				fillChar := getSmartFill(m.currentMap, x, y)
				fillColor := getTileColor(fillChar)
				fillStyle := lipgloss.NewStyle().Foreground(fillColor)
				mapContent.WriteString(fillStyle.Render(string(fillChar)))
			}
		}
		mapContent.WriteString("\n")

		// Second pass: add vertical spacing row with smart fill
		if y < m.currentMap.Height-1 {
			for x := 0; x < rowWidth; x++ {
				fillChar := getSmartFill(m.currentMap, x, y)
				fillColor := getTileColor(fillChar)
				fillStyle := lipgloss.NewStyle().Foreground(fillColor)
				mapContent.WriteString(fillStyle.Render(string(fillChar)))

				// Also fill the diagonal/corner space
				if x < rowWidth-1 {
					mapContent.WriteString(fillStyle.Render(string(fillChar)))
				}
			}
			mapContent.WriteString("\n")
		}
	}

	// Style the map panel with border
	mapPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(2, 3).
		Width(mapWidth).
		Render(mapContent.String())

	// Right panel - player info
	rightPanel := views.RenderPlayerPanel(m.player, rightPanelWidth)

	// Join left (map) and right (empty) panels horizontally
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		mapPanel,
		"  ",
		rightPanel,
	)

	s.WriteString(mainContent)
	s.WriteString("\n\n")

	// Bottom: Controls (full width)
	controls := lipgloss.NewStyle().
		Foreground(mutedColor).
		Align(lipgloss.Center).
		Width(m.width).
		Render("WASD/Arrows: move  •  b: back  •  q: quit")

	s.WriteString(controls)

	return s.String()
}

// getSmartFill returns the appropriate fill character based on surrounding tiles
func getSmartFill(worldMap *town.WorldMap, x, y int) town.TileType {
	currentTile := worldMap.Tiles[y][x]

	// Special tiles (houses, signs, trees, etc.) should NOT be duplicated in fill
	specialTiles := map[town.TileType]bool{
		town.TileHouse:  true,
		town.TileSign:   true,
		town.TileFlower: true,
		town.TileCave:   true,
		town.TileTree:   true,
	}

	if specialTiles[currentTile] {
		// Look at neighbors to determine what to fill with
		neighbors := []town.TileType{}

		// Check right
		if x+1 < len(worldMap.Tiles[y]) {
			neighbors = append(neighbors, worldMap.Tiles[y][x+1])
		}
		// Check down
		if y+1 < worldMap.Height && x < len(worldMap.Tiles[y+1]) {
			neighbors = append(neighbors, worldMap.Tiles[y+1][x])
		}
		// Check left
		if x > 0 {
			neighbors = append(neighbors, worldMap.Tiles[y][x-1])
		}
		// Check up
		if y > 0 && x < len(worldMap.Tiles[y-1]) {
			neighbors = append(neighbors, worldMap.Tiles[y-1][x])
		}

		// Count grass vs path neighbors
		grassCount := 0
		pathCount := 0
		for _, neighbor := range neighbors {
			if neighbor == town.TileGrass {
				grassCount++
			} else if neighbor == town.TilePath {
				pathCount++
			}
		}

		// Return most common neighbor type
		if grassCount > pathCount {
			return town.TileGrass
		} else if pathCount > 0 {
			return town.TilePath
		}
		return town.TileGrass // default
	}

	// For regular tiles, just duplicate them
	return currentTile
}

func getTileColor(tile town.TileType) lipgloss.Color {
	switch tile {
	case town.TileNPC:
		return lipgloss.Color("#FF6B6B")
	case town.TileGrass:
		return lipgloss.Color("#78C850")
	case town.TilePath:
		return lipgloss.Color("#8B9798")
	case town.TileTree:
		return lipgloss.Color("#2D5016")
	case town.TileWater:
		return lipgloss.Color("#6890F0")
	case town.TileBuilding:
		return lipgloss.Color("#705848")
	case town.TileHouse:
		return lipgloss.Color("#FF6B6B")
	case town.TileCave:
		return lipgloss.Color("#4A4A4A")
	case town.TileSign:
		return lipgloss.Color("#FFE66D")
	case town.TileFlower:
		return lipgloss.Color("#EE99AC")
	case town.TileFence:
		return lipgloss.Color("#8B4513")
	default:
		return lipgloss.Color("#F7F7F7")
	}
}

// renderHPBar renders HP bar with colors
func renderHPBar(current, max, width int) string {
	if max == 0 || current < 0 {
		return strings.Repeat("░", width)
	}

	percentage := float64(current) / float64(max)
	if percentage > 1.0 {
		percentage = 1.0
	}
	if percentage < 0.0 {
		percentage = 0.0
	}

	filledWidth := int(float64(width) * percentage)
	emptyWidth := width - filledWidth

	// Ensure non-negative widths
	if filledWidth < 0 {
		filledWidth = 0
	}
	if emptyWidth < 0 {
		emptyWidth = 0
	}

	// Color based on HP percentage
	var color lipgloss.Color
	if percentage > 0.5 {
		color = successColor // Green
	} else if percentage > 0.2 {
		color = accentColor // Yellow
	} else {
		color = primaryColor // Red
	}

	filledStyle := lipgloss.NewStyle().Foreground(color)
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))

	filled := filledStyle.Render(strings.Repeat("█", filledWidth))
	empty := emptyStyle.Render(strings.Repeat("░", emptyWidth))

	return filled + empty
}
