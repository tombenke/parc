package main

import (
	"fmt"
	"github.com/tombenke/parc"
)

func main() {
	input := "Hello World"
	resultState := parc.Cond(parc.IsAsciiLetter).Parse(&input)

	fmt.Printf("%+v\n", resultState)

	// => inputString: 'Hello World', Results: H, Index: 1, Err: <nil>, IsError: false
}
