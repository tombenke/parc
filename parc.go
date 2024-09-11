package main

import (
	"fmt"
	"slices"
	"strings"
)

type ParserState struct {
	TargetString string
	Result       []string
	Index        int
	Err          error
	IsError      bool
}

// ParserFun type represents the generic format of parsers,
// that receives a ParserState as input,
// and returns with a new ParserState as an output
type ParserFun func(parserState ParserState) ParserState

// str is a parser that matches a fixed string value
func str(s string) ParserFun {
	return func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}

		if strings.HasPrefix(parserState.TargetString[parserState.Index:], s) {
			return updateParserState(parserState, parserState.Index+len(s), []string{s})
		}

		return updateParserError(parserState, fmt.Errorf("Could not match '%s' with '%s'", s, parserState.TargetString[parserState.Index:]))
	}
}

func sequenceOf(parsers ...ParserFun) ParserFun {
	return func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		results := []string{}
		nextState := parserState
		for _, parser := range parsers {
			nextState = parser(nextState)
			//results = slices.Concat(results, nextState.Result[len(nextState.Result)-1:])
		}
		return updateParserResult(nextState, results)
		//nextState.Result = results
		//return nextState
	}
}

// Returns with a new copy of state updated with the index and result values
func updateParserState(state ParserState, index int, result []string) ParserState {
	fmt.Printf("\n\n>> updateParserState(%+v, %+v, %+v)\n", state, index, result)
	newState := state
	newState.Index = index
	newState.Result = slices.Concat(state.Result, result)
	fmt.Printf("<< updateParserState(%+v, %+v, %+v) => %+v\n", state, index, result, newState)
	return newState
}

// Returns with a new copy of state updated with result values
func updateParserResult(state ParserState, result []string) ParserState {
	state.Result = slices.Concat(state.Result, result)
	return state
}

// updateParserError returns with a new copy of parser state within an error message
func updateParserError(state ParserState, errorMsg error) ParserState {
	state.IsError = true
	state.Err = errorMsg
	return state
}

// run runs a parser within an initial state on the target string
//func run(parser ParserFun, targetString string) ParserState {
//	initialState := ParserState{TargetString: targetString, Index: 0, Result: []string{}, Err: nil, IsError: false}
//	return parser(initialState)
//}

func main() {
	//helloStrParser := str("Hello")
	//fmt.Printf("%+v\n", run(helloStrParser, "Hello World"))
	sequenceParser := sequenceOf(
		str("Hello "),
		str("World"),
	)
	parser := NewParser(sequenceParser)
	//fmt.Printf("%+v\n", run(sequenceParser, "Hello World"))
	fmt.Printf("%+v\n", parser.Run("Hello World"))
}

type ParserStateTransformerFun ParserFun

type Parser struct {
	parserFun ParserFun
}

func NewParser(parserFun ParserFun) Parser {
	return Parser{parserFun: parserFun}
}

func (p Parser) Run(targetString string) ParserState {
	// run runs a parser within an initial state on the target string
	initialState := ParserState{TargetString: targetString, Index: 0, Result: []string{}, Err: nil, IsError: false}
	return p.parserFun(initialState)
}
