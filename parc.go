package parc

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Result any

// Parser struct represent a parser
type Parser struct {
	ParserFun ParserFun
}

// NewParser is the constructor of the Parser
func NewParser(parserFun ParserFun) *Parser {
	return &Parser{ParserFun: parserFun}
}

// Parse runs the parser with the target string
func (p *Parser) Parse(inputString string) ParserState {
	// It runs a parser within an initial state on the target string
	initialState := ParserState{InputString: inputString, Index: 0, Results: Result(nil), Err: nil, IsError: false}
	return p.ParserFun(initialState)
}

// ParserState represents an actual state of a parser
type ParserState struct {
	InputString string
	Results     Result
	Index       int
	Err         error
	IsError     bool
}

func (ps ParserState) String() string {
	return fmt.Sprintf("InputString: '%s', Results: %+v, Index: %d, Err: %+v, IsError: %+v", ps.InputString, ps.Results, ps.Index, ps.Err, ps.IsError)
}

// ParserFun type represents the generic format of parsers,
// that receives a ParserState as input,
// and returns with a new ParserState as an output
type ParserFun func(parserState ParserState) ParserState

// Str is a parser that matches a fixed string value with the target string
func Str(s string) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		fmt.Printf("\n>> Str('%s', '%s')\n", s, parserState.InputString[parserState.Index:])
		newState := parserState
		if parserState.IsError {
			fmt.Printf("<< Str/newState: %+v\n", newState)
			return newState
		}

		if strings.HasPrefix(parserState.InputString[parserState.Index:], s) {
			newState = updateParserState(parserState, parserState.Index+len(s), Result(s))
			fmt.Printf("<< Str/newState: %+v\n", newState)
			return newState
		}

		newState = updateParserError(parserState, fmt.Errorf("Could not match '%s' with '%s'", s, parserState.InputString[parserState.Index:]))
		fmt.Printf("<< Str/newState: %+v\n", newState)
		return newState
	}
	parser := NewParser(parserFun)
	return parser
}

func Letters() *Parser {
	fmt.Printf("\n>> Letters()\n")
	return RegExp("^[A-Za-z]+", "letters")
}

func Digits() *Parser {
	fmt.Printf("\n>> Digits()\n")
	return RegExp("^[0-9]+", "digits")
}

func Integer() *Parser {
	fmt.Printf("\n>> Integer()\n")
	digitsToIntMapperFn := func(in Result) Result {
		fmt.Printf("\n>> Integer.digitsToIntMapperFn(%+v)\n", in)
		strValue := in.(string)
		intValue, _ := strconv.Atoi(strValue)
		return Result(intValue)
	}

	return Map(Digits(), digitsToIntMapperFn)
}

func RegExp(regexpStr, patternName string) *Parser {
	fmt.Printf("\n>> Regexp('%s', /%s/)\n", patternName, regexpStr)
	parserFun := func(parserState ParserState) ParserState {
		fmt.Printf("\n>> Regexp('%s', /%s/).parserFun('%s')\n", patternName, regexpStr, parserState.InputString[parserState.Index:])
		newState := parserState
		if parserState.IsError {
			fmt.Printf("<< RegExp('%s',/%s/).parserFun(newState: %+v\n", patternName, regexpStr, newState)
			return newState
		}
		slicedInputString := parserState.InputString[parserState.Index:]

		lettersRegexp := regexp.MustCompile(regexpStr)

		loc := lettersRegexp.FindIndex([]byte(slicedInputString))
		if len(loc) > 0 {
			fmt.Printf("   RegExp/%s in: '%s' found: '%s'\n", patternName, slicedInputString, slicedInputString[loc[0]:loc[1]])
		} else {
			fmt.Printf("   RegExp/%s/in: '%s' found: NONE\n", slicedInputString, patternName)
		}

		if loc == nil {
			newState = updateParserError(parserState, fmt.Errorf("Could not match %s at index %d", patternName, parserState.Index))
			fmt.Printf("<< RegExp/%s/newState: %+v\n", patternName, newState)
			return newState
		}

		newState = updateParserState(parserState, parserState.Index+loc[1], Result(slicedInputString[loc[0]:loc[1]]))
		fmt.Printf("<< RegExp/%s/newState: %+v\n", patternName, newState)
		return newState
	}
	parser := NewParser(parserFun)
	return parser
}

