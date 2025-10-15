package errors

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLangParser(t *testing.T) {
	parser := buildParser()

	for _, testCase := range validTestCases {
		input := testCase.Formula
		parseResults := parser.Parse(&input)
		assert.NotNil(t, parseResults.Results)
		assert.Equal(t, testCase.Error, parseResults.Err)
	}

	for _, testCase := range invalidTestCases {
		input := testCase.Formula
		parseResults := parser.Parse(&input)
		assert.Nil(t, parseResults.Results)
		//assert.Equal(t, testCase.Error, parseResults.Err)
		//fmt.Printf("parseResults: %+v\n", parseResults)
		fmt.Printf("parseResults.Err: %+v\n", parseResults.Err)
	}
}
