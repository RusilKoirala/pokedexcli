package ui

import (
	"fmt"
	"image"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/qeesung/image2ascii/convert"
	"github.com/rusilkoirala/pokedexcli/internal/locations"
)

var (
	bigtitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700")).
			Background(lipgloss.Color("#C41E3A")).
			Padding(1, 4).
			MarginBottom(2).
			Align(lipgloss.Center)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Italic(true).
			Align(lipgloss.Center).
			MarginBottom(2)

	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700")).Background(lipgloss.Color("#C41E3A")).Padding(0, 2).MarginBottom(1)

	menuItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			PaddingLeft(2)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true).
			Background(lipgloss.Color("#2C5F8D")).
			PaddingLeft(1).
			PaddingRight(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8B8B8B")).
			Italic(true).
			MarginTop(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFD700")).
			Padding(1, 2).
			MarginTop(1)

	encounterBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("#FF6B6B")).
				Padding(2, 4).
				Align(lipgloss.Center)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true)

	creditsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Padding(0, 2).
			Align(lipgloss.Center)

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
	var s strings.Builder

	// Title with Pokemon-style header
	s.WriteString(titleStyle.Render("⚡ POKÉDEX CLI ⚡") + "\n")

	if m.loading {
		s.WriteString("\n  " + helpStyle.Render("Loading...") + "\n")
		return s.String()
	}

	switch m.currentView {
	case startView:
		s.WriteString(m.renderStartScreen())
	case creditsView:
		s.WriteString(m.renderCredits())
	case menuView:
		s.WriteString(m.renderMenu())
	case listView:
		s.WriteString(m.renderList())
	case detailView:
		s.WriteString(m.renderDetail())
	case myPokedexView:
		s.WriteString(m.renderMyPokedex())
	case exploreView:
		s.WriteString(m.renderExplore())
	case encounterView:
		s.WriteString(m.renderEncounter())
	case pokemonSelectView:
		s.WriteString(m.renderPokemonSelect())
	case battleView:
		s.WriteString(m.renderBattle())
	}

	if m.message != "" {
		s.WriteString("\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true).
			Render("  » "+m.message) + "\n")
	}
	return s.String()
}

func (m Model) renderStartScreen() string {
	var s strings.Builder

	s.WriteString("\n\n")
	asciiTitle := `
  ██████╗  ██████╗ ██╗  ██╗███████╗██████╗ ███████╗██╗  ██╗
  ██╔══██╗██╔═══██╗██║ ██╔╝██╔════╝██╔══██╗██╔════╝╚██╗██╔╝
  ██████╔╝██║   ██║█████╔╝ █████╗  ██║  ██║█████╗   ╚███╔╝
  ██╔═══╝ ██║   ██║██╔═██╗ ██╔══╝  ██║  ██║██╔══╝   ██╔██╗
  ██║     ╚██████╔╝██║  ██╗███████╗██████╔╝███████╗██╔╝ ██╗
  ╚═╝      ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═════╝ ╚══════╝╚═╝  ╚═╝

                      ██████╗██╗     ██╗
                     ██╔════╝██║     ██║
                     ██║     ██║     ██║
                     ██║     ██║     ██║
                     ╚██████╗███████╗██║
                      ╚═════╝╚══════╝╚═╝
`
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true).Align(lipgloss.Center).Width(80).Render(asciiTitle)

	s.WriteString(titleStyle + "\n")

	subtitle := "Catch 'Em All in Your Terminal!"
	subtitleStyled := subtitleStyle.Width(80).Render(subtitle)
	s.WriteString(subtitleStyled + "\n\n")

	// options
	options := []string{
		"Play",
		"Credits",
		"Exit",
	}

	for i, options := range options {
		if m.cursor == i {
			s.WriteString("      " + selectedStyle.Render("  ▶ "+options+"  ") + "\n\n")
		} else {
			s.WriteString("        " + menuItemStyle.Render(options) + "\n\n")
		}
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Align(lipgloss.Center).Width(80).Render("↑/↓: navigate • enter: select • q: quit"))
	s.WriteString("\n\n")

	version := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Italic(true).Align(lipgloss.Center).Width(80).Render("v1.0.0 | Made with ❤️ for Pokémon")

	s.WriteString(version)
	return s.String()
}

