package stats_test

import (
	"financial/functions/stats"
	"testing"

	"github.com/achedges/financial-core-go/pricebar"
	"github.com/achedges/go-assertions"
)

func TestRelativeStrengthIndex_Slide(t *testing.T) {
	bars := []pricebar.PriceBar{
		*pricebar.New(pricebar.Config{BasisPrice: 54.80}),
		*pricebar.New(pricebar.Config{BasisPrice: 56.80}),
		*pricebar.New(pricebar.Config{BasisPrice: 57.85}),
		*pricebar.New(pricebar.Config{BasisPrice: 59.85}),
		*pricebar.New(pricebar.Config{BasisPrice: 60.57}),
		*pricebar.New(pricebar.Config{BasisPrice: 61.10}),
		*pricebar.New(pricebar.Config{BasisPrice: 62.17}),
		*pricebar.New(pricebar.Config{BasisPrice: 60.60}),
		*pricebar.New(pricebar.Config{BasisPrice: 62.35}),
		*pricebar.New(pricebar.Config{BasisPrice: 62.15}),
		*pricebar.New(pricebar.Config{BasisPrice: 62.35}),
		*pricebar.New(pricebar.Config{BasisPrice: 61.45}),
		*pricebar.New(pricebar.Config{BasisPrice: 62.80}),
		*pricebar.New(pricebar.Config{BasisPrice: 61.37}),
	}

	rsi := stats.NewRelativeStrengthIndex(len(bars), bars)
	rsi.Slide(*pricebar.New(pricebar.Config{BasisPrice: 62.50}))

	assertions.CloseEnough(2.39, rsi.RelativeStrength().InexactFloat64(), 0.001, t)
	assertions.CloseEnough(70.50, rsi.RSI().InexactFloat64(), 0.01, t)
}
