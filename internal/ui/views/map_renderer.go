package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rusilkoirala/pokedexcli/internal/npc"
	"github.com/rusilkoirala/pokedexcli/internal/town"
	"github.com/rusilkoirala/pokedexcli/internal/ui/layout"
)

var (
	mapBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4ECDC4")).
		Padding(1, 2)
)

// render the freaking map
func RenderMap(worldMap *town.WorldMap, playerX, playerY int, npcMgr *npc.NPCManager, locationID int) string {
	if worldMap == nil {
		return mapBorderStyle.
			Width(layout.GameWidth).
			Height(layout.GameHeight).
			Render("No map loaded")
	}

	var mapContent strings.Builder

	npcsHere := npcMgr.GetNPCsForLocation(locationID)
	npcPositions := make(map[string]bool)
	trainerPositions := make(map[string]bool)
	for _, n := range npcsHere {
		key := fmt.Sprintf("%d,%d", n.X, n.Y)
		npcPositions[key] = true
		if n.IsTrainer && !n.IsDefeated {
			aboveKey := fmt.Sprintf("%d,%d", n.X, n.Y-1)
			trainerPositions[aboveKey] = true
		}
	}

	for y := 0; y < worldMap.Height; y++ {
		rowWidth := len(worldMap.Tiles[y])

		for x := 0; x < rowWidth; x++ {
			key := fmt.Sprintf("%d,%d", x, y)

			var char string
			var color lipgloss.Color

			if x == playerX && y == playerY {

				char = "⚙"
				color = lipgloss.Color("#FFE66D")
			} else if npcPositions[key] {

				char = "Φ"
				color = lipgloss.Color("#FF6B6B")
			} else if trainerPositions[key] {
				char = "!"
				color = lipgloss.Color("#FF0000")
			} else {

				tile := worldMap.Tiles[y][x]
				char = string(tile)
				color = GetTileColor(tile)
			}

			style := lipgloss.NewStyle().Foreground(color).Bold(true)
			mapContent.WriteString(style.Render(char))
		}

		if y < worldMap.Height-1 {
			mapContent.WriteString("\n")
		}
	}

	return mapBorderStyle.
		Width(layout.GameWidth).
		Height(layout.GameHeight).
		Render(mapContent.String())
}

func GetTileColor(tile town.TileType) lipgloss.Color {
	switch tile {
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
	case town.TileNPC:
		return lipgloss.Color("#FF6B6B")
	default:
		return lipgloss.Color("#F7F7F7")
	}
}
