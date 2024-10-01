package parc

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type Result string

// Parser struct represent a parser
type Parser struct {
	ParserFun ParserFun
}

// NewParser is the consuctor of the Parser
func NewParser(parserFun ParserFun) Parser {
	return Parser{ParserFun: parserFun}
}

// Run runs the parser with the target string
func (p Parser) Run(inputString string) ParserState {
	// It runs a parser within an initial state on the target string
	initialState := ParserState{InputString: inputString, Index: 0, Results: []Result{}, Err: nil, IsError: false}
	return p.ParserFun(initialState)
}

// ParserState represents an actual state of a parser
type ParserState struct {
	InputString string
	Results     []Result
	Index       int
	Err         error
	IsError     bool
}

// ParserFun type represents the generic format of parsers,
// that receives a ParserState as input,
// and returns with a new ParserState as an output
type ParserFun func(parserState ParserState) ParserState

// Str is a parser that matches a fixed string value with the target string
func Str(s string) Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}

		if strings.HasPrefix(parserState.InputString[parserState.Index:], s) {
			return updateParserState(parserState, parserState.Index+len(s), []Result{Result(s)})
		}

		return updateParserError(parserState, fmt.Errorf("Could not match '%s' with '%s'", s, parserState.InputString[parserState.Index:]))
	}
	parser := NewParser(parserFun)
	return parser
}

func Letters() Parser {
	return RegExp("^[A-Za-z]+", "letters")
}

func Digits() Parser {
	return RegExp("^[0-9]+", "digits")
}

func RegExp(regexpStr, patternName string) Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		slicedInputString := parserState.InputString[parserState.Index:]

		lettersRegexp := regexp.MustCompile(regexpStr)

		loc := lettersRegexp.FindIndex([]byte(slicedInputString))
		fmt.Printf("Letters %s => %+v\n", slicedInputString, loc)

		if loc == nil {
			return updateParserError(parserState, fmt.Errorf("Could not match %s at index %d", patternName, parserState.Index))
		}

		return updateParserState(parserState, parserState.Index+loc[1], []Result{Result(slicedInputString[loc[0]:loc[1]])})
	}
	parser := NewParser(parserFun)
	return parser
}

// Map call the map function to the result and returns with the return value of this function
func Map() {
}

// ErrorMap is like Map but it transforms the error value.
// The function passed to ErrorMap gets an object the current error message (error),
// the index (index) that parsing stopped at, and the data (data) from this parsing session.
// Choice tries to execute the series of parser it got as a variadic parameter, and returns with the result of the first
// parser that succeeds. If a parser returns with error, it tries to call the next one, until the last parsers.
// In case none of them succeeds it does returns with error.
func ErrorMap() {
}

// SequenceOf is a parser that executes a sequence of parsers against a parser state
func SequenceOf(parsers ...Parser) Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		////results := []Result{}
		nextState := parserState
		for _, parser := range parsers {
			nextState = parser.ParserFun(nextState)
			////results = slices.Concat(results, nextState.Results[len(nextState.Results)-1:])
		}
		return nextState ////updateParserResults(nextState, results)
	}
	parser := NewParser(parserFun)
	return parser
}

func Choice() {
}

// Many tries to execute the parser given as a parameter, until it succeeds. Aggregate the results and return with it at the end.
// In never returns error either it could run the parser any times without errors or never.
func Many() {
}

// ManyOne is similar to the Many parser, but it must be able to run the parser successfuly at least once, otherwise it return with error.
func ManyOne() {
}

// Chain takes a function which recieves the last matched value and should return a parser.
// That parser is then used to parse the following input, forming a chain of parsers based on previous input.
// Chain is the fundamental way of creating contextual parsers.
func Chain() {
}

// Returns with a new copy of state updated with the index and result values
func updateParserState(state ParserState, index int, result []Result) ParserState {
	fmt.Printf("\n\n>> updateParserState(%+v, %+v, %+v)\n", state, index, result)
	newState := state
	newState.Index = index
	newState.Results = slices.Concat(state.Results, result)
	fmt.Printf("<< updateParserState(%+v, %+v, %+v) => %+v\n", state, index, result, newState)
	return newState
}

//// Returns with a new copy of state updated with result values
//func updateParserResults(state ParserState, result []Result) ParserState {
//	state.Results = slices.Concat(state.Results, result)
//	return state
//}

// updateParserError returns with a new copy of parser state within an error message
func updateParserError(state ParserState, errorMsg error) ParserState {
	state.IsError = true
	state.Err = errorMsg
	return state
}
