package main

import (
	"fmt"
	"github.com/tombenke/parc"
)

func main() {
	parc.Debug(3)

	inputWithText := "Hello World"
	inputWithNumbers := "1342 234 45"

	// The choice parser takes either letters or digits
	choiceParser := parc.Choice(
		parc.Letters,
		parc.Digits,
	)

	// The parser can parse the input string if it begins with letters
	resultState := choiceParser.Parse(&inputWithText)
	fmt.Printf("%+v\n", resultState)

	// => inputString: 'Hello World', Results: Hello, Index: 5, Err: <nil>, IsError: false

	// The parser also can parse the input string if it begins with digits
	resultState = choiceParser.Parse(&inputWithNumbers)
	fmt.Printf("%+v\n", resultState)

	// => inputString: '1342 234 45', Results: 1342, Index: 4, Err: <nil>, IsError: false
}
