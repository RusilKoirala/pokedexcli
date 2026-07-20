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

// RenderMap renders the game map with NO duplicates, fixed size
func RenderMap(worldMap *town.WorldMap, playerX, playerY int, npcMgr *npc.NPCManager, locationID int) string {
	if worldMap == nil {
		return mapBorderStyle.
			Width(layout.GameWidth).
			Height(layout.GameHeight).
			Render("No map loaded")
	}
	
	var mapContent strings.Builder
	
	// Get NPCs for current location
	npcsHere := npcMgr.GetNPCsForLocation(locationID)
	npcPositions := make(map[string]bool)
	for _, npc := range npcsHere {
		key := fmt.Sprintf("%d,%d", npc.X, npc.Y)
		npcPositions[key] = true
	}
	
	// Render ONLY actual tiles - NO SPACING, NO DUPLICATES
	for y := 0; y < worldMap.Height; y++ {
		rowWidth := len(worldMap.Tiles[y])
		
		for x := 0; x < rowWidth; x++ {
			key := fmt.Sprintf("%d,%d", x, y)
			
			var char string
			var color lipgloss.Color
			
			if x == playerX && y == playerY {
				// Player sprite
				char = "⚙"
				color = lipgloss.Color("#FFE66D")
			} else if npcPositions[key] {
				// NPC sprite
				char = "Φ"
				color = lipgloss.Color("#FF6B6B")
			} else {
				// Regular tile
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
	
	// Fixed size with border and padding
	return mapBorderStyle.
		Width(layout.GameWidth).
		Height(layout.GameHeight).
		Render(mapContent.String())
}

// GetTileColor returns color for each tile type
func GetTileColor(tile town.TileType) lipgloss.Color {
	switch tile {
	case town.TileGrass:
		return lipgloss.Color("#78C850") // Green
	case town.TilePath:
		return lipgloss.Color("#8B9798") // Gray
	case town.TileTree:
		return lipgloss.Color("#2D5016") // Dark green
	case town.TileWater:
		return lipgloss.Color("#6890F0") // Blue
	case town.TileBuilding:
		return lipgloss.Color("#705848") // Brown
	case town.TileHouse:
		return lipgloss.Color("#FF6B6B") // Red
	case town.TileCave:
		return lipgloss.Color("#4A4A4A") // Dark gray
	case town.TileSign:
		return lipgloss.Color("#FFE66D") // Yellow
	case town.TileFlower:
		return lipgloss.Color("#EE99AC") // Pink
	case town.TileFence:
		return lipgloss.Color("#8B4513") // Brown
	case town.TileNPC:
		return lipgloss.Color("#FF6B6B") // Red
	default:
		return lipgloss.Color("#F7F7F7") // White
	}
}
