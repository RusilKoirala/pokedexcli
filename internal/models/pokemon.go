package models

// pokemon from pokeapi
type Pokemon struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	Height         int              `json:"height"`
	Weight         int              `json:"Weight"`
	BaseExperience int              `json:"weight"`
	Types          []TypeSlot       `json:"type"`
	Abilites       []AbilitySlot    `json:"abilities"`
	Stats          []StatSlot       `json:"stats"`
	Sprites        Sprites          `json:"sprites"`
	Species        NamedAPIResourse `json:"species"`
}

// type of pokemon
type TypeSlot struct {
	Slot int              `json:"slot"`
	Type NamedAPIResourse `json:"type"`
}

// ability of pokemon
type AbilitySlot struct {
	IsHidden bool             `json:"is_hidden"`
	Slot     int              `json:"slot"`
	Ability  NamedAPIResourse `json:"ability"`
}

// pokemon stat
type StatSlot struct {
	BaseStat int              `json:"base_stat"`
	Effort   int              `json:"effort"`
	Stat     NamedAPIResourse `json:"stat"`
}

// sprites of pokemon
type Sprites struct {
	FrontDefault string `json:"front_default"`
	BackDefault  string `json:"back_default"`
	FrontShiny   string `json:"front_shiny"`
	BackShiny    string `json:"back_shiny"`
}

// api resource
type NamedAPIResourse struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// pokemon listtt
type PokemonList struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []NamedAPIResourse `json:"results"`
}
