package parc

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func init() {
	Debug(0)
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
		Map(Digits, digitsToIntMapperFn),
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
		Digits.Map(digitsToIntMapperFn),
		Str(" "),
		Str("Hello"),
	).Parse(input)

	require.False(t, newState.IsError)
}
