package pokeapi

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	_ "image/jpeg"
	"io"
	"net/http"
	"time"
)

const baseUrl = "https://pokeapi.co/api/v2"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type Pokemon struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	BaseExperience int           `json:"base_experience"`
	Types          []PokemonType `json:"types"`
	Stats          []PokemonStat `json:"stats"`
	Sprites        PokemonSprites `json:"sprites"`
}

type PokemonSprites struct {
	FrontDefault string `json:"front_default"`
}

type PokemonType struct {
	Slot int `json:"slot"`
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type PokemonStat struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type PokemonListResponse struct {
	Count   int `json:"count"`
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// fetch list of pokemon
func (c *Client) GetPokemonList(limit, offset int) (*PokemonListResponse, error) {
	url := fmt.Sprintf("%s/pokemon?limt=%d&offset=%d", baseUrl, limit, offset)

	resp, err := c.httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result PokemonListResponse

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// fetch detailed pokemon data
func (c *Client) GetPokemon(name string) (*Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%s", baseUrl, name)
	resp, err := c.httpClient.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("pokemon not found")
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var pokemon Pokemon

	if err := json.Unmarshal(body, &pokemon); err != nil {
		return nil, err

	}
	return &pokemon, nil
}

// to download sprite
func (c *Client) DownloadSprite(url string) (image.Image, error) {
	if url == "" {
		return nil, fmt.Errorf("no sprite url")
	}

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}
