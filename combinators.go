package parc

import (
	"fmt"
	"slices"
)

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
			if nextState.IsError {
				return updateParserError(parserState, fmt.Errorf("SequenceOf: unable to match the `%s` parser at index %d with '%s'", parser.Name(), parserState.Index, parserState.Remaining()))
			}
			results = slices.Concat(results, []Result{Result(nextState.Results)})
		}
		return updateParserState(nextState, nextState.Index, Result(results))
	}
	return NewParser("SequenceOf("+getParserNames(parsers...)+")", parserFun)
}

// Times is an alias of the Count parser
var Times = Count

// Count tries to execute the parser given as a parameter exactly count times.
// Collects results into an array and returns with it at the end.
// It returns error if it could not run the parser exaclty count times.
// You can use Times parser, instead of Count since that is an alias of this parser.
func Count(parser *Parser, count int) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		results := make([]Result, 0, 10)
		nextState := parserState

		for {
			testState := parser.ParserFun(nextState)
			if testState.IsError || len(results) >= count {
				break
			} else {
				results = slices.Concat(results, []Result{Result(testState.Results)})
				nextState = testState
			}
		}
		if len(results) != count {
			return updateParserError(parserState, fmt.Errorf("Count: unable to match the parser exactly %d times at index %d with '%s'", count, parserState.Index, parserState.Remaining()))
		}
		return updateParserState(nextState, nextState.Index, Result(results))
	}
	return NewParser("Count("+parser.Name()+")", parserFun)
}

// TimesMin is an alias of the CountMin parser
var TimesMin = CountMin

// CountMin tries to execute the parser given as a parameter at least minOccurences times.
// Collects results into an array and returns with it at the end.
// It returns error if it could not run the parser at least minOccurences times.
// You can use TimesMin parser, instead of CountMin since that is an alias of this parser.
func CountMin(parser *Parser, minOccurences int) *Parser {
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
		if len(results) < minOccurences {
			return updateParserError(parserState, fmt.Errorf("CountMin: unable to match the parser at least %d times at index %d with '%s'", minOccurences, parserState.Index, parserState.Remaining()))
		}
		return updateParserState(nextState, nextState.Index, Result(results))
	}
	return NewParser("CountMin("+parser.Name()+")", parserFun)
}

// TimesMinMax is an alias of the CountMinMax parser
var TimesMinMax = CountMinMax

// CountMinMax tries to execute the parser given as a parameter at least minOccurences but maximum maxOccurences times.
// Collects results into an array and returns with it at the end.
// It returns error if it could not run the parser at least minOccurences times.
// You can use TimesMinMax parser, instead of CountMinMax since that is an alias of this parser.
func CountMinMax(parser *Parser, minOccurences int, maxOccurences int) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		results := make([]Result, 0, 10)
		nextState := parserState

		for {
			testState := parser.ParserFun(nextState)
			if testState.IsError || len(results) >= maxOccurences {
				break
			} else {
				results = slices.Concat(results, []Result{Result(testState.Results)})
				nextState = testState
			}
		}
		if len(results) < minOccurences {
			return updateParserError(parserState, fmt.Errorf("CountMinMax: unable to match the parser at least %d times at index %d", minOccurences, parserState.Index))
		}
		return updateParserState(nextState, nextState.Index, Result(results))
	}
	return NewParser("CountMinMax("+parser.Name()+")", parserFun)
}

// ZeroOrOne tries to execute the parser given as a parameter once.
// It returns `nil` if it could not match, or a single result if match occured.
// It never returns error either it could run the parser only once or could not run it at all.
func ZeroOrOne(parser *Parser) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}

		nextState := parser.ParserFun(parserState)
		if nextState.IsError {
			return updateParserState(parserState, nextState.Index, Result(nil))
		}

		return nextState
	}
	return NewParser("ZeroOne("+parser.Name()+")", parserFun)
}

// ZeroOrMore tries to execute the parser given as a parameter, until it succeeds.
// Collects the results into an array and returns with it at the end.
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
// meanwhile it collects the results into an array then returns with it at the end.
func OneOrMore(parser *Parser) *Parser {
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
