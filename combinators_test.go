package parc

import (
	////"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	flightIdentifier  = SequenceOf(airlineDesignator, flightNumber, dateOfDeparture).Map(joinStrResults)
	airlineDesignator = SequenceOf(CondMinMax(IsAlphaNumeric, 2, 2), CondMinMax(IsAlphabetic, 0, 1)).Map(joinStrResults)
	flightNumber      = CondMinMax(IsDigit, 3, 4)
	dateOfDeparture   = SequenceOf(Char("/"), CondMinMax(IsDigit, 1, 2)).Map(joinStrResults)
	joinStrResults    = func(in Result) Result {
		resultsArr := in.([]Result)
		var results string
		for _, v := range resultsArr {
			results = results + v.(string)
		}
		return Result(results)
	}
)

func BenchmarkFlightIdentifier(b *testing.B) {
	input := "UA666/12"
	for i := 0; i < b.N; i++ {
		flightIdentifier.Parse(&input)
	}
}

func TestFlightIdentifier(t *testing.T) {

	validInputs := []string{"LH939/3", "UA666/12"}
	for _, input := range validInputs {
		newState := flightIdentifier.Parse(&input)
		////fmt.Printf("\n%+v\n", newState.Results)
		require.Equal(t, input, newState.Results)
		require.False(t, newState.IsError)
	}
}

func TestFlightNumber(t *testing.T) {

	validInputs := []string{"939", "666"}
	for _, input := range validInputs {
		newState := flightNumber.Parse(&input)
		////fmt.Printf("\n%+v\n", newState.Results)
		require.Equal(t, input, newState.Results)
		require.False(t, newState.IsError)
	}

}

func TestAirlineDesignator(t *testing.T) {
	//Debug(2)
	validInputs := []string{"LH", "X3A"}
	for _, input := range validInputs {
		////fmt.Printf("\ninput: %+v\n", input)
		newState := airlineDesignator.Parse(&input)
		////fmt.Printf("\n%+v\n", newState.Results)
		require.Equal(t, input, newState.Results)
		require.False(t, newState.IsError)
	}
}

func TestDateOfDeparture(t *testing.T) {
	//Debug(2)
	validInputs := []string{"/1", "/12"}
	for _, input := range validInputs {
		////fmt.Printf("\ninput: %+v\n", input)
		newState := dateOfDeparture.Parse(&input)
		////fmt.Printf("\n%+v\n", newState.Results)
		require.Equal(t, input, newState.Results)
		require.False(t, newState.IsError)
	}
}

func TestSequenceOf(t *testing.T) {
	input := "Hello World"
	token1 := "Hello"
	token2 := " "
	token3 := "World"
	expectedIndex := 11
	expectedError := error(nil)
	expectedResults := []Result{token1, token2, token3}
	expectedState := NewParserState(&input, expectedResults, expectedIndex, expectedError)

	sequenceParser := SequenceOf(
		Str(token1),
		Str(token2),
		Str(token3),
	)
	newState := sequenceParser.Parse(&input)
	require.Equal(t, expectedState, newState)

	wrongSequenceParser := SequenceOf(
		Str(token1),
		Str(token2),
		Str(token3),
		Str(token3), // Try one more at the end of the input
	)
	newState = wrongSequenceParser.Parse(&input)
	require.True(t, newState.IsError)
}

func TestTimes(t *testing.T) {
	input := "Hello Hello Hello Hello Hello "
	tokenNOk := "XXX "
	tokenOk := "Hello "
	expectedIndex := 24
	expectedError := error(nil)
	expectedResultsNOk := Result(nil)
	expectedResultsOk := []Result{tokenOk, tokenOk, tokenOk, tokenOk}
	expectedStateOk := NewParserState(&input, expectedResultsOk, expectedIndex, expectedError)

	timesParser := Times(Str(tokenOk), 4)

	newState := timesParser.Parse(&input)
	require.Equal(t, expectedStateOk, newState)

	timesParser = Times(Str(tokenNOk), 4)

	newState = timesParser.Parse(&input)
	require.Equal(t, expectedResultsNOk, newState.Results)
}

func TestTimesMin(t *testing.T) {
	input := "Hello Hello Hello Hello Hello "
	tokenNOk := "XXX "
	tokenOk := "Hello "
	expectedIndex := 30
	expectedError := error(nil)
	expectedResultsNOk := Result(nil)
	expectedResultsOk := []Result{tokenOk, tokenOk, tokenOk, tokenOk, tokenOk}
	expectedStateOk := NewParserState(&input, expectedResultsOk, expectedIndex, expectedError)

	timesMinParser := TimesMin(Str(tokenOk), 2)

	newState := timesMinParser.Parse(&input)
	require.Equal(t, expectedStateOk, newState)

	timesMinParser = TimesMin(Str(tokenNOk), 2)

	newState = timesMinParser.Parse(&input)
	require.Equal(t, expectedResultsNOk, newState.Results)
}

func TestTimesMinMax(t *testing.T) {
	input := "Hello Hello Hello Hello Hello "
	tokenNOk := "XXX "
	tokenOk := "Hello "
	expectedIndex := 18
	expectedError := error(nil)
	expectedResultsNOk := Result(nil)
	expectedResultsOk := []Result{tokenOk, tokenOk, tokenOk}
	expectedStateOk := NewParserState(&input, expectedResultsOk, expectedIndex, expectedError)

	timesMinMaxParser := TimesMinMax(Str(tokenOk), 2, 3)

	newState := timesMinMaxParser.Parse(&input)
	require.Equal(t, expectedStateOk, newState)

	timesMinMaxParser = TimesMinMax(Str(tokenNOk), 2, 3)

	newState = timesMinMaxParser.Parse(&input)
	require.Equal(t, expectedResultsNOk, newState.Results)
}

