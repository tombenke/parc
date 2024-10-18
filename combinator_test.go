package parc

import (
	"fmt"
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
	for i := 0; i < b.N; i++ {
		flightIdentifier.Parse("UA666/12")
	}
}

func TestFlightIdentifier(t *testing.T) {

	validInputs := []string{"LH939/3", "UA666/12"}
	for _, input := range validInputs {
		newState := flightIdentifier.Parse(input)
		fmt.Printf("\n%+v\n", newState.Results)
		require.Equal(t, input, newState.Results)
		require.False(t, newState.IsError)
	}
}

func TestFlightNumber(t *testing.T) {

	validInputs := []string{"939", "666"}
	for _, input := range validInputs {
		newState := flightNumber.Parse(input)
		fmt.Printf("\n%+v\n", newState.Results)
		require.Equal(t, input, newState.Results)
		require.False(t, newState.IsError)
	}

}

func TestAirlineDesignator(t *testing.T) {
	//Debug(2)
	validInputs := []string{"LH", "X3A"}
	for _, input := range validInputs {
		fmt.Printf("\ninput: %+v\n", input)
		newState := airlineDesignator.Parse(input)
		fmt.Printf("\n%+v\n", newState.Results)
		require.Equal(t, input, newState.Results)
		require.False(t, newState.IsError)
	}
}

func TestDateOfDeparture(t *testing.T) {
	//Debug(2)
	validInputs := []string{"/1", "/12"}
	for _, input := range validInputs {
		fmt.Printf("\ninput: %+v\n", input)
		newState := dateOfDeparture.Parse(input)
		fmt.Printf("\n%+v\n", newState.Results)
		require.Equal(t, input, newState.Results)
		require.False(t, newState.IsError)
	}
}
