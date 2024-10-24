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
// the index (index) that parsing stopped at from this parsing session.
func (p *Parser) ErrorMap(mapperFn func(ParserState) error) *Parser {

	parserFun := func(parserState ParserState) ParserState {
		newState := p.ParserFun(parserState)
		if !newState.IsError {
			return newState
		}

		return updateParserError(newState, mapperFn(newState))
	}

	return NewParser("ErrorMap("+p.Name()+")", parserFun)
}
