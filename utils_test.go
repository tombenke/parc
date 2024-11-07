package parc

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAnyChar(t *testing.T) {
	input := "ű"
	expectedIndex := 2
	expectedError := error(nil)
	expectedResults := input
	expectedState := NewParserState(&input, expectedResults, expectedIndex, expectedError)

	newState := AnyChar.Parse(&input)
	require.Equal(t, expectedState, newState)
}

func TestAnyStr(t *testing.T) {
	input := "Ez egy szöveg. Mindenféle betűt és számot pl.: 42, illetve írásjeleket (?!%'*) is tartalmaz"
	expectedIndex := 98
	expectedError := error(nil)
	expectedResults := input
	expectedState := NewParserState(&input, expectedResults, expectedIndex, expectedError)

	newState := AnyStr.Parse(&input)
	require.Equal(t, expectedState, newState)
}

func TestInteger(t *testing.T) {

	numInput := "42"
	expectedIndex := 2
	expectedResult := Result(int(42))
	expectedError := error(nil)

	newState := Integer.Parse(&numInput)

	require.Equal(t, expectedResult, newState.Results)
	require.Equal(t, expectedIndex, newState.Index)
	require.Equal(t, expectedError, newState.Err)
	require.False(t, newState.IsError)

	textInput := "Hello World!"

	newState = Integer.Parse(&textInput)

	require.Equal(t, 0, newState.Index)
	require.True(t, newState.IsError)
}

func TestLetters(t *testing.T) {

	textInput := "Hello World!"
	expectedIndex := 5
	expectedError := error(nil)
	expectedResults := "Hello"
	expectedState := NewParserState(&textInput, expectedResults, expectedIndex, expectedError)

	newState := Letters.Parse(&textInput)
	require.Equal(t, expectedState, newState)

	numInput := "42 is the number of the Universe!"
	newState = Letters.Parse(&numInput)
	require.Equal(t, 0, newState.Index)
	require.True(t, newState.IsError)
}

func TestDigits(t *testing.T) {

	numInput := "42 is the number of the Universe!"
	expectedIndex := 2
	expectedError := error(nil)
	expectedResults := "42"
	expectedState := NewParserState(&numInput, expectedResults, expectedIndex, expectedError)

	newState := Digits.Parse(&numInput)

	require.Equal(t, expectedState, newState)

	textInput := "Hello World!"

	newState = Digits.Parse(&textInput)
	require.Equal(t, 0, newState.Index)
	require.True(t, newState.IsError)
}

func TestRestOfLine(t *testing.T) {
	singleLine := "Ez a szöveg első sora"
	multiLine := singleLine + "\nMindenféle betűt és számot pl.: 42, illetve írásjeleket (?!%'*) is tartalmaz.\nTöbb sorból áll"

	newState := RestOfLine.Parse(&singleLine)
	require.Equal(t, singleLine, newState.Results.(string))
	require.False(t, newState.IsError)

	newState = RestOfLine.Parse(&multiLine)
	require.Equal(t, singleLine, newState.Results.(string))
	require.False(t, newState.IsError)
}
