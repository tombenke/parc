package parc

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

// IsAlphabetic is an alias of IsAsciiLetter
var IsAlphabetic = IsAsciiLetter

// IsAlphabetic tests if rune is ASCII alphabetic: [A-Za-z]
func IsAsciiLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

// IsDigit is an alias of IsDecimalDigit
var IsDigit = IsDecimalDigit

// IsDecimalDigit tests if rune is ASCII digit: [0-9]
func IsDecimalDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// IsHexadecimalDigit tests if rune is ASCII hex digit: [0-9A-Fa-f]
func IsHexadecimalDigit(r rune) bool {
	return r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F'
}

// IsOctDigit tests if rune is ASCII octal digit: [0-7]
func IsOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

// IsBinaryDigit tests if rune is ASCII binary digit: [0-1]
func IsBinaryDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

// IsAlphaNumeric tests if rune is ASCII alphanumeric: [A-Za-z0-9]
func IsAlphaNumeric(r rune) bool {
	return IsAsciiLetter(r) || IsDecimalDigit(r)
}

// IsWhitespace tests if rune is ASCII space, newline or tab
func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

// IsSpace tests if rune is ASCII space or tab
func IsSpace(r rune) bool {
	return r == ' '
}

// IsTab Tests if rune is ASCII space or tab: [\t]
func IsTab(r rune) bool {
	return r == '\t'
}

// IsNewline tests if rune is ASCII newline: [\n]
func IsNewline(r rune) bool {
	return r == '\n'
}

// IsCarriageReturn tests if rune is ASCII newline: [\r]
func IsCarriageReturn(r rune) bool {
	return r == '\r'
}

// IsAnyChar tests if rune is any char. Actually it always returns with true.
func IsAnyChar(r rune) bool {
	return true
}

// Cond returns a Parser which tests the next rune in the input with the condition function.
// If the condition is met, the rune is consumed from the input and the parser succeeds.
// Otherwise the parser fails.
func Cond(conditionFn func(rune) bool) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}

		if parserState.AtTheEnd() {
			return updateParserError(parserState, fmt.Errorf("Cond: got Unexpected end of input."))
		}

		// Try to take a single occurence
		r, nextState := parserState.NextRune()
		if !conditionFn(r) {
			return updateParserError(parserState, fmt.Errorf("Cond: could not match %c at index %d", r, parserState.Index))
		}
		return updateParserState(parserState, nextState.Index, Result(string(r)))
	}
	return NewParser("Cond()", parserFun)
}

// CondMin returns a Parser which tests the next rune in the input with the condition function.
// If the condition is met, the rune is consumed from the input and the parser succeeds as many times as possible,
// but at least `minOccurences` times.
// Otherwise the parser fails.
func CondMin(conditionFn func(rune) bool, minOccurences int) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}

		if parserState.AtTheEnd() && minOccurences > 0 {
			return updateParserError(parserState, fmt.Errorf("CondMin: got Unexpected end of input."))
		}

		if minOccurences < 0 {
			return updateParserError(parserState, fmt.Errorf("CondMin: wrong minOccurences value %d", minOccurences))
		}

		currentState := parserState
		numFound := 0
		results := make([]byte, 0, 10)

		// Try to take as many occurences as possible, but at least minOccurences
		for {
			if parserState.AtTheEnd() {
				break
			}

			r, nextState := currentState.NextRune()
			if r == utf8.RuneError || parserState.AtTheEnd() || !conditionFn(r) {
				break
			}
			numFound = numFound + 1
			currentState = nextState
			results = utf8.AppendRune(results, r)
		}
		if numFound < minOccurences {
			return updateParserError(parserState, fmt.Errorf("CondMin: %d number of found are less then minOccurences: %d", numFound, minOccurences))
		}
		return updateParserState(parserState, currentState.Index, Result(string(results)))
	}
	return NewParser("CondMin()", parserFun)
}

// CondMinMax returns a Parser which tests the next rune in the input with the condition function.
// If the condition is met, the rune is consumed from the input and the parser succeeds at minimum of `minOccurences` times,
// but maximum of `maxOccurences` times.
// Otherwise the parser fails.
func CondMinMax(conditionFn func(rune) bool, minOccurences, maxOccurences int) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}

		if parserState.AtTheEnd() && minOccurences > 0 {
			return updateParserError(parserState, fmt.Errorf("CondMinMax: got Unexpected end of input."))
		}

		if minOccurences < 0 || minOccurences > maxOccurences {
			return updateParserError(parserState, fmt.Errorf("CondMinMax: wrong range of occurences min.: %d, max.: %d", minOccurences, maxOccurences))
		}

		currentState := parserState
		numFound := 0
		results := make([]byte, 0, 10)

		// Try to take as many occurences as possible, but at least minOccurences
		for {
			if numFound >= maxOccurences {
				break
			}
			r, nextState := currentState.NextRune()
			if nextState.IsError || !conditionFn(r) {
				break
			}
			numFound = numFound + 1
			currentState = nextState
			results = utf8.AppendRune(results, r)
		}
		if numFound < minOccurences {
			return updateParserError(parserState, fmt.Errorf("CondMinMax: %d number of found are less then minOccurences: %d", numFound, minOccurences))
		}
		return updateParserState(parserState, currentState.Index, Result(string(results)))
	}
	return NewParser("CondMinMax()", parserFun)
}

// AnyChar matches any character
var AnyChar = Cond(IsAnyChar)

// AnyStr matches any characters
var AnyStr = CondMin(IsAnyChar, 1)

// Letter is a parser that matches a single letter character with the target string
var Letter = Cond(IsAsciiLetter)

// Letters is a parser that matches one or more letter characters with the target string
var Letters = CondMin(IsAsciiLetter, 1)

// Digit is a parser that matches a singl digit character with the target string
var Digit = Cond(IsDigit)

// Digits is a parser that matches one or more digit characters with the target string
var Digits = CondMin(IsDigit, 1)

// Integer is a parser that matches one or more digit characters with the target string and returns with an int value
var Integer = Digits.Map(func(in Result) Result {
	strValue := in.(string)
	intValue, _ := strconv.Atoi(strValue)
	return Result(intValue)
})
