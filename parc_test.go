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

	//initialState := state{data: input, offset: 0}
	//expectedState := state{data: input, offset: len(input)}
	//expectedError := error(nil)
	//expectedResult := int(42)

	//parserT := GetString(Exactly(token))
	//handlerFn := func(digits string) Parser[int] {
	//	if len(digits) > 1 && digits[0] == '0' {
	//		return Fail[int]
	//	}
	//	v, err := strconv.Atoi(digits)
	//	if err != nil {
	//		return Fail[int]
	//	}
	//	return Succeed(v)
	//}

	//andThenParser := AndThen(parserT, handlerFn)

	//result, newState, err := andThenParser(initialState)
	//fmt.Printf("\n  result: %+v\nnewState: %+v\n     err: %+v\n\n", result, newState, err)

	//require.Equal(t, expectedResult, result)
	//require.Equal(t, expectedState, newState)
	//require.Equal(t, expectedError, err)
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
