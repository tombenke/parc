package parc

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Result represents the type of the result that is produced by calling the parser function of a parser.
// It is stored in the parser state the parser's parser function returns with.
type Result any

var (
	// parseDepth defines the actual call-depth of a specific parser during the parsing
	parseDepth = 0

	// debugLevel sets the actual debug-level. 0=NO-DEBUG, 1=minimum, 2=medium, 3=detailed
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

// Name returns the name of the parser
func (p *Parser) Name() string {
	return p.name
}

// Parse runs the parser with the target string
func (p *Parser) Parse(inputString string) ParserState {
	// It runs a parser within an initial state on the target string
	initialState := ParserState{InputString: inputString, Index: 0, Results: Result(nil), Err: nil, IsError: false}
	return p.ParserFun(initialState)
}

// Map call the map function to the result and returns with the return value of this function
func (p *Parser) Map(mapper func(Result) Result) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		newState := p.ParserFun(parserState)
		if newState.IsError {
			return newState
		}

		result := mapper(newState.Results)

		return updateParserState(newState, newState.Index, Result(result))
	}

	return NewParser("Map("+p.Name()+")", parserFun)
}

// Chain takes a function which receieves the last matched value and should return a parser.
// That parser is then used to parse the following input, forming a chain of parsers based on previous input.
// Chain is the fundamental way of creating contextual parsers.
func (p *Parser) Chain(parserMakerFn func(Result) *Parser) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		newState := p.ParserFun(parserState)
		if newState.IsError {
			return newState
		}

		nextParser := parserMakerFn(newState.Results)
		result := nextParser.ParserFun(newState)

		return updateParserState(newState, newState.Index, Result(result))
	}

	return NewParser("Chain("+p.Name()+")", parserFun)
}

// Map call the map function to the result and returns with the return value of this function
func Map(parser *Parser, mapper func(Result) Result) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		newState := parser.ParserFun(parserState)
		if newState.IsError {
			return newState
		}

		result := mapper(newState.Results)

		return updateParserState(newState, newState.Index, Result(result))
	}

	return NewParser("Map("+parser.Name()+")", parserFun)
}

// Chain takes a function which receieves the last matched value and should return a parser.
// That parser is then used to parse the following input, forming a chain of parsers based on previous input.
// Chain is the fundamental way of creating contextual parsers.
func Chain(parser *Parser, parserMakerFn func(Result) *Parser) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		newState := parser.ParserFun(parserState)
		if newState.IsError {
			return newState
		}

		nextParser := parserMakerFn(newState.Results)
		result := nextParser.ParserFun(newState)

		return updateParserState(newState, newState.Index, Result(result))
	}

	return NewParser("Chain("+parser.Name()+")", parserFun)
}

// ErrorMap is like Map but it transforms the error value.
// The function passed to ErrorMap gets an object the current error message (error),
// the index (index) that parsing stopped at, and the data (data) from this parsing session.
// Choice tries to execute the series of parser it got as a variadic parameter, and returns with the result of the first
// parser that succeeds. If a parser returns with error, it tries to call the next one, until the last parsers.
// In case none of them succeeds it does returns with error.
func (p *Parser) ErrorMap(mapperFn func(err error, index int) error) *Parser {

	parserFun := func(parserState ParserState) ParserState {
		newState := p.ParserFun(parserState)
		if !newState.IsError {
			return newState
		}

		return updateParserError(newState, mapperFn(newState.Err, newState.Index))
	}

	return NewParser("ErrorMap("+p.Name()+")", parserFun)
}

// ParserState represents an actual state of a parser
type ParserState struct {
	InputString string
	Results     Result
	Index       int
	Err         error
	IsError     bool
}

// String returns with the string fromat of the parser state
func (ps ParserState) String() string {
	return fmt.Sprintf("InputString: '%s', Results: %+v, Index: %d, Err: %+v, IsError: %+v", ps.InputString, ps.Results, ps.Index, ps.Err, ps.IsError)
}

// ParserFun type represents the generic format of parsers,
// that receives a ParserState as input,
// and returns with a new ParserState as an output
type ParserFun func(parserState ParserState) ParserState

// Char is a parser that matches a fixed, single character value with the target string exactly one time
func Char(s string) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		if len(s) != 1 {
			return updateParserError(parserState, fmt.Errorf("Wrong argument for Char('%s'). It must be a single character", s))
		}

		if strings.HasPrefix(parserState.InputString[parserState.Index:], s) {
			return updateParserState(parserState, parserState.Index+len(s), Result(s))
		}

		return updateParserError(parserState, fmt.Errorf("Could not match '%s' with '%s'", s, parserState.InputString[parserState.Index:]))
	}
	return NewParser("Char('"+s+"')", parserFun)
}

// StartOfInput is a parser that only succeeds when the parser is at the beginning of the input.
func StartOfInput() *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}

		if parserState.Index > 0 {
			return updateParserError(
				parserState,
				fmt.Errorf("StartOfInput: expect start of input but index position is %d", parserState.Index))
		}
		return parserState
	}
	return NewParser("StartOfInput()", parserFun)
}

// EndOfInput is a parser that only succeeds when there is no more input to be parsed.
func EndOfInput() *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}

		inputLength := len(parserState.InputString)
		if parserState.Index != inputLength {
			return updateParserError(
				parserState,
				fmt.Errorf("EndOfInput: expect end of input but index position is %d to the end", parserState.Index-inputLength))
		}
		return parserState
	}
	return NewParser("EndOfInput()", parserFun)
}