// Map call the map function to the result and returns with the return value of this function
func Map(parser *Parser, mapper func(Result) Result) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		fmt.Printf("\n>> Map('%s')\n", parserState.InputString[parserState.Index:])
		newState := parserState
		if parserState.IsError {
			fmt.Printf("<< Map/newState: %+v\n", newState)
			return newState
		}
		newState = parser.ParserFun(parserState)
		if newState.IsError {
			fmt.Printf("<< Map/newState: %+v\n", newState)
			return newState
		}

		result := mapper(newState.Results)

		newState = updateParserState(newState, newState.Index, Result(result))
		fmt.Printf("<< Map/newState: %+v\n", newState)
		return newState
	}
	mapParser := NewParser(parserFun)
	return mapParser
}

// // ErrorMap is like Map but it transforms the error value.
// // The function passed to ErrorMap gets an object the current error message (error),
// // the index (index) that parsing stopped at, and the data (data) from this parsing session.
// // Choice tries to execute the series of parser it got as a variadic parameter, and returns with the result of the first
// // parser that succeeds. If a parser returns with error, it tries to call the next one, until the last parsers.
// // In case none of them succeeds it does returns with error.
// func ErrorMap() {
// }

// SequenceOf is a parser that executes a sequence of parsers against a parser state
func SequenceOf(parsers ...*Parser) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		fmt.Printf("\n>> SequenceOf('%s')\n", parserState.InputString[parserState.Index:])
		if parserState.IsError {
			fmt.Printf("<< SequenceOf() => %+v\n", parserState)
			return parserState
		}
		results := make([]Result, 0, 10)
		nextState := parserState
		for _, parser := range parsers {
			nextState = (*parser).ParserFun(nextState)
			results = slices.Concat(results, []Result{Result(nextState.Results)})
		}
		newState := updateParserState(nextState, nextState.Index, Result(results))
		fmt.Printf("<< SequenceOf() => %+v\n", newState)
		return newState
	}
	parser := NewParser(parserFun)
	return parser
}

// Choice is a parser that executes a sequence of parsers against a parser state,
// and returns the first successful result if there is any
func Choice(parsers ...*Parser) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		fmt.Printf("\n>> Choice('%s')\n", parserState.InputString[parserState.Index:])
		if parserState.IsError {
			return parserState
		}
		var nextState ParserState
		for _, parser := range parsers {
			nextState = parser.ParserFun(parserState)
			if !nextState.IsError {
				fmt.Printf("<< Choice() => %+v\n", nextState)
				return nextState
			}
			fmt.Printf("   Choice try next with : %+v\n", parserState)
		}
		return updateParserError(parserState, fmt.Errorf("choice: Unable to match with any parser at index %d", parserState.Index))
	}
	parser := NewParser(parserFun)
	return parser
}

// // Many tries to execute the parser given as a parameter, until it succeeds. Aggregate the results and return with it at the end.
// // In never returns error either it could run the parser any times without errors or never.
// func Many() {
// }

// // ManyOne is similar to the Many parser, but it must be able to run the parser successfuly at least once, otherwise it return with error.
// func ManyOne() {
// }

// // Chain takes a function which recieves the last matched value and should return a parser.
// // That parser is then used to parse the following input, forming a chain of parsers based on previous input.
// // Chain is the fundamental way of creating contextual parsers.
// func Chain() {
// }

// Returns with a new copy of state updated with the index and result values
func updateParserState(state ParserState, index int, result Result) ParserState {
	newState := state
	newState.Index = index
	newState.Results = result
	//fmt.Printf("   updateParserState(%s, %+v, %+v)\n                  => %+v\n", state, index, result, newState)
	return newState
}

// updateParserError returns with a new copy of parser state within an error message
func updateParserError(state ParserState, errorMsg error) ParserState {
	newState := state
	newState.IsError = true
	newState.Err = errorMsg
	//fmt.Printf("   updateParserError(%s, %+v, %+v)\n                  => %+v\n", state, index, result, newState)
	return newState
}
