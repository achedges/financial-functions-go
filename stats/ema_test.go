package stats_test

import (
	"financial/functions/stats"
	"testing"

	"github.com/achedges/go-assertions"
)

// we can't let the rolling calculation wrap all the way around, since the initial EMA calculation is an SMA

func checkStatsEma(function *stats.ExponentialMovingAvg, newindex int, t *testing.T) {
	function.Slide(data[newindex])
	testindex := (function.Buffer.Index - function.Buffer.Period + 1) % len(ema)
	assertions.CloseEnough(ema[testindex], function.GetValueFloat(), 0.0001, t)
}

func TestExponentialMovingAvg_Linear(t *testing.T) {
	prices := getPriceBarList(len(data))
	ema := stats.NewExponentialMovingAvg(period, mapClosePricesFromBars(prices))
	for ema.Buffer.Index < len(data)+2 {
		newindex := (ema.Buffer.Index + 1) % ema.Buffer.Length
		checkStatsEma(ema, newindex, t)
	}
}

func TestExponentialMovingAvg_Ring(t *testing.T) {
	prices := getPriceBarList(period)
	ema := stats.NewExponentialMovingAvg(period, mapClosePricesFromBars(prices))
	for ema.Buffer.Index < len(data)+2 {
		newindex := (ema.Buffer.Index + 1) % len(data)
		checkStatsEma(ema, newindex, t)
	}
}
