package main

import (
	"fmt"
	"github.com/tombenke/parc"
)

func main() {
	input := "Hello World"
	resultState := parc.RegExp("^[A-Za-z]{5} [A-Za-z]{5}$").As("HelloWorld").Parse(&input)
	fmt.Printf("\n%+v\n", resultState)

	// => inputString: 'Hello World', Results: Hello World, Index: 11, Err: <nil>, IsError: false
}
