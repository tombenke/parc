package parc

import (
	"fmt"
	"unicode/utf8"
)

// ParserState represents an actual state of a parser
type ParserState struct {
	InputString string
	Results     Result
	Index       int
	Err         error
	IsError     bool
}

// NextRune returns the next rune in the input,
// as well as a new state in which the rune has been consumed.
func (ps ParserState) NextRune() (rune, ParserState) {
	r, w := utf8.DecodeRuneInString(ps.Remaining())
	return r, ps.Consume(w)
}

// Remaining returns the a string which is just the unconsumed input
func (ps ParserState) Remaining() string {
	return ps.InputString[ps.Index:]
}

// Consume returns a new state in which the index pointer is advanced by n bytes
func (ps ParserState) Consume(n int) ParserState {
	ps.Index += n
	return ps
}

// String returns with the string fromat of the parser state
func (ps ParserState) String() string {
	return fmt.Sprintf("InputString: '%s', Results: %+v, Index: %d, Err: %+v, IsError: %+v", ps.InputString, ps.Results, ps.Index, ps.Err, ps.IsError)
}

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