func (m Model) renderMenu() string {
	var s strings.Builder

	s.WriteString("\n")

	options := []string{
		"📚 Browse Pokemon",
		"🎒 My Pokedex",
		"🗺️  Go Exploring",
		"🚪 Exit",
	}

	for i, option := range options {
		cursor := "  "
		if m.cursor == i {
			s.WriteString("  " + selectedStyle.Render(" ▶ "+option+" ") + "\n")
		} else {
			s.WriteString(cursor + menuItemStyle.Render(option) + "\n")
		}
		if i < len(options)-1 {
			s.WriteString("\n")
		}
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Render("  ↑/↓: navigate • enter: select • q: quit"))
	return s.String()
}

func (m Model) renderList() string {
	var s strings.Builder

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		Render(fmt.Sprintf("📋 Pokemon List (Page %d)", m.page+1))

	s.WriteString("\n" + header + "\n\n")

	for i, name := range m.pokemonList {
		caught := ""
		if m.pokedex.Has(name) {
			caught = " ✓"
		}

		if m.cursor == i {
			s.WriteString("  " + selectedStyle.Render(" ▶ "+name+caught+" ") + "\n")
		} else {
			s.WriteString("    " + menuItemStyle.Render(name+caught) + "\n")
		}
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Render("  ↑/↓: navigate • enter: view • n: next • p: prev • b: back"))
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
		Foreground(lipgloss.Color("#FFD700")).
		Render(fmt.Sprintf("🎒 My Pokedex (%d caught)", m.pokedex.Count()))

	s.WriteString("\n" + header + "\n\n")

	if len(m.pokemonList) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8B8B8B")).
			Italic(true).
			Render("  Your Pokedex is empty!\n  Go exploring to catch some Pokemon!")
		s.WriteString(emptyMsg + "\n")
	} else {
		for i, name := range m.pokemonList {
			if m.cursor == i {
				s.WriteString("  " + selectedStyle.Render(" ▶ "+name+" ") + "\n")
			} else {
				s.WriteString("    " + menuItemStyle.Render(name) + "\n")
			}
		}
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Render("  ↑/↓: navigate • enter: view • b: back"))
	return s.String()
}

func (m Model) renderExplore() string {
	var s strings.Builder

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		Render("🗺️  Exploration Mode")

	s.WriteString("\n" + header + "\n")
	s.WriteString(helpStyle.Render("  Choose a location and press 'e' to explore") + "\n\n")

	for i := 0; i < locations.GetLocationCount(); i++ {
		location := locations.GetLocation(i)

		if m.cursor == i {
			s.WriteString("  " + selectedStyle.Render(" ▶ "+location.Name+" ") + "\n")
			s.WriteString("    " + helpStyle.Render("→ "+location.Description) + "\n\n")
		} else {
			s.WriteString("    " + menuItemStyle.Render(location.Name) + "\n\n")
		}
	}

	// Stats box
	statsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#6890F0")).
		Padding(0, 2).
		Render(fmt.Sprintf("Encounters: %d  |  Caught: %d", m.totalEncounters, m.pokedex.Count()))

	s.WriteString(statsBox + "\n\n")
	s.WriteString(helpStyle.Render("  ↑/↓: navigate • e: explore • b: back"))
	return s.String()
}

