package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
	}
}

type config struct {
	Next     string
	Previous string
}

type mapResponse struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func cleanInput(text string) []string {
	stepOne := strings.TrimSpace(text)
	stepTwo := strings.ToLower(stepOne)
	return strings.Fields(stepTwo)
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config) error {
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

func commandMapBack(cfg *config) error {
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
			if cmd, ok := commands[cmdName]; ok {
				err := cmd.callback(cfg)
				if err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Println("unknown command")
			}
		}
	}
}
