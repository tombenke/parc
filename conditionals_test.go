package parc

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	Debug(0)
}

func TestCond_Ok(t *testing.T) {

	input := "Hello World"
	expectedIndex := 1
	expectedError := error(nil)
	expectedResults := "H"
	expectedState := NewParserState(input, expectedResults, expectedIndex, expectedError)

	newState := Cond(IsAsciiLetter).Parse(input)
	require.Equal(t, expectedState, newState)
}

func TestCond_Err(t *testing.T) {
	input := "Hello World"
	newState := Cond(IsDecimalDigit).Parse(input)
	require.True(t, newState.IsError)
}

func TestCondMin_Ok(t *testing.T) {

	input := "Hello World"
	expectedIndex := 5
	expectedError := error(nil)
	expectedResults := "Hello"
	expectedState := NewParserState(input, expectedResults, expectedIndex, expectedError)

	newState := CondMin(IsAsciiLetter, 3).Parse(input)
	require.Equal(t, expectedState, newState)
}

func TestCondMin_Err(t *testing.T) {
	input := "Hello World"
	newState := CondMin(IsAsciiLetter, 8).Parse(input)
	require.True(t, newState.IsError)
}

func TestCondMin0_Ok(t *testing.T) {

	input := "Hello World"
	expectedIndex := 0
	expectedError := error(nil)
	expectedResults := ""
	expectedState := NewParserState(input, expectedResults, expectedIndex, expectedError)

	newState := CondMin(IsDecimalDigit, 0).Parse(input)
	require.Equal(t, expectedState, newState)
}

func TestCondMinMax_Ok(t *testing.T) {

	input := "Hello World"
	expectedError := error(nil)
	expectedState := NewParserState(input, "Hell", 4, expectedError)

	newState := CondMinMax(IsAsciiLetter, 3, 4).Parse(input)
	require.Equal(t, expectedState, newState)

	newState = CondMinMax(IsAsciiLetter, 3, 20).Parse(input)
	expectedState = NewParserState(input, "Hello", 5, expectedError)
	require.Equal(t, expectedState, newState)
}