func (m Model) renderEncounter() string {
	if m.encounterPokemon == nil {
		return "Loading encounter..."
	}

	var s strings.Builder
	location := locations.GetLocation(m.currentLocation)

	locationHeader := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#78C850")).
		Bold(true).
		Render(fmt.Sprintf("📍 %s", location.Name))

	s.WriteString("\n" + locationHeader + "\n\n")

	switch m.encounterState {
	case appearing, choosing:
		// Wild Pokemon appeared box
		wildText := fmt.Sprintf("A wild %s appeared!", strings.ToUpper(m.encounterPokemon.Name))

		encounterHeader := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true).
			Align(lipgloss.Center).
			Width(50).
			Render(wildText)

		s.WriteString(encounterHeader + "\n\n")

		if m.encounterSprite != nil {
			s.WriteString(convertImageToASCII(m.encounterSprite) + "\n")
		}

		// Pokemon info line
		infoLine := fmt.Sprintf("#%03d  ", m.encounterPokemon.ID)
		for i, t := range m.encounterPokemon.Types {
			color := typeColors[t.Type.Name]
			typeStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(color)).
				Bold(true).
				Padding(0, 1).
				Background(lipgloss.Color("#2C3E50"))
			infoLine += typeStyle.Render(strings.ToUpper(t.Type.Name))
			if i < len(m.encounterPokemon.Types)-1 {
				infoLine += " "
			}
		}
		s.WriteString("\n" + infoLine + "\n\n")

		if m.pokedex.Has(m.encounterPokemon.Name) {
			s.WriteString(helpStyle.Render("  ⚠️  Already in your Pokedex") + "\n\n")
		}

		catchRate := m.calculateCatchRate() * 100
		catchRateText := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Render(fmt.Sprintf("  Catch Rate: %.0f%%", catchRate))
		s.WriteString(catchRateText + "\n\n")

		actionsBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFD700")).
			Padding(0, 2).
			Render("[b] Battle  [c] Catch  [r] Run  [esc] Back")
		s.WriteString(actionsBox)

	case throwing:
		animation := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(50).
			Render("You threw a Pokéball!\n\n     ●\n       →\n         ⚫")
		s.WriteString("\n" + animation + "\n")

	case shaking:
		shakes := strings.Repeat("... ", m.shakeCount+1)
		animation := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(50).
			Render("The Pokéball is shaking...\n\n" + shakes)
		s.WriteString("\n" + animation + "\n")

	case caught:
		successMsg := successStyle.Render(fmt.Sprintf("🎉 Gotcha! %s was caught! 🎉", strings.ToUpper(m.encounterPokemon.Name)))

		successBox := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#00FF00")).
			Padding(1, 2).
			Align(lipgloss.Center).
			Width(50).
			Render(successMsg + "\n\n" + fmt.Sprintf("%s was added to your Pokedex!", m.encounterPokemon.Name))

		s.WriteString("\n" + successBox + "\n\n")
		s.WriteString(helpStyle.Render("  Press 'b' to continue exploring"))

	case escaped:
		escapeMsg := errorStyle.Render(fmt.Sprintf("Oh no! %s broke free!", strings.ToUpper(m.encounterPokemon.Name)))

		escapeBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B6B")).
			Padding(1, 2).
			Align(lipgloss.Center).
			Width(50).
			Render(escapeMsg)

		s.WriteString("\n" + escapeBox + "\n\n")
		s.WriteString(helpStyle.Render("  Press 'b' to continue"))
	}

	return s.String()
}

func (m Model) renderPokemonSelect() string {
	var s strings.Builder

	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700")).Render("⚔️  Choose Your Pokemon")

	s.WriteString("\n" + header + "\n\n")

	if len(m.pokemonList) == 0 {
		s.WriteString(helpStyle.Render(" You don't have any Pokemon!\n"))
	} else {
		for i, name := range m.pokemonList {
			if m.cursor == i {
				s.WriteString("  " + selectedStyle.Render(" ▶ "+name+" ") + "\n")
			} else {
				s.WriteString("    " + menuItemStyle.Render(name) + "\n")
			}
		}
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Render("  ↑/↓: navigate • enter: select • b: back"))
	return s.String()
}

