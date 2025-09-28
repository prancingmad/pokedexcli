package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)
		if len(words) > 0 {
			fmt.Println("Your command was:", words[0])
		}
	}
}

func cleanInput(text string) []string {
	stepOne := strings.TrimSpace(text)
	stepTwo := strings.ToLower(stepOne)
	return strings.Fields(stepTwo)
}
