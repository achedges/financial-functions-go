package window_test

import (
	"financial/functions/window"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/achedges/go-assertions"
)

func checkStats(function *window.BollingerBands, newindex int, t *testing.T) {
	function.Slide(data[newindex])
	testindex := (function.Buffer.Index - function.Buffer.Period + 1) % len(dev)

	mavg := function.GetMovingAvgValue()
	stdev := function.GetValue()
	ub := mavg.Add(stdev.Mul(decimal.NewFromInt(2)))
	lb := mavg.Sub(stdev.Mul(decimal.NewFromInt(2)))

	assertions.CloseEnough(sma[testindex], mavg.InexactFloat64(), 0.0001, t)
	assertions.CloseEnough(dev[testindex], stdev.InexactFloat64(), 0.0001, t)
	assertions.CloseEnough(ub.InexactFloat64(), function.UpperBand.InexactFloat64(), 0.0001, t)
	assertions.CloseEnough(lb.InexactFloat64(), function.LowerBand.InexactFloat64(), 0.0001, t)
}

func TestBollingerBands_LinearBuffer(t *testing.T) {
	bars := getPriceBarList(len(data))
	bb := window.NewBollingerBands(period, mapClosePricesFromBars(bars))
	for bb.Buffer.Index < len(data)+10 {
		newindex := (bb.Buffer.Index + 1) % bb.Buffer.Length
		checkStats(bb, newindex, t)
	}
}

func TestBollingerBands_RingBuffer(t *testing.T) {
	bars := getPriceBarList(period)
	bb := window.NewBollingerBands(period, mapClosePricesFromBars(bars))
	for bb.Buffer.Index < len(data)+10 {
		newindex := (bb.Buffer.Index + 1) % len(data)
		checkStats(bb, newindex, t)
	}
}