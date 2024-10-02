package parc

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStr(t *testing.T) {

	input := "Hello World"
	token := "Hello"
	expectedIndex := 5
	expectedError := error(nil)

	results := Str(token).Run(input)
	fmt.Printf("\n  results: %+v\n", results)

	//require.Equal(t, expectedResults, results)
	require.Equal(t, expectedIndex, results.Index)
	require.Equal(t, expectedError, results.Err)
	require.False(t, results.IsError)
}

func TestLetters(t *testing.T) {

	textInput := "Hello World!"
	expectedIndex := 5
	expectedError := error(nil)

	results := Letters().Run(textInput)
	fmt.Printf("\n  results: %+v\n", results)

	//require.Equal(t, expectedResults, results)
	require.Equal(t, expectedIndex, results.Index)
	require.Equal(t, expectedError, results.Err)
	require.False(t, results.IsError)

	numInput := "42 is the number of the Universe!"
	results = Letters().Run(numInput)
	fmt.Printf("\n  results: %+v\n", results)

	require.Equal(t, 0, results.Index)
	//require.Equal(t, expectedError, results.Err)
	require.True(t, results.IsError)
}

func TestDigits(t *testing.T) {

	numInput := "42 is the number of the Universe!"
	expectedIndex := 2
	expectedError := error(nil)

	results := Digits().Run(numInput)
	fmt.Printf("\n  results: %+v\n", results)

	require.Equal(t, expectedIndex, results.Index)
	require.Equal(t, expectedError, results.Err)
	require.False(t, results.IsError)

	textInput := "Hello World!"

	results = Digits().Run(textInput)
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
	results := sequenceParser.Run(input)
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
	results := choiceParser.Run(inputWithText)
	fmt.Printf("\n  results with text: %+v\n", results)

	require.Equal(t, expectedIndexWithText, results.Index)
	require.Equal(t, expectedError, results.Err)
	require.False(t, results.IsError)

	results = choiceParser.Run(inputWithNumbers)
	fmt.Printf("\n  results with numbers: %+v\n", results)

	require.Equal(t, expectedIndexWithNumbers, results.Index)
	require.Equal(t, expectedError, results.Err)
	require.False(t, results.IsError)

	results = choiceParser.Run(inputWithPunct)
	fmt.Printf("\n  results with numbers: %+v\n", results)

	require.True(t, results.IsError)
}
