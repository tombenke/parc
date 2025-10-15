package errors

import (
	"github.com/tombenke/parc"
)

type StackData float64

type Atom struct {
	Tag   string
	Value parc.Result
}

type Term struct {
	Tag       string
	Operand_A parc.Result
	Operand_B parc.Result
	Operator  Operator
}

type Expression struct {
	Tag       string
	Operand_A parc.Result
	Operand_B parc.Result
	Operator  Operator
}

type Number struct {
	Tag   string
	Value StackData
}

type Operator struct {
	Tag   string
	Value string
}

type Constant struct {
	Tag  string
	Name string
}

// The parser instance
var parser parc.Parser

func init() {
	parc.Debug(0)
	parser = buildParser()
}

// Parse uses the internal parser object that parses the input string and returns with the results in the form of AST tree.
func Parse(source string) (parc.Result, error) {

	parseResults := parser.Parse(&source)
	if parseResults.IsError {
		return nil, parseResults.Err
	}

	return parseResults.Results, nil
}

// buildParser creates a parser of the language. It is called by the init function.
func buildParser() parc.Parser {
	var expression parc.Parser

	constant := *parc.Map(parc.Choice(parc.Str("pi"), parc.Str("phi"), parc.Str("e")).As("constant"),
		func(in parc.Result) parc.Result {
			return parc.Result(
				Constant{
					Tag:  "CONSTANT",
					Name: in.(string),
				},
			)
		})

	intNumber := *parc.Map(parc.Integer.As("intNumber"), func(in parc.Result) parc.Result {
		operand := Number{
			Tag:   "NUMBER",
			Value: StackData(in.(int)),
		}
		return parc.Result(operand)
	})

	realNumber := *parc.Map(parc.RealNumber.As("realNumber"), func(in parc.Result) parc.Result {
		operand := Number{
			Tag:   "NUMBER",
			Value: StackData(in.(float64)),
		}
		return parc.Result(operand)
	})

	number := *parc.Choice(&realNumber, &intNumber).As("number")

	mulOperator := *parc.Map(parc.Choice(parc.Str("*"), parc.Str("/")).As("mulOperator"),
		func(in parc.Result) parc.Result {
			return parc.Result(Operator{
				Tag:   "OPERATOR",
				Value: in.(string),
			})
		})

	addOperator := *parc.Map(parc.Choice(parc.Str("+"), parc.Str("-")).As("addOperator"),
		func(in parc.Result) parc.Result {
			return parc.Result(Operator{
				Tag:   "OPERATOR",
				Value: in.(string),
			})
		})

	bracketed_expression := *parc.Map(parc.SequenceOf(
		parc.Str("("),
		parc.ZeroOrMore(parc.Cond(parc.IsWhitespace)),
		&expression,
		parc.ZeroOrMore(parc.Cond(parc.IsWhitespace)),
		parc.Str(")"),
	).As("bracketed_expression"), func(in parc.Result) parc.Result {
		arrResult := in.([]parc.Result)
		return parc.Result(arrResult[2])
	})

	atom := *parc.Choice(&number, &constant, &bracketed_expression).As("atom")

	term := *parc.Map(parc.SequenceOf(
		&atom,
		parc.ZeroOrMore(
			parc.SequenceOf(
				parc.ZeroOrMore(parc.Cond(parc.IsWhitespace)),
				&mulOperator,
				parc.ZeroOrMore(parc.Cond(parc.IsWhitespace)),
				&atom,
			),
		),
	).As("term"), func(in parc.Result) parc.Result {
		arr := in.([]parc.Result)
		arr1 := arr[1].([]parc.Result)
		if len(arr1) == 0 {
			return parc.Result(arr[0])
		}
		arr11 := arr1[0].([]parc.Result)
		return parc.Result(Term{Tag: "TERM", Operand_A: arr[0], Operator: arr11[1].(Operator), Operand_B: arr11[3]})
	})

	expression = *parc.Map(parc.SequenceOf(
		&term,
		parc.ZeroOrMore(
			parc.SequenceOf(
				parc.ZeroOrMore(parc.Cond(parc.IsWhitespace)),
				&addOperator,
				parc.ZeroOrMore(parc.Cond(parc.IsWhitespace)),
				&term,
			),
		),
	).As("expression"), func(in parc.Result) parc.Result {
		arr := in.([]parc.Result)
		arr1 := arr[1].([]parc.Result)
		if len(arr1) == 0 {
			return parc.Result(arr[0])
		}
		arr11 := arr1[0].([]parc.Result)
		return parc.Result(Expression{Tag: "EXPRESSION", Operand_A: arr[1], Operator: arr11[1].(Operator), Operand_B: arr11[3]})
	})

	formula := *parc.Map(parc.SequenceOf(
		parc.StartOfInput(),
		&expression,
		parc.EndOfInput(),
	).As("formula"), func(in parc.Result) parc.Result {
		arrResult := in.([]parc.Result)
		return parc.Result(arrResult[1])
	})
	return formula
}
