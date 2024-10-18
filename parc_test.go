package parc

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func init() {
	Debug(0)
}

func TestStartOfInput(t *testing.T) {
	newState := StartOfInput().Parse("")
	require.False(t, newState.IsError)

	newState = SequenceOf(
		StartOfInput(),
		Str("Hello"),
		Str(" "),
		Str("World"),
	).Parse("Hello World")
	require.False(t, newState.IsError)

	newState = SequenceOf(
		Str("Hello"),
		StartOfInput(), // Not at the start
		Str(" "),
		Str("World"),
	).Parse("Hello World")
	require.True(t, newState.IsError)
}

func TestEndOfInput(t *testing.T) {
	newState := EndOfInput().Parse("")
	require.False(t, newState.IsError)

	newState = SequenceOf(
		Str("Hello"),
		Str(" "),
		Str("World"),
		EndOfInput(),
	).Parse("Hello World")
	require.False(t, newState.IsError)

	newState = SequenceOf(
		Str("Hello"),
		EndOfInput(), // Not at the end
		Str(" "),
		Str("World"),
	).Parse("Hello World")
	require.True(t, newState.IsError)
}

func TestStr(t *testing.T) {

	input := "Hello World"
	token := "Hello"
	expectedIndex := 5
	expectedError := error(nil)
	expectedResults := token
	expectedState := ParserState{InputString: input, Results: expectedResults, Index: expectedIndex, Err: expectedError, IsError: false}

	newState := Str(token).Parse(input)
	require.Equal(t, expectedState, newState)

	// Try with an empty input
	newState = Str(token).Parse("")
	require.True(t, newState.IsError)
}

func TestInteger(t *testing.T) {

	numInput := "42"
	expectedIndex := 2
	expectedResult := Result(int(42))
	expectedError := error(nil)

	newState := Integer().Parse(numInput)

	require.Equal(t, expectedResult, newState.Results)
	require.Equal(t, expectedIndex, newState.Index)
	require.Equal(t, expectedError, newState.Err)
	require.False(t, newState.IsError)

	textInput := "Hello World!"

	newState = Integer().Parse(textInput)

	require.Equal(t, 0, newState.Index)
	require.True(t, newState.IsError)
}

func TestLetters(t *testing.T) {

	textInput := "Hello World!"
	expectedIndex := 5
	expectedError := error(nil)
	expectedResults := "Hello"
	expectedState := ParserState{InputString: textInput, Results: expectedResults, Index: expectedIndex, Err: expectedError, IsError: false}

	newState := Letters().Parse(textInput)
	require.Equal(t, expectedState, newState)

	numInput := "42 is the number of the Universe!"
	newState = Letters().Parse(numInput)
	require.Equal(t, 0, newState.Index)
	require.True(t, newState.IsError)
}

func TestDigits(t *testing.T) {

	numInput := "42 is the number of the Universe!"
	expectedIndex := 2
	expectedError := error(nil)
	expectedResults := "42"
	expectedState := ParserState{InputString: numInput, Results: expectedResults, Index: expectedIndex, Err: expectedError, IsError: false}

	newState := Digits().Parse(numInput)

	require.Equal(t, expectedState, newState)

	textInput := "Hello World!"

	newState = Digits().Parse(textInput)
	require.Equal(t, 0, newState.Index)
	require.True(t, newState.IsError)
}

func TestSequenceOf(t *testing.T) {
	input := "Hello World"
	token1 := "Hello"
	token2 := " "
	token3 := "World"
	expectedIndex := 11
	expectedError := error(nil)
	expectedResults := []Result{token1, token2, token3}
	expectedState := ParserState{InputString: input, Results: expectedResults, Index: expectedIndex, Err: expectedError, IsError: false}

	sequenceParser := SequenceOf(
		Str(token1),
		Str(token2),
		Str(token3),
	)
	newState := sequenceParser.Parse(input)
	require.Equal(t, expectedState, newState)

	wrongSequenceParser := SequenceOf(
		Str(token1),
		Str(token2),
		Str(token3),
		Str(token3), // Try one more at the end of the input
	)
	newState = wrongSequenceParser.Parse(input)
	require.True(t, newState.IsError)
}

func TestZeroOrMore(t *testing.T) {
	input := "Hello Hello Hello Hello Hello "
	tokenZero := "XXX "
	tokenMore := "Hello "
	expectedIndex := 30
	expectedError := error(nil)
	expectedResultsZero := []Result{}
	expectedStateZero := ParserState{InputString: input, Results: expectedResultsZero, Index: 0, Err: expectedError, IsError: false}
	expectedResultsMore := []Result{tokenMore, tokenMore, tokenMore, tokenMore, tokenMore}
	expectedStateMore := ParserState{InputString: input, Results: expectedResultsMore, Index: expectedIndex, Err: expectedError, IsError: false}

	zeroOrMoreParser := ZeroOrMore(Str(tokenZero))

	newState := zeroOrMoreParser.Parse(input)
	require.Equal(t, expectedStateZero, newState)

	zeroOrMoreParser = ZeroOrMore(Str(tokenMore))

	newState = zeroOrMoreParser.Parse(input)
	require.Equal(t, expectedStateMore, newState)
}

