package parc

import (
	"fmt"
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
	testCases := []TestCase{
		TestCase{Input: "42", ExpectedResult: Result(int(42))},
		TestCase{Input: "+41", ExpectedResult: Result(int(41))},
		TestCase{Input: "-24", ExpectedResult: Result(int(-24))},
		TestCase{Input: "0", ExpectedResult: Result(int(0))},
		TestCase{Input: "-0", ExpectedResult: Result(int(0))},
		TestCase{Input: "+0", ExpectedResult: Result(int(0))},
	}

	for _, tc := range testCases {
		newState := Integer.Parse(&tc.Input)
		require.Equal(t, tc.ExpectedResult, newState.Results)
		require.False(t, newState.IsError)
	}

	testErrorCases := []TestCase{
		TestCase{Input: "some text", ExpectedResult: Result(nil)},
	}

	for _, tc := range testErrorCases {
		newState := Integer.Parse(&tc.Input)
		require.Equal(t, tc.ExpectedResult, newState.Results)
		require.True(t, newState.IsError)
	}
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

func TestSign(t *testing.T) {
	testCases := []TestCase{
		TestCase{Input: "", ExpectedResult: "+"},
		TestCase{Input: "+", ExpectedResult: "+"},
		TestCase{Input: "-", ExpectedResult: "-"},
	}

	for _, tc := range testCases {
		newState := Sign.Parse(&tc.Input)
		require.Equal(t, tc.ExpectedResult, newState.Results)
		require.False(t, newState.IsError)
	}
}

func TestExponent(t *testing.T) {
	testCases := []TestCase{
		TestCase{Input: "", ExpectedResult: Result(int(0))},
		TestCase{Input: "e1", ExpectedResult: Result(int(1))},
		TestCase{Input: "E0", ExpectedResult: Result(int(0))},
		TestCase{Input: "e-2", ExpectedResult: Result(int(-2))},
		TestCase{Input: "E+3", ExpectedResult: Result(int(+3))},
	}

	for _, tc := range testCases {
		newState := Exponent.Parse(&tc.Input)
		fmt.Printf("\nresult: %+v\n", newState.Results)
		require.Equal(t, tc.ExpectedResult, newState.Results)
		require.False(t, newState.IsError)
	}
}
func TestRealNumber(t *testing.T) {
	testCases := []TestCase{
		TestCase{Input: "0.", ExpectedResult: Result(float64(0.))},
		TestCase{Input: "0.0", ExpectedResult: Result(float64(0.))},
		TestCase{Input: "3.1415", ExpectedResult: Result(float64(3.1415))},
		TestCase{Input: "-3.1415", ExpectedResult: Result(float64(-3.1415))},
		TestCase{Input: "-3.14E2", ExpectedResult: Result(float64(-314.))},
		TestCase{Input: "-2500.e-2", ExpectedResult: Result(float64(-25.))},
		TestCase{Input: "2500.e0", ExpectedResult: Result(float64(2500.))},
	}

	for _, tc := range testCases {
		newState := RealNumber.Parse(&tc.Input)
		require.Equal(t, tc.ExpectedResult, newState.Results)
		require.False(t, newState.IsError)
	}
}
