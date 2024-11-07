package parc

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStartOfInput(t *testing.T) {
	input := ""
	newState := StartOfInput().Parse(&input)
	require.False(t, newState.IsError)

	hwInput := "Hello World"
	newState = SequenceOf(
		StartOfInput(),
		Str("Hello"),
		Str(" "),
		Str("World"),
	).Parse(&hwInput)
	require.False(t, newState.IsError)

	newState = SequenceOf(
		Str("Hello"),
		StartOfInput(), // Not at the start
		Str(" "),
		Str("World"),
	).Parse(&hwInput)
	require.True(t, newState.IsError)
}

func TestEndOfInput(t *testing.T) {
	input := ""
	newState := EndOfInput().Parse(&input)
	require.False(t, newState.IsError)

	hwInput := "Hello World"
	newState = SequenceOf(
		Str("Hello"),
		Str(" "),
		Str("World"),
		EndOfInput(),
	).Parse(&hwInput)
	require.False(t, newState.IsError)

	newState = SequenceOf(
		Str("Hello"),
		EndOfInput(), // Not at the end
		Str(" "),
		Str("World"),
	).Parse(&hwInput)
	require.True(t, newState.IsError)
}

func TestRest(t *testing.T) {
	input := "Hello World"
	expectedError := error(nil)
	expectedResults := Result(input)
	expectedIndex := 11
	expectedState := NewParserState(&input, expectedResults, expectedIndex, expectedError)

	newState := Rest().Parse(&input)
	require.Equal(t, expectedState, newState)
	require.False(t, newState.IsError)

	newState = SequenceOf(Str("Hello"), Rest()).Parse(&input)
	expectedResults = []Result([]Result{"Hello", " World"})
	expectedState = NewParserState(&input, expectedResults, expectedIndex, expectedError)
	require.Equal(t, expectedState, newState)
	require.False(t, newState.IsError)

	newState = SequenceOf(Str("Hello"), Str(" "), Str("World"), Rest()).Parse(&input)
	expectedResults = []Result([]Result{"Hello", " ", "World", ""})
	expectedState = NewParserState(&input, expectedResults, expectedIndex, expectedError)
	require.Equal(t, expectedState, newState)
	require.False(t, newState.IsError)
}

func TestStr(t *testing.T) {

	input := "Hello World"
	token := "Hello"
	expectedIndex := 5
	expectedError := error(nil)
	expectedResults := token
	expectedState := NewParserState(&input, expectedResults, expectedIndex, expectedError)

	newState := Str(token).Parse(&input)
	require.Equal(t, expectedState, newState)

	// Try with an empty input
	emptyInput := ""
	newState = Str(token).Parse(&emptyInput)
	require.True(t, newState.IsError)
}
