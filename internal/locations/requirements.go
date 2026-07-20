package locations

type AreaRequirement struct {
	LocationID    int
	RequiredLevel int
	Name          string
	Description   string
}

var AreaRequirements = []AreaRequirement{
	{
		LocationID:    0,
		RequiredLevel: 1,
		Name:          "Pallet Town",
		Description:   "Your starting town",
	},
	{
		LocationID:    1,
		RequiredLevel: 3,
		Name:          "Viridian Forest",
		Description:   "A forest full of bug Pokemon",
	},
	{
		LocationID:    2,
		RequiredLevel: 8,
		Name:          "Mt. Moon",
		Description:   "A dark cave with rock Pokemon",
	},
	{
		LocationID:    3,
		RequiredLevel: 5,
		Name:          "Route 1",
		Description:   "A path with common Pokemon",
	},
	{
		LocationID:    4,
		RequiredLevel: 10,
		Name:          "Safari Zone",
		Description:   "Rare Pokemon roam here!",
	},
}

// return the requirement of the loaction
func GetRequirement(locationID int) *AreaRequirement {
	for i := range AreaRequirements {
		if AreaRequirements[i].LocationID == locationID {
			return &AreaRequirements[i]
		}
	}
	return nil
}

// check if play can access a location
func CanAccess(locationID int, playerLevel int) bool {
	req := GetRequirement(locationID)
	if req == nil {
		return true
	}
	return playerLevel >= req.RequiredLevel
}

// get message for locked area
func GetLockMessage(locationID int, playerLevel int) string {
	req := GetRequirement(locationID)
	if !CanAccess(locationID, playerLevel) {
		return "🔒 Requires Level " + string(rune(req.RequiredLevel+'0'))
	}
	if req == nil {
		return ""
	}
	return ""
}
