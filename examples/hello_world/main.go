package main

import (
	"fmt"
	"github.com/tombenke/parc"
)

func main() {
	RunHelloParser()

	RunSequenceParser()
}

func RunHelloParser() {
	input := "Hello World"
	results := parc.Str("Hello").Parse(&input)
	fmt.Printf("%+v\n", results)
}

func RunSequenceParser() {
	input := "Hello World"
	sequenceParser := parc.SequenceOf(

		parc.Str("Hello "),
		parc.Str("World"),
	)
	results := sequenceParser.Parse(&input)
	fmt.Printf("%+v\n", results)
}
