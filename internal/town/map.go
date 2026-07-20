package town

import (
	"strings"
)

type TileType rune

const (
	TileGrass    TileType = '░'
	TilePath     TileType = '▓'
	TileTree     TileType = '♣'
	TileWater    TileType = '~'
	TileBuilding TileType = '■'
	TileHouse    TileType = '🏠'
	TileCave     TileType = '◘'
	TileSign     TileType = '⚑'
	TileFlower   TileType = '❀'
	TileFence    TileType = '═'
	TilePlayer   TileType = '⚙'
	TileNPC      TileType = 'Φ'
)

type WorldMap struct {
	Name         string
	Width        int
	Height       int
	Tiles        [][]TileType
	StartX       int
	StartY       int
	LocationID   int
	MinPokemonID int
	MaxPokemonID int
}

// check if its walkable
func (m *WorldMap) IsWalkable(x, y int) bool {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return false
	}

	tile := m.Tiles[y][x]
	walkable := tile == TileGrass || tile == TilePath || tile == TileSign || tile == TileFlower
	return walkable
}

// tile triggers encounters
func (m *WorldMap) IsGrass(x, y int) bool {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return false
	}
	return m.Tiles[y][x] == TileGrass
}

// parse map converts ascii art to string
func parseMap(mapString string) [][]TileType {
	lines := strings.Split(strings.TrimSpace(mapString), "\n")

	// Filter out empty lines
	nonEmptyLines := []string{}
	for _, line := range lines {
		if len(strings.TrimSpace(line)) > 0 {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	if len(nonEmptyLines) == 0 {
		return [][]TileType{}
	}

	tiles := make([][]TileType, len(nonEmptyLines))

	for y, line := range nonEmptyLines {
		runes := []rune(line)
		tiles[y] = make([]TileType, len(runes))
		for x, char := range runes {
			tiles[y][x] = TileType(char)
		}
	}

	return tiles
}

// getMapWidth returns the maximum width across all rows
func getMapWidth(tiles [][]TileType) int {
	maxWidth := 0
	for _, row := range tiles {
		if len(row) > maxWidth {
			maxWidth = len(row)
		}
	}
	return maxWidth
}

func GetMap(locationID int) *WorldMap {
	maps := []WorldMap{
		PalletTownMap(),
		ViridianForestMap(),
		MtMoonMap(),
		Route1Map(),
		SafariZoneMap(),
	}

	if locationID < 0 || locationID >= len(maps) {
		palletTown := PalletTownMap()
		return &palletTown
	}

	return &maps[locationID]
}

func PalletTownMap() WorldMap {
	mapArt := `
■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
■🏠░░░░░░░░░░░░░░░░░░░░░░🏠■
■░░░░░░░░░░░░░░░░░░░░░░░░░░■
■░░░❀░░░░░░⚑░░░░░░░░❀░░░░■
■░░░░░░░░░░░░░░░░░░░░░░░░░░■
■░░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░■
■░░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░■
■░░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░■
■♣♣♣░░▓▓▓▓▓▓▓▓▓▓▓▓▓░░♣♣♣■
■♣♣♣░░▓▓▓▓▓▓▓▓▓▓▓▓▓░░♣♣♣■
■░░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░■
■░░░░░░░░░░░░░░░░░░░░░░░░░░■
■░░❀░░░░░░░░░░░░░░░░░❀░░░■
■░░░░░░░░░░🏠░░░░░░░░░░░░░■
■■■■■■■■■■■═══■■■■■■■■■■■■
                ▓▓▓
                ▓▓▓
                ▓▓▓`

	tiles := parseMap(mapArt)

	return WorldMap{
		Name:         "Pallet Town",
		Width:        getMapWidth(tiles),
		Height:       len(tiles),
		Tiles:        tiles,
		StartX:       15,
		StartY:       10,
		LocationID:   0,
		MinPokemonID: 1,
		MaxPokemonID: 20,
	}
}

// ViridianForestMap creates the Viridian Forest map
func ViridianForestMap() WorldMap {
	mapArt := `
♣♣♣♣♣♣♣♣♣♣♣♣♣♣♣♣♣♣♣♣♣♣♣♣♣
♣░░░░░░♣♣♣░░░░░♣♣♣░░░░░♣
♣░░❀░░░♣♣░░❀░░░♣♣░░░░░♣
♣░░░░░░░♣░░░░░░░♣░░░░░░♣
♣♣░░░░░░░░░░░░░░░░░░░░♣♣
♣♣♣░░░░░░▓▓▓░░░░░░░░♣♣♣
♣♣░░░░░░░▓▓▓░░░░░░░░░░♣♣
♣░░░░░░░░▓▓▓░░░░░░░░░░░♣
♣░░░♣♣░░░░░░░░░♣♣░░░░░♣
♣░░░♣♣♣░░░░░░♣♣♣░░❀░░♣
♣░░░░♣♣░░⚑░░░♣♣░░░░░░♣
♣░░░░░░░░░░░░░░░░░░░░░░♣
♣░❀░░░░░░░░░░░░░░░░░░░░♣
♣░░░░░░░░░░░░░░░░░░░░░░♣
♣♣♣♣♣♣♣♣♣♣▓▓▓♣♣♣♣♣♣♣♣♣
              ▓▓▓
              ▓▓▓`

	tiles := parseMap(mapArt)

	return WorldMap{
		Name:         "Viridian Forest",
		Width:        getMapWidth(tiles),
		Height:       len(tiles),
		Tiles:        tiles,
		StartX:       12,
		StartY:       13,
		LocationID:   1,
		MinPokemonID: 10,
		MaxPokemonID: 30,
	}
}

// MtMoonMap creates Mt. Moon cave map
func MtMoonMap() WorldMap {
	mapArt := `
■■■■■■■■■■■■■■■■■■■■■■■■■■
■░░░░░░■■■░░░░░░░░■■■░░■
■░░❀░░░░░░░░░░░░░░░░░░░■
■░░░░░░░░░░░░░░░░░░░░░░■
■░░░░■■░░░░░░░░░■■░░░░■
■░░░░■■░░░░░░░░░■■░░░░■
■░░░░░░░░░░◘░░░░░░░░░░■
■░░░░░░░░░░░░░░░░░░░░░■
■░░░░░░░■■■■■░░░░░░░░░■
■░░░░░░░░░░░░░░░░░░░░░■
■░░░░░░░░░░░░░░░░░░❀░■
■░░░░░░░░░⚑░░░░░░░░░░■
■░░░░░░░░░░░░░░░░░░░░░■
■■■■■■■■■▓▓▓■■■■■■■■■■■
            ▓▓▓
            ▓▓▓`

	tiles := parseMap(mapArt)

	return WorldMap{
		Name:         "Mt. Moon",
		Width:        getMapWidth(tiles),
		Height:       len(tiles),
		Tiles:        tiles,
		StartX:       12,
		StartY:       11,
		LocationID:   2,
		MinPokemonID: 35,
		MaxPokemonID: 75,
	}
}

// route1Map creates Route 1 map
func Route1Map() WorldMap {
	mapArt := `
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
▓░░░░░░░░░░░░░░░░░░░░░░▓
▓░░░░░░░░░░░░░░░░░░░░░░▓
▓░░❀░░░░░░░░░░░░░░❀░░░▓
▓░░░░░░░░░░░░░░░░░░░░░░▓
▓░░░░░░░░░░░░░░░░░░░░░░▓
▓░░░░░░░░⚑░░░░░░░░░░░░▓
▓░░░░░░░░░░░░░░░░░░░░░░▓
▓░░░░░░░░░░░░░░░░░░░░░░▓
▓░░░░░░░░░░░░░░░░░░░░░░▓
▓░░░░░░░░░░░░░░░░░░░░░░▓
▓░░❀░░░░░░░░░░░░░░❀░░░▓
▓░░░░░░░░░░░░░░░░░░░░░░▓
▓░░░░░░░░░░░░░░░░░░░░░░▓
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓`

	tiles := parseMap(mapArt)

	return WorldMap{
		Name:         "Route 1",
		Width:        getMapWidth(tiles),
		Height:       len(tiles),
		Tiles:        tiles,
		StartX:       12,
		StartY:       12,
		LocationID:   3,
		MinPokemonID: 16,
		MaxPokemonID: 50,
	}
}

// safariZoneMap creates Safari Zone map
func SafariZoneMap() WorldMap {
	mapArt := `
═══════════════════════════
═░░░░░░░~~~~~░░░░░░░░░░░═
═░░❀░░░~~~~~░░░❀░░░░░░░═
═░░░░░░~~~~~░░░░░░░░░░░░═
═░░░░░░~~~~~░░░░░░░░░░░░═
═░░░░░░░░░░░░░░░░░░░░░░░═
═░░░░♣♣♣░░⚑░░░♣♣♣░░░░░═
═░░░░♣♣♣░░░░░░░♣♣♣░░░░░═
═░░░░░░░░░░░░░░░░░░░░░░░═
═░░░░░░░░░░░░░░░░░░~~~~~═
═░░❀░░░░░░░░░░░░░░~~~~~═
═░░░░░░░░░░░░░░░░░~~~~~═
═░░░░░░░░░░░░░░░░░░░░░░░═
═░░░░░░░░░░░░░░░░░░░░░░░═
═══════════════════════════`

	tiles := parseMap(mapArt)

	return WorldMap{
		Name:         "Safari Zone",
		Width:        getMapWidth(tiles),
		Height:       len(tiles),
		Tiles:        tiles,
		StartX:       12,
		StartY:       12,
		LocationID:   4,
		MinPokemonID: 100,
		MaxPokemonID: 150,
	}
}
