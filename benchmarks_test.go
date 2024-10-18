package parc

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	letters          = "ARelativelyLongTextToParseWithRegExpAndWithCond"
	regexpParser     = Letters()
	condMinParser    = CondMin(IsAsciiLetter, 0)
	condMinMaxParser = CondMinMax(IsAsciiLetter, 0, len(letters))
)

func TestCompareResultsOfBenchmarks(t *testing.T) {

	regexpResults := regexpParser.Parse(letters)
	condMinResults := condMinParser.Parse(letters)
	condMinMaxResults := condMinMaxParser.Parse(letters)

	require.Equal(t, regexpResults, condMinResults)
	require.Equal(t, regexpResults, condMinMaxResults)

}

func BenchmarkRegExp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Letters().Parse(letters)
	}
}

func BenchmarkCondMin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CondMin(IsAsciiLetter, 0).Parse(letters)
	}
}

func BenchmarkCondMinMax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CondMinMax(IsAsciiLetter, 0, len(letters)).Parse(letters)
	}
}
