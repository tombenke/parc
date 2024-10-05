package main

import (
	"fmt"
	"github.com/tombenke/parc"
)

func main() {
	input := "(+ (* 10 2) (- (/ 50 3) 2))"
	//input := "(+ 1 2)"
	UseMicroLanguage(input)
}

func UseMicroLanguage(input string) {
	var expr parc.Parser

	operator := parc.Choice(parc.Str("+"), parc.Str("-"), parc.Str("*"), parc.Str("/"))

	lit1 := parc.Str("(")
	lit2 := parc.Str(" ")
	lit3 := parc.Str(" ")
	lit4 := parc.Str(")")

	operation := parc.Map(parc.RefSequenceOf(
		&lit1, // parc.Str("("),
		&operator,
		&lit2, // parc.Str(" "),
		&expr,
		&lit3, // parc.Str(" "),
		&expr,
		&lit4, // parc.Str(")"),
	), func(in parc.Result) parc.Result {
		return parc.Result(in)
	})

	expr = parc.Choice(
		parc.Integer(),
		operation,
	)

	results := expr.Parse(input)
	fmt.Printf("\nresults: %+v\n\n", results)
}
