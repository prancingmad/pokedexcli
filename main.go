package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

type cliCommand struct {
	name		string
	description	string
	callback	func() error
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
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)
		if len(words) > 0 {
			cmdName := words[0]
			if cmd, ok := commands[cmdName]; ok {
				err := cmd.callback()
				if err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Println("unknown command")
			}
		}
	}
}

func cleanInput(text string) []string {
	stepOne := strings.TrimSpace(text)
	stepTwo := strings.ToLower(stepOne)
	return strings.Fields(stepTwo)
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd:= range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}