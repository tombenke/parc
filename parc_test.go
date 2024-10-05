package parc

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestStr(t *testing.T) {

	input := "Hello World"
	token := "Hello"
	expectedIndex := 5
	expectedError := error(nil)

	results := Str(token).Parse(input)
	fmt.Printf("\n  results: %+v\n", results)

	//require.Equal(t, expectedResults, results)
	require.Equal(t, expectedIndex, results.Index)
	require.Equal(t, expectedError, results.Err)
	require.False(t, results.IsError)
}

func TestInteger(t *testing.T) {

	numInput := "42"
	expectedIndex := 2
	expectedResult := Result(int(42))
	expectedError := error(nil)

	results := Integer().Parse(numInput)
	fmt.Printf("\n  results: %+v\n", results)

	require.Equal(t, expectedResult, results.Results)
	require.Equal(t, expectedIndex, results.Index)
	require.Equal(t, expectedError, results.Err)
	require.False(t, results.IsError)

	textInput := "Hello World!"

	results = Integer().Parse(textInput)
	fmt.Printf("\n  results: %+v\n", results)

	//require.Equal(t, expectedResults, results)
	require.Equal(t, 0, results.Index)
	//require.Equal(t, expectedError, results.Err)
	require.True(t, results.IsError)
}

func TestLetters(t *testing.T) {

	textInput := "Hello World!"
	expectedIndex := 5
	expectedError := error(nil)

	results := Letters().Parse(textInput)
	fmt.Printf("\n  results: %+v\n", results)

	//require.Equal(t, expectedResults, results)
	require.Equal(t, expectedIndex, results.Index)
	require.Equal(t, expectedError, results.Err)
	require.False(t, results.IsError)

	numInput := "42 is the number of the Universe!"
	results = Letters().Parse(numInput)
	fmt.Printf("\n  results: %+v\n", results)

	require.Equal(t, 0, results.Index)
	//require.Equal(t, expectedError, results.Err)
	require.True(t, results.IsError)
}

func TestDigits(t *testing.T) {

	numInput := "42 is the number of the Universe!"
	expectedIndex := 2
	expectedError := error(nil)

	results := Digits().Parse(numInput)
	fmt.Printf("\n  results: %+v\n", results)

	require.Equal(t, expectedIndex, results.Index)
	require.Equal(t, expectedError, results.Err)
	require.False(t, results.IsError)

	textInput := "Hello World!"

	results = Digits().Parse(textInput)
	fmt.Printf("\n  results: %+v\n", results)

	//require.Equal(t, expectedResults, results)
	require.Equal(t, 0, results.Index)
	//require.Equal(t, expectedError, results.Err)
	require.True(t, results.IsError)
}

func TestSequenceOf(t *testing.T) {
	input := "Hello World"
	token1 := "Hello "
	token2 := "World"
	expectedIndex := 11
	expectedError := error(nil)

	sequenceParser := SequenceOf(
		Str(token1),
		Str(token2),
	)
	results := sequenceParser.Parse(input)
	fmt.Printf("\n  results: %+v\n", results)

	require.Equal(t, expectedIndex, results.Index)
	require.Equal(t, expectedError, results.Err)
	require.False(t, results.IsError)
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
	results := choiceParser.Parse(inputWithText)
	fmt.Printf("\n  results with text: %+v\n", results)

	require.Equal(t, expectedIndexWithText, results.Index)
	require.Equal(t, expectedError, results.Err)
	require.False(t, results.IsError)

	results = choiceParser.Parse(inputWithNumbers)
	fmt.Printf("\n  results with numbers: %+v\n", results)

	require.Equal(t, expectedIndexWithNumbers, results.Index)
	require.Equal(t, expectedError, results.Err)
	require.False(t, results.IsError)

	results = choiceParser.Parse(inputWithPunct)
	fmt.Printf("\n  results with numbers: %+v\n", results)

	require.True(t, results.IsError)
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

	results := SequenceOf(
		Map(Digits(), digitsToIntMapperFn),
		Str(" "),
		Str("Hello"),
	).Parse(input)
	fmt.Printf("\n  results: %+v\n", results)
	require.False(t, results.IsError)
}
