package parc

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	Debug(0)
}

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
