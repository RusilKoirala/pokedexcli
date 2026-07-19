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
	tiles := make([][]TileType, len(lines))

	for y, line := range lines {
		tiles[y] = make([]TileType, len([]rune(line)))
		for x, char := range line {
			tiles[y][x] = TileType(char)
		}
	}
	return tiles
}

func GetMap(locationID int) *WorldMap {
	maps := []WorldMap{}
	if locationID < 0 || locationID >= len(maps) {
		return &maps[0]
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
		Width:        len(tiles[0]),
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
		Width:        len(tiles[0]),
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
		Width:        len(tiles[0]),
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
		Width:        len(tiles[0]),
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
		Width:        len(tiles[0]),
		Height:       len(tiles),
		Tiles:        tiles,
		StartX:       12,
		StartY:       12,
		LocationID:   4,
		MinPokemonID: 100,
		MaxPokemonID: 150,
	}
}
