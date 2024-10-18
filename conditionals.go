package parc

import (
	"fmt"
	"unicode/utf8"
)

// IsAlphabetic is an alias of IsAsciiLetter
var IsAlphabetic = IsAsciiLetter

// IsAlphabetic Tests if byte is ASCII alphabetic: [A-Za-z]
func IsAsciiLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

// IsDigit is an alias of IsDecimalDigit
var IsDigit = IsDecimalDigit

// IsDecimalDigit Tests if byte is ASCII digit: [0-9]
func IsDecimalDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// IsHexadecimalDigit Tests if byte is ASCII hex digit: [0-9A-Fa-f]
func IsHexadecimalDigit(r rune) bool {
	return r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F'
}

// IsOctDigit Tests if byte is ASCII octal digit: [0-7]
func IsOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

// IsBinaryDigit Tests if byte is ASCII binary digit: [0-1]
func IsBinaryDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

// IsAlphaNumeric: Tests if byte is ASCII alphanumeric: [A-Za-z0-9]
func IsAlphaNumeric(r rune) bool {
	return IsAsciiLetter(r) || IsDecimalDigit(r)
}

func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

// IsSpace Tests if byte is ASCII space or tab: [ ]
func IsSpace(r rune) bool {
	return r == ' '
}

// IsTab Tests if byte is ASCII space or tab: [\t]
func IsTab(r rune) bool {
	return r == '\t'
}

// IsNewline Tests if byte is ASCII newline: [\n]
func IsNewline(r rune) bool {
	return r == '\n'
}

// IsCarriageReturn Tests if byte is ASCII newline: [\r]
func IsCarriageReturn(r rune) bool {
	return r == '\r'
}

// Cond returns a Parser which tests the next rune in the input with the condition function.
// If the condition is met, the rune is consumed from the input and the parser succeeds.
// Otherwise the parser fails.
func Cond(conditionFn func(rune) bool) *Parser {
	parserFun := func(parserState ParserState) ParserState {
		if parserState.IsError {
			return parserState
		}

		// TODO: Replace with Remaining()
		slicedInputString := parserState.InputString[parserState.Index:]

		// TODO: Replace with end-check
		if len(slicedInputString) == 0 {
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

		// TODO: Replace with Remaining()
		slicedInputString := parserState.InputString[parserState.Index:]

		// TODO: Replace with end-check
		if len(slicedInputString) == 0 {
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
			r, nextState := currentState.NextRune()
			if nextState.IsError || !conditionFn(r) {
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

		// TODO: Replace with Remaining()
		slicedInputString := parserState.InputString[parserState.Index:]

		// TODO: Replace with end-check
		if len(slicedInputString) == 0 && minOccurences > 0 {
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
