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
	results := parc.Str("Hello").Run("Hello World")
	fmt.Printf("%+v\n", results)
}

func RunSequenceParser() {
	sequenceParser := parc.SequenceOf(
		parc.Str("Hello "),
		parc.Str("World"),
	)
	results := sequenceParser.Run("Hello World")
	fmt.Printf("%+v\n", results)
}
