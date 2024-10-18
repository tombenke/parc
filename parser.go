package parc

import (
	"fmt"
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