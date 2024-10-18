package parc

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

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