func TestZeroOrOne(t *testing.T) {
	input := "Hello Hello Hello Hello Hello "
	tokenZero := "XXX "
	tokenOne := "Hello "
	expectedIndex := 6
	expectedError := error(nil)
	expectedResultsZero := Result(nil)
	expectedStateZero := NewParserState(&input, expectedResultsZero, 0, expectedError)
	expectedResultsOne := Result(tokenOne)
	expectedStateOne := NewParserState(&input, expectedResultsOne, expectedIndex, expectedError)

	zeroOrOneParser := ZeroOrOne(Str(tokenZero))

	newState := zeroOrOneParser.Parse(&input)
	require.Equal(t, expectedStateZero, newState)

	zeroOrOneParser = ZeroOrOne(Str(tokenOne))

	newState = zeroOrOneParser.Parse(&input)
	require.Equal(t, expectedStateOne, newState)
}

func TestZeroOrMore(t *testing.T) {
	input := "Hello Hello Hello Hello Hello "
	tokenZero := "XXX "
	tokenMore := "Hello "
	expectedIndex := 30
	expectedError := error(nil)
	expectedResultsZero := []Result{}
	expectedStateZero := NewParserState(&input, expectedResultsZero, 0, expectedError)
	expectedResultsMore := []Result{tokenMore, tokenMore, tokenMore, tokenMore, tokenMore}
	expectedStateMore := NewParserState(&input, expectedResultsMore, expectedIndex, expectedError)

	zeroOrMoreParser := ZeroOrMore(Str(tokenZero))

	newState := zeroOrMoreParser.Parse(&input)
	require.Equal(t, expectedStateZero, newState)

	zeroOrMoreParser = ZeroOrMore(Str(tokenMore))

	newState = zeroOrMoreParser.Parse(&input)
	require.Equal(t, expectedStateMore, newState)
}

func TestOneOrMore(t *testing.T) {
	input := "Hello Hello "
	tokenOne := "Hello Hello "
	tokenMore := "Hello "
	tokenNone := "XXX "
	expectedIndex := 12
	expectedError := error(nil)
	expectedResultsOne := []Result{tokenOne}
	expectedStateOne := NewParserState(&input, expectedResultsOne, expectedIndex, expectedError)
	expectedResultsMore := []Result{tokenMore, tokenMore}
	expectedStateMore := NewParserState(&input, expectedResultsMore, expectedIndex, expectedError)

	oneOrMoreParser := OneOrMore(Str(tokenOne))

	newState := oneOrMoreParser.Parse(&input)
	require.Equal(t, expectedStateOne, newState)

	oneOrMoreParser = OneOrMore(Str(tokenMore))

	newState = oneOrMoreParser.Parse(&input)
	require.Equal(t, expectedStateMore, newState)

	oneOrMoreParser = OneOrMore(Str(tokenNone))

	newState = oneOrMoreParser.Parse(&input)
	require.Equal(t, Result(nil), newState.Results)
	require.True(t, newState.IsError)
}

func TestChoice(t *testing.T) {
	inputWithText := "Hello World"
	inputWithNumbers := "1342 234 45"
	inputWithPunct := "!., 1342 234 45"
	expectedIndexWithText := 5
	expectedIndexWithNumbers := 4
	expectedError := error(nil)

	choiceParser := Choice(
		Letters,
		Digits,
	)
	newState := choiceParser.Parse(&inputWithText)

	require.Equal(t, expectedIndexWithText, newState.Index)
	require.Equal(t, expectedError, newState.Err)
	require.False(t, newState.IsError)

	newState = choiceParser.Parse(&inputWithNumbers)

	require.Equal(t, expectedIndexWithNumbers, newState.Index)
	require.Equal(t, expectedError, newState.Err)
	require.False(t, newState.IsError)

	newState = choiceParser.Parse(&inputWithPunct)

	require.True(t, newState.IsError)
}

func TestBetween(t *testing.T) {
	input := "(42)"
	expectedResult := int(42)

	betweenParser := Between(Char("("), Char(")"))(Integer)
	newState := betweenParser.Parse(&input)
	require.Equal(t, expectedResult, newState.Results)
	require.False(t, newState.IsError)
}

func TestChain(t *testing.T) {
	stringInput := "string:Hello"
	numberInput := "number:42"
	dicerollInput := "diceroll:2d8"

	stringParser := Letters
	numberParser := Digits
	dicerollParser := SequenceOf(
		Integer,
		Char("d"),
		Integer,
	)
	parser := Chain(
		SequenceOf(Letters, Char(":")),
		func(result Result) *Parser {
			arr := result.([]Result)
			leftValue := arr[0].(string)
			switch leftValue {
			case "string":
				return stringParser
			case "number":
				return numberParser
			default:
				return dicerollParser
			}
		})

	newState := parser.Parse(&stringInput)
	require.False(t, newState.IsError)

	newState = parser.Parse(&numberInput)
	require.False(t, newState.IsError)

	newState = parser.Parse(&dicerollInput)
	require.False(t, newState.IsError)
}

func TestParser_Chain(t *testing.T) {
	//stringInput := "string:Hello"
	//numberInput := "number:42"
	dicerollInput := "diceroll:2d8"

	stringParser := Letters
	numberParser := Digits
	dicerollParser := SequenceOf(
		Integer,
		Char("d"),
		Integer,
	)
	parser := SequenceOf(Letters, Char(":")).Chain(
		func(result Result) *Parser {
			arr := result.([]Result)
			leftValue := arr[0].(string)
			switch leftValue {
			case "string":
				return stringParser
			case "number":
				return numberParser
			default:
				return dicerollParser
			}
		})
	newState := parser.Parse(&dicerollInput)
	require.False(t, newState.IsError)
}
