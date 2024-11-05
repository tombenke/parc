package main

import (
	"fmt"
	"github.com/tombenke/parc"
)

func main() {
	parc.Debug(2)

	// There are several format of inputs,
	// and the first segment before the `:` character
	// determines how to parse the rest of the input string.
	stringInput := "string:Hello"
	numberInput := "number:42"
	dicerollInput := "diceroll:2d8"

	// These are the parsers for the second part of the input string after the `:` character.

	// It parses a string value
	stringParser := parc.Letters

	// it parses a number value
	numberParser := parc.Digits

	// It parses a special format of dice-roll value
	dicerollParser := parc.SequenceOf(
		parc.Integer,
		parc.Char("d"),
		parc.Integer,
	)

	// This is the parser for the first part of the input string including the `:` separator character
	prefixParser := parc.SequenceOf(parc.Letters, parc.Char(":"))

	// This is the main parser, that chains the result of the prefix parser,
	// and the chain function selects that which parser should it use for the rest of the input string.
	parser := parc.Chain(
		prefixParser,
		func(result parc.Result) *parc.Parser {
			arr := result.([]parc.Result)
			leftValue := arr[0].(string)
			switch leftValue {
			case "string":
				return stringParser
			case "number":
				return numberParser
			default:
				return dicerollParser
			}
		})

	// Parse string input
	resultState := parser.Parse(&stringInput)
	fmt.Printf("%+v\n", resultState)

	// => inputString: 'string:Hello', Results: inputString: 'string:Hello', Results: Hello, Index: 12, Err: <nil>, IsError: false, Index: 7, Err: <nil>, IsError: false

	// Parse number input
	resultState = parser.Parse(&numberInput)
	fmt.Printf("%+v\n", resultState)

	// => inputString: 'number:42', Results: inputString: 'number:42', Results: 42, Index: 9, Err: <nil>, IsError: false, Index: 7, Err: <nil>, IsError: false

	// Parse diceroll input
	resultState = parser.Parse(&dicerollInput)
	fmt.Printf("%+v\n", resultState)

	// => inputString: 'diceroll:2d8', Results: inputString: 'diceroll:2d8', Results: [2 d 8], Index: 12, Err: <nil>, IsError: false, Index: 9, Err: <nil>, IsError: false
}
