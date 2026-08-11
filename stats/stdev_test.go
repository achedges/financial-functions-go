package stats_test

import (
	"financial/functions/stats"
	"testing"

	"github.com/achedges/go-assertions"
)

func checkStatsStdev(function *stats.StandardDeviation, newindex int, t *testing.T) {
	function.Slide(data[newindex])
	testindex := (function.Buffer.Index - function.Buffer.Period + 1) % len(dev)
	assertions.CloseEnough(dev[testindex], function.GetValueFloat(), 0.0001, t)
}

func TestStandardDeviation_Linear(t *testing.T) {
	prices := getPriceBarList(len(data))
	stdev := stats.NewStandardDeviation(period, mapClosePricesFromBars(prices))
	for stdev.Buffer.Index < len(data)+10 {
		newindex := (stdev.Buffer.Index + 1) % stdev.Buffer.Length
		checkStatsStdev(stdev, newindex, t)
	}
}

func TestStandardDeviation_Ring(t *testing.T) {
	prices := getPriceBarList(period)
	stdev := stats.NewStandardDeviation(period, mapClosePricesFromBars(prices))
	for stdev.Buffer.Index < len(data)+10 {
		newindex := (stdev.Buffer.Index + 1) % len(data)
		checkStatsStdev(stdev, newindex, t)
	}
}
