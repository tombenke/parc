package parc

import (
	"github.com/stretchr/testify/require"
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
