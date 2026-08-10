package functions_test

import (
	"financial/functions"
	"testing"

	"github.com/achedges/go-assertions"
)

func checkStatsSma(function *functions.SimpleMovingAvg, newindex int, t *testing.T) {
	function.Slide(data[newindex])
	testindex := (function.Buffer.Index - function.Buffer.Period + 1) % len(sma)
	assertions.EqualFloats(sma[testindex], function.GetValueFloat(), t)
}

func TestSimpleMovingAvg_Linear(t *testing.T) {
	prices := getPriceBarList(len(data))
	sma := functions.NewSimpleMovingAvg(period, mapClosePricesFromBars(prices))
	for sma.Buffer.Index < len(data)+10 {
		newindex := (sma.Buffer.Index + 1) % sma.Buffer.Length
		checkStatsSma(sma, newindex, t)
	}
}

func TestSimpleMovingAvg_Ring(t *testing.T) {
	prices := getPriceBarList(period)
	sma := functions.NewSimpleMovingAvg(period, mapClosePricesFromBars(prices))
	for sma.Buffer.Index < len(data)+10 {
		newindex := (sma.Buffer.Index + 1) % len(data)
		checkStatsSma(sma, newindex, t)
	}
}
