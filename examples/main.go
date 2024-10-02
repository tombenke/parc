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
	results := parc.Str("Hello").Parse("Hello World")
	fmt.Printf("%+v\n", results)
}

func RunSequenceParser() {
	sequenceParser := parc.SequenceOf(
		parc.Str("Hello "),
		parc.Str("World"),
	)
	results := sequenceParser.Parse("Hello World")
	fmt.Printf("%+v\n", results)
}
