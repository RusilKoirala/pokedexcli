package locations

import "math/rand"

type Location struct {
	Name        string
	Description string
	PokemonPool []string
	MinID       int
	MaxID       int
}

var Locations = []Location{
	{
		Name:        "Pallet Town",
		Description: "A quiet town with starter Pokemon",
		MinID:       1,
		MaxID:       20,
	},
	{
		Name:        "Viridian Forest",
		Description: "A forest full of bug-type Pokemon",
		MinID:       10,
		MaxID:       30,
	},
	{
		Name:        "Mt. Moon",
		Description: "A dark cave with rock and ground types",
		MinID:       35,
		MaxID:       75,
	},
	{
		Name:        "Route 1",
		Description: "A path with common Pokemon",
		MinID:       16,
		MaxID:       50,
	},
	{
		Name:        "Safari Zone",
		Description: "Rare Pokemon roam here!",
		MinID:       100,
		MaxID:       150,
	},
}

// get random pokemon id from a location
func (l *Location) GetRandomPokemonID() int {
	return rand.Intn(l.MaxID-l.MinID+1) + l.MinID
}

// return a location  by index
func GetLocation(index int) Location {
	if index < 0 || index >= len(Locations) {
		return Locations[0]
	}
	return Locations[index]
}

// the number of location
func GetLocationCount() int {
	return len(Locations)
}