func TestOneOrMore(t *testing.T) {
	input := "Hello Hello "
	tokenOne := "Hello Hello "
	tokenMore := "Hello "
	tokenNone := "XXX "
	expectedIndex := 12
	expectedError := error(nil)
	expectedResultsOne := []Result{tokenOne}
	expectedStateOne := ParserState{InputString: input, Results: expectedResultsOne, Index: expectedIndex, Err: expectedError, IsError: false}
	expectedResultsMore := []Result{tokenMore, tokenMore}
	expectedStateMore := ParserState{InputString: input, Results: expectedResultsMore, Index: expectedIndex, Err: expectedError, IsError: false}

	oneOrMoreParser := OneOrMore(Str(tokenOne))

	newState := oneOrMoreParser.Parse(input)
	require.Equal(t, expectedStateOne, newState)

	oneOrMoreParser = OneOrMore(Str(tokenMore))

	newState = oneOrMoreParser.Parse(input)
	require.Equal(t, expectedStateMore, newState)

	oneOrMoreParser = OneOrMore(Str(tokenNone))

	newState = oneOrMoreParser.Parse(input)
	require.Equal(t, Result(nil), newState.Results)
	require.True(t, newState.IsError)
}

func TestChoice(t *testing.T) {
	inputWithText := "Hello World"
	inputWithNumbers := "1342 234 45"
	inputWithPunct := "!., 1342 234 45"
	expectedIndexWithText := 5
	expectedIndexWithNumbers := 4
	expectedError := error(nil)

	choiceParser := Choice(
		Letters(),
		Digits(),
	)
	newState := choiceParser.Parse(inputWithText)

	require.Equal(t, expectedIndexWithText, newState.Index)
	require.Equal(t, expectedError, newState.Err)
	require.False(t, newState.IsError)

	newState = choiceParser.Parse(inputWithNumbers)

	require.Equal(t, expectedIndexWithNumbers, newState.Index)
	require.Equal(t, expectedError, newState.Err)
	require.False(t, newState.IsError)

	newState = choiceParser.Parse(inputWithPunct)

	require.True(t, newState.IsError)
}

func TestMap(t *testing.T) {
	type MapResult struct {
		Tag   string
		Value int
	}
	input := "42 Hello"
	digitsToIntMapperFn := func(in Result) Result {
		strValue := in.(string)
		intValue, _ := strconv.Atoi(strValue)
		result := MapResult{
			Tag:   "INTEGER",
			Value: intValue,
		}
		return Result(result)
	}

	newState := SequenceOf(
		Map(Digits(), digitsToIntMapperFn),
		Str(" "),
		Str("Hello"),
	).Parse(input)

	require.False(t, newState.IsError)
}

func TestParser_Map(t *testing.T) {
	type MapResult struct {
		Tag   string
		Value int
	}
	input := "42 Hello"
	digitsToIntMapperFn := func(in Result) Result {
		strValue := in.(string)
		intValue, _ := strconv.Atoi(strValue)
		result := MapResult{
			Tag:   "INTEGER",
			Value: intValue,
		}
		return Result(result)
	}

	newState := SequenceOf(
		Digits().Map(digitsToIntMapperFn),
		Str(" "),
		Str("Hello"),
	).Parse(input)

	require.False(t, newState.IsError)
}

func TestBetween(t *testing.T) {
	input := "(42)"
	expectedResult := int(42)

	betweenParser := Between(Char("("), Char(")"))(Integer())
	newState := betweenParser.Parse(input)
	require.Equal(t, expectedResult, newState.Results)
	require.False(t, newState.IsError)
}

func TestChain(t *testing.T) {
	//stringInput := "string:Hello"
	//numberInput := "number:42"
	dicerollInput := "diceroll:2d8"

	stringParser := Letters()
	numberParser := Digits()
	dicerollParser := SequenceOf(
		Integer(),
		Char("d"),
		Integer(),
	)
	parser := Chain(
		SequenceOf(Letters(), Char(":")),
		func(result Result) *Parser {
			arr := result.([]Result)
			leftValue := arr[0].(string)
			switch leftValue {
			case "string":
				return stringParser
			case "number":
				return numberParser
			default:
				return dicerollParser
			}
		})
	newState := parser.Parse(dicerollInput)
	require.False(t, newState.IsError)
}

func TestParser_Chain(t *testing.T) {
	//stringInput := "string:Hello"
	//numberInput := "number:42"
	dicerollInput := "diceroll:2d8"

	stringParser := Letters()
	numberParser := Digits()
	dicerollParser := SequenceOf(
		Integer(),
		Char("d"),
		Integer(),
	)
	parser := SequenceOf(Letters(), Char(":")).Chain(
		func(result Result) *Parser {
			arr := result.([]Result)
			leftValue := arr[0].(string)
			switch leftValue {
			case "string":
				return stringParser
			case "number":
				return numberParser
			default:
				return dicerollParser
			}
		})
	newState := parser.Parse(dicerollInput)
	require.False(t, newState.IsError)
}
