package functions_test

import (
	"financial/functions"
	"testing"

	"github.com/achedges/financial-core-go/pricebar"
	"github.com/achedges/go-assertions"
)

var priceData = [][]float64{
	{50.00, 51.20, 49.80, 50.90},
	{50.70, 51.80, 50.30, 51.50},
	{51.70, 52.90, 51.70, 52.80},
	{52.50, 53.70, 52.30, 53.50},
	{53.60, 54.80, 53.50, 54.70},
	{54.40, 54.40, 52.90, 53.00},
	{52.90, 53.20, 52.00, 52.00},
	{52.00, 52.70, 52.00, 52.20},
}

var knownRanges = []float64{1.40, 1.50, 1.40, 1.40, 1.30, 1.80, 1.20, 0.70}

func getPriceBars() *[]pricebar.PriceBar {
	bars := make([]pricebar.PriceBar, len(priceData))
	for i, v := range priceData {
		bar := pricebar.New(pricebar.Config{Symbol: "TEST"})
		bar.Open = v[0]
		bar.High = v[1]
		bar.Low = v[2]
		bar.Close = v[3]
		bars[i] = *bar
	}
	return &bars
}

func TestAverageTrueRange_CalculateTrueRange(t *testing.T) {
	bars := *getPriceBars()
	for i, v := range bars {
		var prevbar *pricebar.PriceBar = nil
		if i > 0 {
			prevbar = &bars[i-1]
		}
		tr, _ := functions.CalculateTrueRange(&v, prevbar).Float64()
		assertions.EqualFloats(knownRanges[i], tr, t)
	}
}

func TestAverageTrueRange_Slide(t *testing.T) {
	bars := *getPriceBars()
	atr := functions.NewAverageTrueRange(bars[0:7])
	assertions.CloseEnough(1.428571, atr.GetValueFloat(), 0.000001, t)

	atr.Slide(&bars[7])
	assertions.CloseEnough(1.324489, atr.GetValueFloat(), 0.000001, t)
}
