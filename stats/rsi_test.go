package stats_test

import (
	"testing"

	"github.com/achedges/financial-core-go/candle"
	"github.com/achedges/financial-functions-go/stats"
	"github.com/achedges/go-assertions"
)

func TestRelativeStrengthIndex_Slide(t *testing.T) {
	bars := []candle.Candle{
		candle.New(candle.Config{BasisPrice: 54.80}),
		candle.New(candle.Config{BasisPrice: 56.80}),
		candle.New(candle.Config{BasisPrice: 57.85}),
		candle.New(candle.Config{BasisPrice: 59.85}),
		candle.New(candle.Config{BasisPrice: 60.57}),
		candle.New(candle.Config{BasisPrice: 61.10}),
		candle.New(candle.Config{BasisPrice: 62.17}),
		candle.New(candle.Config{BasisPrice: 60.60}),
		candle.New(candle.Config{BasisPrice: 62.35}),
		candle.New(candle.Config{BasisPrice: 62.15}),
		candle.New(candle.Config{BasisPrice: 62.35}),
		candle.New(candle.Config{BasisPrice: 61.45}),
		candle.New(candle.Config{BasisPrice: 62.80}),
		candle.New(candle.Config{BasisPrice: 61.37}),
	}

	rsi := stats.NewRelativeStrengthIndex(len(bars), bars)
	rsi.Slide(candle.New(candle.Config{BasisPrice: 62.50}))

	assertions.CloseEnough(2.39, rsi.RelativeStrength().InexactFloat64(), 0.001, t)
	assertions.CloseEnough(70.50, rsi.RSI().InexactFloat64(), 0.01, t)
}
