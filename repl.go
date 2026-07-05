package main

import (
	"bufio"
	"fmt"
	"os"
)

func StartRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("> ")
		scanner.Scan()
		text := scanner.Text()
		fmt.Println(text)
	}

}
