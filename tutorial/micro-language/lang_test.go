package main

import (
	"testing"
)

var (
	formula = "(+ (* 10 2) (- (/ 50 3) 2))"
)

func BenchmarkMicroLanguage(b *testing.B) {
	parser := buildParser()
	parseResults := parser.Parse(&formula)
	b.ResetTimer()
	evaluate(parseResults.Results)
}
