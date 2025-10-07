package testing

import (
	"github.com/stretchr/testify/require"
	"github.com/tombenke/parc"
	"testing"
)

func TestIntNumber(t *testing.T) {

	validTestCases := []parc.TestCase{
		parc.TestCase{
			Input:          "42",
			ExpectedResult: ASTNode{Tag: "INTEGER", Value: NumberValue(42)},
		},
		parc.TestCase{
			Input:          "12342",
			ExpectedResult: ASTNode{Tag: "INTEGER", Value: NumberValue(12342)},
		},
	}

	for _, tc := range validTestCases {
		newState := IntNumber.Parse(&tc.Input)
		require.Equal(t, tc.ExpectedResult, newState.Results.(ASTNode))
		require.False(t, newState.IsError)
	}

	invalidInputs := []string{"?23", "AB23", "A33"}
	for _, input := range invalidInputs {
		newState := IntNumber.Parse(&input)
		require.True(t, newState.IsError)
	}
}

func TestRealNumber(t *testing.T) {

	validTestCases := []parc.TestCase{
		parc.TestCase{
			Input:          "42.",
			ExpectedResult: ASTNode{Tag: "REAL", Value: NumberValue(42.)},
		},
		parc.TestCase{
			Input:          "-42.e12",
			ExpectedResult: ASTNode{Tag: "REAL", Value: NumberValue(-42.e12)},
		},
		parc.TestCase{
			Input:          "42.e-2",
			ExpectedResult: ASTNode{Tag: "REAL", Value: NumberValue(0.42)},
		},
	}

	for _, tc := range validTestCases {
		newState := RealNumber.Parse(&tc.Input)
		require.Equal(t, tc.ExpectedResult, newState.Results.(ASTNode))
		require.False(t, newState.IsError)
	}

	invalidInputs := []string{"?23", "AB23", "A33"}
	for _, input := range invalidInputs {
		newState := RealNumber.Parse(&input)
		require.True(t, newState.IsError)
	}
}

func TestNumber(t *testing.T) {

	validTestCases := []parc.TestCase{
		parc.TestCase{
			Input:          "42.",
			ExpectedResult: ASTNode{Tag: "REAL", Value: NumberValue(42.)},
		},
		parc.TestCase{
			Input:          "42",
			ExpectedResult: ASTNode{Tag: "INTEGER", Value: NumberValue(42)},
		},
		parc.TestCase{
			Input:          "-42.e12",
			ExpectedResult: ASTNode{Tag: "REAL", Value: NumberValue(-42.e12)},
		},
		parc.TestCase{
			Input:          "215",
			ExpectedResult: ASTNode{Tag: "INTEGER", Value: NumberValue(215)},
		},
	}

	for _, tc := range validTestCases {
		newState := Number.Parse(&tc.Input)
		require.Equal(t, tc.ExpectedResult, newState.Results.(ASTNode))
		require.False(t, newState.IsError)
	}

	invalidInputs := []string{"?23", "AB23", "A33"}
	for _, input := range invalidInputs {
		newState := Number.Parse(&input)
		require.True(t, newState.IsError)
	}
}
