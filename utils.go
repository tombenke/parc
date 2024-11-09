package parc

import (
	"fmt"
	"math"
	"strconv"
)

var (
	// Newline matches a space character ` `
	Space = Char(" ").As("Space")

	// Newline matches a newline character \n
	Newline = Char("\n").As("Newline")

	// Tab matches a tab character \t
	Tab = Char("\t").As("Tab")

	// Crlf recognizes the string \r\n
	Crlf = Str("\r\n").As("Crlf")

	// AnyChar matches any character
	AnyChar = Cond(IsAnyChar).As("AnyChar")

	// AnyStr matches any characters
	AnyStr = CondMin(IsAnyChar, 1).As("AnyStr")

	// Letter is a parser that matches a single letter character with the target string
	Letter = Cond(IsAsciiLetter).As("Letter")

	// Letters is a parser that matches one or more letter characters with the target string
	Letters = CondMin(IsAsciiLetter, 1).As("Letters")

	// Digit is a parser that matches a singl digit character with the target string
	Digit = Cond(IsDigit).As("Digit")

	// Digits is a parser that matches one or more digit characters with the target string
	Digits = CondMin(IsDigit, 1).As("Digits")

	// NonNegativeInteger is a parser that matches one or more digit characters with the target string and returns with a non negative int value
	NonNegativeInteger = Digits.Map(func(in Result) Result {
		strValue := in.(string)
		intValue, _ := strconv.Atoi(strValue)
		return Result(intValue)
	})

	// WholeNumber is an aliad for a NonNegativeInteger
	WholeNumber = NonNegativeInteger

	// Integer is a parser that matches one or more digit characters with the target string and returns with an int value
	Integer = SequenceOf(Sign, NonNegativeInteger).Map(func(result Result) Result {

		if resultsArr, ok := result.([]Result); ok {
			if integerValue, ok := resultsArr[1].(int); ok {
				if resultsArr[0] == "-" {
					return integerValue * -1
				} else {
					return integerValue
				}
			}
		}
		return Result(nil)
	})

	// RealNumber is a parser for a double type real number
	RealNumber = SequenceOf(Sign, Digits, Char("."), Optional(Digits), Exponent).Map(func(result Result) Result {
		resultArr := result.([]Result)

		sign := resultArr[0].(string)
		intPart := resultArr[1].(string)
		fractPart := "0"
		if fractPartValue, ok := resultArr[3].(string); ok {
			fractPart = fractPartValue
		}
		exponent := resultArr[4].(int)

		realValue, _ := strconv.ParseFloat(fmt.Sprintf("%s%s.%s", sign, intPart, fractPart), 64)
		realValue = realValue * math.Pow10(exponent)

		return Result(realValue)
	}).As("RealNumber")

	// Sign is a parser for the sign of a number value. It has a default value, that is "+".
	Sign = Optional(Choice(Char("+"), Char("-"))).Map(func(result Result) Result {
		sign := "+"
		if signValue, ok := result.(string); ok {
			sign = signValue
		}
		return sign

	}).As("Sign")

	// Exponent is a parser for the exponent part of a real number
	Exponent = Optional(SequenceOf(Choice(Char("e"), Char("E")), Sign, Integer)).Map(func(result Result) Result {

		exponent := 0
		if resultsArr, ok := result.([]Result); ok {
			if exponentValue, ok := resultsArr[2].(int); ok {
				if resultsArr[1] == "-" {
					exponent = exponentValue * -1
				} else {
					exponent = exponentValue
				}
			}
		}

		return exponent
	})

	// RestOfLine returns with the content of the input string until the next newline (\n) character,
	// or until the end of the input string, if there is no newline found.
	RestOfLine = SequenceOf(
		CondMin(func(r rune) bool { return r != '\n' }, 1),
		Choice(Newline, EndOfInput()),
	).As("RestOfLine").Map(func(in Result) Result {
		sequence := in.([]Result)
		return sequence[0]
	})
)

// JoinStrResults merges the a string-array result into a single string. The items of the array must be string type.
// TODO: Handle non-string type items, e.g. skip nils, stringify other types, etc.
func JoinStrResults(in Result) Result {
	resultsArr := in.([]Result)
	var results string
	for _, v := range resultsArr {
		results = results + v.(string)
	}
	return Result(results)
}

// Ref creates a reference to any value
// It is useful to define reference values of fixtures in test cases
func Ref[T any](value T) *T {
	var v T = value
	return &v
}

// Generic TestCase struct to help writing test cases for sub-parsers
type TestCase struct {
	Input          string
	ExpectedResult Result
}

// GetResultsItem takes the nth item from the results array, if there is any, otherwise it returns nil value
func GetResultsItem[T any](result Result, itemIdx int) *T {
	if debugLevel >= 3 {
		fmt.Printf("\nGetResultsItem(%+v, %d) => ", result, itemIdx)
	}
	if resultArr, ok := result.([]Result); ok {
		if value, ok := resultArr[itemIdx].(T); ok {
			//var result T
			result := value
			if debugLevel >= 3 {
				fmt.Printf("%+v\n", &result)
			}
			return &result
		}
	}
	return nil
}
