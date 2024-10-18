package parc

import ()

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