// Str is a parser that matches a fixed string value with the target string exactly one time
func Str(s string) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}

		if len(parserState.InputString[parserState.Index:]) == 0 {
			return updateParserError(parserState, fmt.Errorf("Str: tried to match '%s', but got Unexpected end of input.", s))
		}

		if strings.HasPrefix(parserState.InputString[parserState.Index:], s) {
			return updateParserState(parserState, parserState.Index+len(s), Result(s))
		}

		return updateParserError(parserState, fmt.Errorf("Str: could not match '%s' with '%s'", s, parserState.InputString[parserState.Index:]))
	}
	return NewParser("Str('"+s+"')", parserFun)
}

// Letters is a parser that matches a single letter character with the target string
func Letter() *Parser {
	return RegExp("Letter", "^[A-Za-z]")
}

// Letters is a parser that matches one or more letter characters with the target string
func Letters() *Parser {
	return RegExp("Letters", "^[A-Za-z]+")
}

// Digit is a parser that matches a singl digit character with the target string
func Digit() *Parser {
	return RegExp("Digit", "^[0-9]")
}

// Digits is a parser that matches one or more digit characters with the target string
func Digits() *Parser {
	return RegExp("Digits", "^[0-9]+")
}

// Integer is a parser that matches one or more digit characters with the target string and returns with an int value
func Integer() *Parser {
	digitsToIntMapperFn := func(in Result) Result {
		strValue := in.(string)
		intValue, _ := strconv.Atoi(strValue)
		return Result(intValue)
	}

	return Digits().Map(digitsToIntMapperFn)
}

// RexExp is a parser that matches the regexpStr regular expression with the target string and returns with the first match.
// The patternName parameter defines a name for the expression for debugging purposes
func RegExp(patternName, regexpStr string) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		slicedInputString := parserState.InputString[parserState.Index:]
		if len(slicedInputString) == 0 {
			return updateParserError(parserState, fmt.Errorf("RegExp: tried to match /%s/, but got Unexpected end of input.", regexpStr))
		}

		lettersRegexp := regexp.MustCompile(regexpStr)

		loc := lettersRegexp.FindIndex([]byte(slicedInputString))

		if loc == nil {
			return updateParserError(parserState, fmt.Errorf("RegExp: could not match %s at index %d", patternName, parserState.Index))
		}

		return updateParserState(parserState, parserState.Index+loc[1], Result(slicedInputString[loc[0]:loc[1]]))
	}
	return NewParser(fmt.Sprintf("RegExp('%s', /%s/)", patternName, regexpStr), parserFun)
}

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
		return updateParserState(nextState, nextState.Index, Result(results))
	}
	return NewParser("SequenceOf("+getParserNames(parsers...)+")", parserFun)
}

// ZeroOrMore tries to execute the parser given as a parameter, until it succeeds.
// Aggregate the results and returns with it at the end.
// It never returns error either it could run the parser any times without errors or never.
func ZeroOrMore(parser *Parser) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}

		results := make([]Result, 0, 10)
		nextState := parserState

		for {
			testState := parser.ParserFun(nextState)
			if testState.IsError {
				break
			} else {
				results = slices.Concat(results, []Result{Result(testState.Results)})
				nextState = testState
			}
		}
		return updateParserState(nextState, nextState.Index, Result(results))
	}
	return NewParser("ZeroOrMore("+parser.Name()+")", parserFun)
}

// OneOrMore is similar to the ZeroOrMore parser,
// but it must be able to run the parser successfuly at least once, otherwise it return with error.
// It executes the parser given as a parameter, until it succeeds,
// meanwhile it aggregate the results then returns with it at the end.
func OneOrMore(parser *Parser) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}

		results := make([]Result, 0, 10)
		nextState := parserState
		fmt.Printf("\n=== nextState: %+v\n", nextState)

		for {
			testState := parser.ParserFun(nextState)
			fmt.Printf("\n=== testState: %+v\n", testState)
			if testState.IsError {
				break
			} else {
				results = slices.Concat(results, []Result{Result(testState.Results)})
				nextState = testState
			}
		}
		if len(results) == 0 {
			return updateParserError(parserState, fmt.Errorf("ZeroOrMore: unable to match any input using parser at index %d", parserState.Index))
		}
		return updateParserState(nextState, nextState.Index, Result(results))
	}
	return NewParser("ZeroOrMore("+parser.Name()+")", parserFun)
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
	return NewParser("Choice("+getParserNames(parsers...)+")", parserFun)
}

// Between is a utility function that takes two parsers as arguments that defines a starting and ending pattern of a content,
// and returns a function that takes a content parser as argument.
// Using the resulted parser will provide a result that is the outcome of the content parser.
func Between(leftParser, rightParser *Parser) func(*Parser) *Parser {
	return func(contentParser *Parser) *Parser {
		return SequenceOf(
			leftParser,
			contentParser,
			rightParser,
		).Map(func(result Result) Result {
			arrResults := result.([]Result)
			return arrResults[1]
		})
	}
}

// TODO
// func StartOfInput() *Parser {}
// func EndOfInput() *Parser {}

// Template is a sample code block to create a new parser
//func Template() *Parser {
//	parserFun := func(parserState ParserState) ParserState {
//		if parserState.IsError {
//			return parserState
//		}
//		// TODO: Add logic here
//		if OK {
//			return updateParserState(parserState, newIndex, Result(theResult))
//		}
//
//		return updateParserError(parserState, fmt.Errorf("error ... at %d with input: %s", parserState.Index, parserState.InputString[parserState.Index:]))
//	}
//	return NewParser("Template(...)", parserFun)
//}

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
