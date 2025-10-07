package testing

import (
	"github.com/tombenke/parc"
)

type NumberValue float64

type ASTNode struct {
	Tag   string
	Value NumberValue
}

var (
	Number = *parc.Choice(&RealNumber, &IntNumber)

	IntNumber = *parc.Map(parc.Integer, func(in parc.Result) parc.Result {
		node := ASTNode{
			Tag:   "INTEGER",
			Value: NumberValue(in.(int)),
		}
		return parc.Result(node)
	})

	RealNumber = *parc.Map(parc.RealNumber, func(in parc.Result) parc.Result {
		node := ASTNode{
			Tag:   "REAL",
			Value: NumberValue(in.(float64)),
		}
		return parc.Result(node)
	})
)
