package parc

import (
	"fmt"
	"regexp"
	"strings"
)

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

		inputLength := parserState.InputLength()
		if parserState.Index != inputLength {
			return updateParserError(
				parserState,
				fmt.Errorf("EndOfInput: expect end of input but index position is %d to the end", parserState.Index-inputLength))
		}
		return parserState
	}
	return NewParser("EndOfInput()", parserFun)
}

// Rest is a parser that returns the remaining input
func Rest() *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}

		inputLength := parserState.InputLength()
		if parserState.Index > inputLength {
			return updateParserError(
				parserState,
				fmt.Errorf("Rest: expect index %d less then or equal to the length of input %d", parserState.Index, inputLength))
		}
		return updateParserState(parserState, inputLength, Result(parserState.Remaining()))
	}
	return NewParser("Rest()", parserFun)
}

// Char is a parser that matches a fixed, single character value with the target string exactly one time
func Char(s string) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		if len(s) != 1 {
			return updateParserError(parserState, fmt.Errorf("Wrong argument for Char('%s'). It must be a single character", s))
		}

		if strings.HasPrefix(parserState.Remaining(), s) {
			return updateParserState(parserState, parserState.Index+len(s), Result(s))
		}

		return updateParserError(parserState, fmt.Errorf("Could not match '%s' with '%s'", s, parserState.Remaining()))
	}
	return NewParser("Char('"+s+"')", parserFun)
}

// Newline matches a space character ` `
var Space = Char(" ")

// Newline matches a newline character \n
var Newline = Char("\n")

// Tab matches a tab character \t
var Tab = Char("\t")

// Str is a parser that matches a fixed string value with the target string exactly one time
func Str(s string) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}

		slicedInput := parserState.Remaining()
		if len(slicedInput) == 0 {
			return updateParserError(parserState, fmt.Errorf("Str: tried to match '%s', but got Unexpected end of input.", s))
		}

		if strings.HasPrefix(slicedInput, s) {
			return updateParserState(parserState, parserState.Index+len(s), Result(s))
		}

		return updateParserError(parserState, fmt.Errorf("Str: could not match '%s' with '%s'", s, parserState.Remaining()))
	}
	return NewParser("Str('"+s+"')", parserFun)
}

// Crlf recognizes the string \r\n
var Crlf = Str("\r\n")

// RexExp is a parser that matches the regexpStr regular expression with the target string and returns with the first match.
// The patternName parameter defines a name for the expression for debugging purposes
func RegExp(patternName, regexpStr string) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}
		slicedInput := parserState.Remaining()
		if len(slicedInput) == 0 {
			return updateParserError(parserState, fmt.Errorf("RegExp: tried to match /%s/, but got Unexpected end of input.", regexpStr))
		}

		lettersRegexp := regexp.MustCompile(regexpStr)

		loc := lettersRegexp.FindIndex([]byte(slicedInput))

		if loc == nil {
			return updateParserError(parserState, fmt.Errorf("RegExp: could not match %s at index %d", patternName, parserState.Index))
		}

		return updateParserState(parserState, parserState.Index+loc[1], Result(slicedInput[loc[0]:loc[1]]))
	}
	return NewParser(fmt.Sprintf("RegExp('%s', /%s/)", patternName, regexpStr), parserFun)
}