func (m Model) renderBattle() string {
	if m.currentBattle == nil {
		return "No battle active"
	}

	var s strings.Builder
	s.WriteString("\n")

	// Top row: Enemy Pokemon (right aligned)
	enemyName := strings.ToUpper(m.currentBattle.WildPokemon.Pokemon.Name)
	enemyLevel := fmt.Sprintf("Lv %d", m.currentBattle.WildPokemon.Level)

	// Enemy info and sprite side by side
	enemyInfo := fmt.Sprintf("%s %s\nHP: %d/%d\n%s",
		enemyName,
		enemyLevel,
		m.currentBattle.WildPokemon.CurrentHP,
		m.currentBattle.WildPokemon.MaxHP,
		renderHPBar(m.currentBattle.WildPokemon.CurrentHP, m.currentBattle.WildPokemon.MaxHP, 20),
	)

	enemyInfoStyled := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF6B6B")).
		Padding(0, 1).
		Render(enemyInfo)

	enemySprite := ""
	if m.enemyBattleSprite != nil {
		enemySprite = convertImageToASCIISmall(m.enemyBattleSprite)
	}

	// Join enemy info and sprite horizontally
	enemyRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		enemySprite,
		"  ",
		enemyInfoStyled,
	)

	// Right align the whole enemy section
	s.WriteString(lipgloss.NewStyle().
		Width(80).
		Align(lipgloss.Right).
		Render(enemyRow))

	s.WriteString("\n\n")

	// Middle spacer
	s.WriteString(strings.Repeat(" ", 32) + "⚔️\n\n")

	// Bottom row: Your Pokemon (left aligned)
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
		BorderForeground(lipgloss.Color("#4ECDC4")).
		Padding(0, 1).
		Render(playerInfo)

	playerSprite := ""
	if m.playerBattleSprite != nil {
		playerSprite = convertImageToASCIISmall(m.playerBattleSprite)
	}

	// Join player info and sprite horizontally
	playerRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		playerInfoStyled,
		"  ",
		playerSprite,
	)

	s.WriteString(playerRow)
	s.WriteString("\n\n")

	// Battle log
	if m.battleLog != "" {
		logBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#888888")).
			Padding(0, 1).
			Width(60).
			Render(m.battleLog)
		s.WriteString(logBox + "\n\n")
	}

	// Move boxes at bottom
	if !m.currentBattle.IsOver {
		s.WriteString(m.renderMoveBoxes())
	} else {
		s.WriteString(helpStyle.Render("  Press 'b' to return"))
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

		typeColor := typeColors[move.Type]
		if typeColor == "" {
			typeColor = "#CCCCCC"
		}

		moveContent := fmt.Sprintf("%s\n%s | PP:%d/%d",
			strings.ToUpper(move.Name),
			strings.ToUpper(move.Type),
			move.PP,
			move.MaxPP)

		if i == m.selectedMoveIndex {
			// Selected move - highlighted
			boxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FFD700")).
				Background(lipgloss.Color("#2C5F8D")).
				Foreground(lipgloss.Color("#FFD700")).
				Bold(true).
				Padding(0, 1).
				Width(16).
				Align(lipgloss.Center)
		} else {
			// Unselected move
			boxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(typeColor)).
				Foreground(lipgloss.Color("#CCCCCC")).
				Padding(0, 1).
				Width(16).
				Align(lipgloss.Center)
		}

		boxes = append(boxes, boxStyle.Render(moveContent))
	}

	// Join boxes horizontally with spacing
	movesRow := lipgloss.JoinHorizontal(lipgloss.Top, boxes...)

	help := helpStyle.Render("  ←/→: select move • enter: attack • c: catch • r: run • b: back")

	return "\n" + movesRow + "\n\n" + help
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

	// ensure non-negative widths
	if filledWidth < 0 {
		filledWidth = 0
	}
	if emptyWidth < 0 {
		emptyWidth = 0
	}

	// color based on HP percentage
	var color lipgloss.Color
	if percentage > 0.5 {
		color = lipgloss.Color("#00FF00") // green
	} else if percentage > 0.2 {
		color = lipgloss.Color("#FFD700") // yellow
	} else {
		color = lipgloss.Color("#FF0000") // red
	}

	filledStyle := lipgloss.NewStyle().Foreground(color)
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))

	filled := filledStyle.Render(strings.Repeat("█", filledWidth))
	empty := emptyStyle.Render(strings.Repeat("░", emptyWidth))

	return filled + empty
}
func (m Model) renderCredits() string {
	var s strings.Builder
	s.WriteString("\n\n")

	title := bigtitleStyle.Width(80).Render("CREDITS")
	s.WriteString(title + "\n\n")

	credits := `
	Created By : Rusil Koirala

	Built With : 
	Go 
	PokeAPI
	Bubble Tea TUI Framework
	Lipgloss for Styling

	Github: github.com/rusilkoirala/pokedexcli
	`

	creditsBox := boxStyle.
		Width(60).
		Align(lipgloss.Center).
		Render(creditsStyle.Render(credits))

	s.WriteString(lipgloss.NewStyle().Width(80).Align(lipgloss.Center).Render(creditsBox))
	s.WriteString("\n\n")
	s.WriteString(helpStyle.Align(lipgloss.Center).Width(80).Render("Press 'b' to go back"))

	return s.String()
}
