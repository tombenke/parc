package parc

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

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
		Map(Digits, digitsToIntMapperFn),
		Str(" "),
		Str("Hello"),
	).Parse(&input)

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
		Digits.Map(digitsToIntMapperFn),
		Str(" "),
		Str("Hello"),
	).Parse(&input)

	require.False(t, newState.IsError)
}

func TestParser_ErrorMap(t *testing.T) {
	type MapResult struct {
		Tag   string
		Value int
	}
	input := "42 Hello World!"
	digitsToIntMapperFn := func(in Result) Result {
		strValue := in.(string)
		intValue, _ := strconv.Atoi(strValue)
		result := MapResult{
			Tag:   "INTEGER",
			Value: intValue,
		}
		return Result(result)
	}

	expectedError := fmt.Errorf("Catch SequenceOf error")
	newState := SequenceOf(
		Digits.Map(digitsToIntMapperFn),
		Choice(Str(", "), Str(" ")),
		Str("Hello"),
		// Generates error via missing Str(" "),
		Str("World"),
	).ErrorMap(func(state ParserState) error {
		return expectedError
	}).Parse(&input)

	require.Equal(t, expectedError, newState.Err)
	require.True(t, newState.IsError)
}
