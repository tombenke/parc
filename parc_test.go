package parc

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestStr(t *testing.T) {

	input := "Hello World"
	token := "Hello"
	expectedIndex := 5
	expectedError := error(nil)
	expectedResults := token
	expectedState := ParserState{InputString: input, Results: expectedResults, Index: expectedIndex, Err: expectedError, IsError: false}

	newState := Str(token).Parse(input)
	require.Equal(t, expectedState, newState)
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

func TestBetween(t *testing.T) {
	Debug(1)
	input := "(42)"
	expectedResult := int(42)

	betweenParser := Between(Char("("), Char(")"))(Integer())
	newState := betweenParser.Parse(input)
	require.Equal(t, expectedResult, newState.Results)
	require.False(t, newState.IsError)
}

func TestChain(t *testing.T) {
	Debug(2)
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
