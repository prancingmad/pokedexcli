package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "List all available commands",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the next page of 20 maps",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous page of 20 maps",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Displays all available pokemon in the location given",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a Pokemon by name",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Shows detailed information about a caught Pokemon",
			callback:    commandInspect,
		},
	}
}

type config struct {
	Next     string
	Previous string
	Caught   map[string]Pokemon
}

type mapResponse struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Pokemon struct {
	Name           string
	BaseExperience int
	Height         int
	Weight         int
	Stats          map[string]int
	Types          []string
}

func cleanInput(text string) []string {
	stepOne := strings.TrimSpace(text)
	stepTwo := strings.ToLower(stepOne)
	return strings.Fields(stepTwo)
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config, args []string) error {
	url := cfg.Next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?limit=20"
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch locations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad response: %s - %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var parsed mapResponse
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	for _, loc := range parsed.Results {
		fmt.Println(loc.Name)
	}

	cfg.Next = parsed.Next
	cfg.Previous = parsed.Previous

	return nil
}

func commandMapBack(cfg *config, args []string) error {
	url := cfg.Previous
	if url == "" {
		fmt.Println("No previous page available.")
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch previous locations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad response: %s - %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var parsed mapResponse
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	for _, loc := range parsed.Results {
		fmt.Println(loc.Name)
	}

	cfg.Next = parsed.Next
	cfg.Previous = parsed.Previous

	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: explore <location-name>")
		return nil
	}

	location := strings.Join(args, "-")

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", location)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch location: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad response: %s - %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type exploreResponse struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}

	var parsed exploreResponse
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(parsed.PokemonEncounters) == 0 {
		fmt.Println("No Pokemon found in this location.")
		return nil
	}

	fmt.Printf("Pokemon in %s:\n", location)
	for _, encounter := range parsed.PokemonEncounters {
		fmt.Println("- " + encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: catch <pokemon-name>")
		return nil
	}

	name := strings.Join(args, "-")
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", name)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch pokemon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad response: %s - %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data struct {
		Name           string `json:"name"`
		BaseExperience int    `json:"base_experience"`
		Height         int    `json:"height"`
		Weight         int    `json:"weight"`
		Stats          []struct {
			Stat struct {
				Name string `json:"name"`
			} `json:"stat"`
			BaseStat int `json:"base_stat"`
		} `json:"stats"`
		Types []struct {
			Type struct {
				Name string `json:"name"`
			} `json:"type"`
		} `json:"types"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", data.Name)

	chance := 50.0 - float64(data.BaseExperience)/2.0
	if chance < 5 {
		chance = 5
	}
	if chance > 95 {
		chance = 95
	}

	if rand.Float64()*100 < chance {
		fmt.Printf("%s was caught!\n", data.Name)
		if cfg.Caught == nil {
			cfg.Caught = make(map[string]Pokemon)
		}

		statsMap := make(map[string]int)
		for _, s := range data.Stats {
			statsMap[s.Stat.Name] = s.BaseStat
		}

		var types []string
		for _, t := range data.Types {
			types = append(types, t.Type.Name)
		}

		cfg.Caught[data.Name] = Pokemon{
			Name:           data.Name,
			BaseExperience: data.BaseExperience,
			Height:         data.Height,
			Weight:         data.Weight,
			Stats:          statsMap,
			Types:          types,
		}
	} else {
		fmt.Printf("%s escaped!\n", data.Name)
	}

	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: inspect <pokemon-name>")
		return nil
	}

	name := strings.Join(args, "-")

	pokemon, found := cfg.Caught[name]
	if !found {
		fmt.Println("You have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for statName, statValue := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", statName, statValue)
	}

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t)
	}

	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)
		if len(words) > 0 {
			cmdName := words[0]
			args := words[1:]

			if cmd, ok := commands[cmdName]; ok {
				err := cmd.callback(cfg, args)
				if err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Println("unknown command")
			}
		}
	}
}
