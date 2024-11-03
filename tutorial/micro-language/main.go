package main

import (
	"fmt"
	"github.com/tombenke/parc"
)

func main() {
	parc.Debug(0)
	input := "(+ (* 10 2) (- (/ 50 3) 2))"
	//input := "(+ 1 2)"
	fmt.Printf("%s => interpreter => %d\n", input, interpreter(input))
}

func interpreter(input string) int {
	parser := buildParser()
	parseResults := parser.Parse(&input)
	endResult := evaluate(parseResults.Results)
	return endResult
}

type Operation struct {
	Tag       string
	Operation string
	Operand_A parc.Result
	Operand_B parc.Result
}

type Operand struct {
	Tag   string
	Value int
}

func buildParser() parc.Parser {
	var expr parc.Parser
	var operation parc.Parser

	operator := parc.Choice(parc.Str("+"), parc.Str("-"), parc.Str("*"), parc.Str("/"))

	expr = *parc.Choice(
		parc.Map(parc.Integer, func(in parc.Result) parc.Result {
			operand := Operand{
				Tag:   "INTEGER",
				Value: in.(int),
			}
			return parc.Result(operand)
		}),
		&operation,
	)

	operation = *parc.Map(parc.SequenceOf(
		parc.Str("("),
		operator,
		parc.Str(" "),
		&expr,
		parc.Str(" "),
		&expr,
		parc.Str(")"),
	), func(in parc.Result) parc.Result {
		arr := in.([]parc.Result)
		op := Operation{
			Tag:       "OPERATION",
			Operation: arr[1].(string),
			Operand_A: arr[3],
			Operand_B: arr[5],
		}
		return parc.Result(op)
	})

	return expr
}

func evaluate(node parc.Result) int {
	////fmt.Printf("\nevaluateParseResults: %+v\n\n", node)
	switch n := node.(type) {
	case Operation:
		switch n.Operation {
		case "+":
			return evaluate(n.Operand_A) + evaluate(n.Operand_B)
		case "-":
			return evaluate(n.Operand_A) - evaluate(n.Operand_B)
		case "/":
			return evaluate(n.Operand_A) / evaluate(n.Operand_B)
		case "*":
			return evaluate(n.Operand_A) * evaluate(n.Operand_B)
		}
	case Operand:
		return n.Value
	}
	return 0
}
