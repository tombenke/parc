package parc

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Result any

var (
	parseDepth = 0
	debugLevel = 0
)

// Parser struct represent a parser
type Parser struct {
	name      string
	ParserFun ParserFun
}

// Debug switches debugging ON with the given level. Level=0 means, Debug is switched off.
func Debug(level int) {
	debugLevel = level
}

// NewParser is the constructor of the Parser
func NewParser(parserName string, parserFun ParserFun) *Parser {
	wrapperFn := func(parserState ParserState) ParserState {
		var indent string
		if debugLevel > 0 {
			indent = strings.Repeat("|   ", parseDepth)
			fmt.Printf("%s+-> %s <= Input: '%s'\n", indent, parserName, parserState.InputString[parserState.Index:])
			parseDepth = parseDepth + 1
		}
		newState := parserFun(parserState)
		if debugLevel > 0 {
			parseDepth = parseDepth - 1
			fmt.Printf("%s+<- %s =>\n", indent, parserName)
			if debugLevel > 1 {
				fmt.Printf("%s    Err: %+v, Result: '%+v'\n", indent, newState.Err, newState.Results)
			} else {
				fmt.Printf("%s    Err: %+v\n", indent, newState.Err)
			}
		}
		return newState
	}
	return &Parser{name: parserName, ParserFun: wrapperFn}
}

// Parse runs the parser with the target string
func (p *Parser) Parse(inputString string) ParserState {
	// It runs a parser within an initial state on the target string
	initialState := ParserState{InputString: inputString, Index: 0, Results: Result(nil), Err: nil, IsError: false}
	return p.ParserFun(initialState)
}

// Name returns the name of the parser
func (p *Parser) Name() string {
	return p.name
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
		newState := parserState
		if parserState.IsError {
			return newState
		}

		if strings.HasPrefix(parserState.InputString[parserState.Index:], s) {
			newState = updateParserState(parserState, parserState.Index+len(s), Result(s))
			return newState
		}

		newState = updateParserError(parserState, fmt.Errorf("Could not match '%s' with '%s'", s, parserState.InputString[parserState.Index:]))
		return newState
	}
	parser := NewParser("Str('"+s+"')", parserFun)
	return parser
}

func Letters() *Parser {
	return RegExp("Letters", "^[A-Za-z]+")
}

func Digits() *Parser {
	return RegExp("Digits", "^[0-9]+")
}

func Integer() *Parser {
	digitsToIntMapperFn := func(in Result) Result {
		strValue := in.(string)
		intValue, _ := strconv.Atoi(strValue)
		return Result(intValue)
	}

	return Map(Digits(), digitsToIntMapperFn)
}

func RegExp(patternName, regexpStr string) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		newState := parserState
		if parserState.IsError {
			return newState
		}
		slicedInputString := parserState.InputString[parserState.Index:]

		lettersRegexp := regexp.MustCompile(regexpStr)

		loc := lettersRegexp.FindIndex([]byte(slicedInputString))

		if loc == nil {
			newState = updateParserError(parserState, fmt.Errorf("Could not match %s at index %d", patternName, parserState.Index))
			return newState
		}

		newState = updateParserState(parserState, parserState.Index+loc[1], Result(slicedInputString[loc[0]:loc[1]]))
		return newState
	}
	parser := NewParser(fmt.Sprintf("Regexp('%s', /%s/)", patternName, regexpStr), parserFun)

	return parser
}

// Map call the map function to the result and returns with the return value of this function
func Map(parser *Parser, mapper func(Result) Result) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		newState := parserState
		if parserState.IsError {
			return newState
		}
		newState = parser.ParserFun(parserState)
		if newState.IsError {
			return newState
		}

		result := mapper(newState.Results)

		newState = updateParserState(newState, newState.Index, Result(result))
		return newState
	}

	mapParser := NewParser("Map("+parser.Name()+")", parserFun)
	return mapParser
}

// ErrorMap is like Map but it transforms the error value.
// The function passed to ErrorMap gets an object the current error message (error),
// the index (index) that parsing stopped at, and the data (data) from this parsing session.
// Choice tries to execute the series of parser it got as a variadic parameter, and returns with the result of the first
// parser that succeeds. If a parser returns with error, it tries to call the next one, until the last parsers.
// In case none of them succeeds it does returns with error.
// TODO
// func ErrorMap() {
// }

// SequenceOf is a parser that executes a sequence of parsers against a parser state
func SequenceOf(parsers ...*Parser) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		results := make([]Result, 0, 10)
		nextState := parserState
		for _, parser := range parsers {
			nextState = (*parser).ParserFun(nextState)
			results = slices.Concat(results, []Result{Result(nextState.Results)})
		}
		newState := updateParserState(nextState, nextState.Index, Result(results))
		return newState
	}
	parser := NewParser("SequenceOf("+getParserNames(parsers...)+")", parserFun)
	return parser
}

// getParserNames returns a string of the comma separated list of parser names
func getParserNames(parsers ...*Parser) string {
	parserNames := ""

	if debugLevel > 1 {
		for i, parser := range parsers {
			parserNames = parserNames + parser.Name()
			if i < len(parsers)-1 {
				parserNames = parserNames + ", "
			}
		}
	}
	return parserNames
}

// Choice is a parser that executes a sequence of parsers against a parser state,
// and returns the first successful result if there is any
func Choice(parsers ...*Parser) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		var nextState ParserState
		for _, parser := range parsers {
			nextState = parser.ParserFun(parserState)
			if !nextState.IsError {
				return nextState
			}
		}
		return updateParserError(parserState, fmt.Errorf("choice: Unable to match with any parser at index %d", parserState.Index))
	}
	parser := NewParser("Choice("+getParserNames(parsers...)+")", parserFun)
	return parser
}

// Many tries to execute the parser given as a parameter, until it succeeds. Aggregate the results and return with it at the end.
// In never returns error either it could run the parser any times without errors or never.
// TODO
// func Many() {
// }

// ManyOne is similar to the Many parser, but it must be able to run the parser successfuly at least once,
// otherwise it return with error.
// TODO
// func ManyOne() {
// }

// Chain takes a function which receieves the last matched value and should return a parser.
// That parser is then used to parse the following input, forming a chain of parsers based on previous input.
// Chain is the fundamental way of creating contextual parsers.
// TODO
// func Chain() {
// }

// TODO
// func EndOfInput *Parser {}

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
