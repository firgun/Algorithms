package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"
)

type Card struct {
	front, back string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: deck <file-name>.deck")
		os.Exit(1)
	}
	
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("deck: failed to load deck: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	
	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("deck: failed to load deck: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(string(b), "\n")
	nonEmptyLines := []string{}
	for _, ln := range lines {
		ln = strings.Trim(ln, " \t\n\r\v")
		if len(ln) == 0 {
			continue
		}
		nonEmptyLines = append(nonEmptyLines, ln)
	}

	if len(nonEmptyLines) % 2 != 0 {
		fmt.Printf("deck: failed to load deck: invalid format")
		os.Exit(1)
	}

	cards := []Card{}
	for i := 0; i < len(nonEmptyLines); i += 2 {
		cards = append(cards, Card{nonEmptyLines[i], nonEmptyLines[i+1]})
	}

	for index, card := range cards {
		var d string
		fmt.Println(card.front)
		fmt.Println("\npress return to reveal back")
		fmt.Scanln(&d)
		fmt.Println(card.back)
		if index != len(cards)-1 {
			fmt.Println("\npress return for next card")
			fmt.Scanln(&d)
		}
	}
}