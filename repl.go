package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func StartRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("> ")
		scanner.Scan()
		text := scanner.Text()
		cleaned := cleanInput(text)

		if len(cleaned) == 0 {
			continue
		}

		command := cleaned[0]

		switch command {
		case "help":
			fmt.Println("Pokedex Menu")
			fmt.Println("1. - help")
			fmt.Println("2. - exit")
			fmt.Println(" ")
		case "exit":
			fmt.Println("Exiting...")
			os.Exit(0)

		default:
			fmt.Println("Invalid Command")
		}

	}

}

func cleanInput(str string) []string {
	loweredString := strings.ToLower(str)
	words := strings.Fields(loweredString)
	return words
}
