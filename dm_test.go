package functions_test

import (
	"financial/functions"
	"testing"

	"github.com/achedges/financial-core-go/pricebar"
	"github.com/achedges/go-assertions"
)

func getDirectionalPriceBars(source [][]float64) []pricebar.PriceBar {
	bars := make([]pricebar.PriceBar, len(source))
	for i, v := range source {
		barConfig := pricebar.Config{
			Symbol: "TEST",
		}
		bar := pricebar.New(barConfig)
		bar.High = v[0]
		bar.Low = v[1]
		bar.Close = v[2]
		bars[i] = *bar
	}
	return bars
}

func TestDirectionalMovement_Init(t *testing.T) {
	bars := getDirectionalPriceBars(initPrices)
	dm := functions.NewDirectionalMovement(14, bars)
	assertions.CloseEnough(57.05, dm.TrueRangeSum.InexactFloat64(), 0.01, t)
	assertions.CloseEnough(12.29, dm.PosDirectionalMovementSum.InexactFloat64(), 0.01, t)
	assertions.CloseEnough(21.75, dm.NegDirectionalMovementSum.InexactFloat64(), 0.01, t)
}

func TestDirectionalMovement_Slide(t *testing.T) {
	bars := getDirectionalPriceBars(initPrices)
	slideBars := getDirectionalPriceBars(slidePrices)

	dm := functions.NewDirectionalMovement(14, bars)
	for i, v := range slideBars {
		dm.Slide(&v)
		z := i + (14 * 2) - 1
		assertions.CloseEnough(known_tr[z], dm.TrueRange.InexactFloat64(), 0.001, t)
		assertions.CloseEnough(known_pdm[z], dm.PosDirectionalMovement.InexactFloat64(), 0.001, t)
		assertions.CloseEnough(known_ndm[z], dm.NegDirectionalMovement.InexactFloat64(), 0.001, t)
		assertions.CloseEnough(known_tr14[z], dm.TrueRangeSum.InexactFloat64(), 0.1, t)
		assertions.CloseEnough(known_pdm14[z], dm.PosDirectionalMovementSum.InexactFloat64(), 0.1, t)
		assertions.CloseEnough(known_ndm14[z], dm.NegDirectionalMovementSum.InexactFloat64(), 0.1, t)
		assertions.CloseEnough(known_pdi[z], dm.PosDirectionalIndex.InexactFloat64(), 0.1, t)
		assertions.CloseEnough(known_ndi[z], dm.NegDirectionalIndex.InexactFloat64(), 0.1, t)
		assertions.CloseEnough(known_dx[z], dm.DirectionalIndex.InexactFloat64(), 0.1, t)
		assertions.CloseEnough(known_adx[z], dm.AvgDirectionalIndex.InexactFloat64(), 0.1, t)
	}
}

var initPrices = [][]float64{
	{274.00, 272.00, 272.75},
	{273.25, 270.25, 270.75},
	{272.00, 269.75, 270.00},
	{270.75, 268.00, 269.25},
	{270.00, 269.00, 269.75},
	{270.50, 268.00, 270.00},
	{268.50, 266.50, 266.50},
	{265.50, 263.00, 263.25},
	{262.50, 259.00, 260.25},
	{263.50, 260.00, 263.00},
	{269.50, 263.00, 266.50},
	{267.25, 265.00, 267.00},
	{267.50, 265.50, 265.75},
	{269.75, 266.00, 268.50},
	{268.25, 263.25, 264.25},
	{264.00, 261.50, 264.00},
	{268.00, 266.25, 266.50},
	{266.00, 264.25, 265.25},
	{274.00, 267.00, 273.00},
	{277.50, 273.50, 276.75},
	{277.00, 272.50, 273.00},
	{272.00, 269.50, 270.25},
	{267.75, 264.00, 266.75},
	{269.25, 263.00, 263.00},
	{266.00, 263.50, 265.50},
	{265.00, 262.00, 262.25},
	{264.75, 261.50, 262.75},
	{261.00, 255.50, 255.50},
}

var slidePrices = [][]float64{
	{257.50, 253.00, 253.00},
	{259.00, 254.00, 257.50},
	{259.75, 257.50, 257.50},
	{257.25, 250.00, 250.00},
	{250.00, 247.00, 249.75},
	{254.25, 252.75, 253.75},
	{254.00, 250.50, 251.25},
	{253.25, 250.25, 250.50},
	{253.25, 251.00, 253.00},
	{251.75, 250.50, 251.50},
	{253.00, 249.50, 250.00},
	{251.50, 245.25, 245.75},
	{246.25, 240.00, 242.75},
	{244.25, 241.25, 243.50},
}

var known_tr = []float64{3.0, 2.25, 2.75, 1.0, 2.5, 3.5, 3.5, 4.25, 3.5, 6.5, 2.25, 2.0, 4.0, 5.25, 2.75, 4.0, 2.25, 8.75, 4.5, 4.5, 3.5, 6.25, 6.25, 3.0, 3.5, 3.25, 7.25, 4.5, 6.0, 2.25, 7.5, 3.0, 4.5, 3.5, 3.0, 2.75, 2.5, 3.5, 6.25, 6.25, 3.0}
var known_pdm = []float64{0.0, 0.0, 0.0, 0.0, 0.5, 0.0, 0.0, 0.0, 1.0, 6.0, 0.0, 0.25, 2.25, 0.0, 0.0, 4.0, 0.0, 8.0, 3.5, 0.0, 0.0, 0.0, 1.5, 0.0, 0.0, 0.0, 0.0, 0.0, 1.5, 0.75, 0.0, 0.0, 4.25, 0.0, 0.0, 0.0, 0.0, 1.25, 0.0, 0.0, 0.0}
var known_ndm = []float64{1.75, 0.5, 1.75, 0.0, 1.0, 1.5, 3.5, 4.0, 0.0, 0.0, 0.0, 0.0, 0.0, 2.75, 1.75, 0.0, 2.0, 0.0, 0.0, 1.0, 3.0, 5.5, 1.0, 0.0, 1.5, 0.5, 6.0, 2.5, 0.0, 0.0, 7.5, 3.0, 0.0, 2.25, 0.25, 0.0, 0.5, 1.0, 4.25, 5.25, 0.0}
var known_tr14 = []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 43.32, 42.98, 43.91, 43.02, 48.70, 49.72, 50.67, 50.55, 53.19, 55.64, 54.67, 54.26, 53.63, 57.05, 57.47, 59.36, 57.37, 60.77, 59.43, 59.68, 58.92, 57.71, 56.34, 54.82, 54.40, 56.76, 58.96, 57.75}
var known_pdm14 = []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 8.82, 8.19, 11.61, 10.78, 18.01, 20.22, 18.78, 17.44, 16.19, 16.53, 15.35, 14.26, 13.24, 12.29, 11.40, 12.09, 11.98, 11.12, 10.33, 13.84, 12.85, 11.93, 11.08, 10.29, 10.80, 10.03, 9.31, 8.64}
var known_ndm14 = []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 15.75, 16.38, 15.21, 16.12, 14.97, 13.90, 13.91, 15.91, 20.28, 18.83, 17.48, 17.73, 16.97, 21.76, 22.70, 21.08, 19.57, 25.67, 26.84, 24.92, 25.39, 23.83, 22.13, 21.05, 19.55, 22.40, 26.05, 24.19}
var known_pdi = []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 20, 19, 26, 25, 37, 41, 37, 34, 30, 30, 28, 26, 25, 22, 19.85, 20.37, 20.88, 18.3, 17.38, 23.19, 21.81, 20.68, 19.67, 18.77, 19.86, 17.67, 15.80, 15}
var known_ndi = []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 36, 38, 35, 37, 31, 28, 27, 31, 38, 34, 32, 33, 32, 38, 39.49, 35.50, 34.11, 42.24, 45.15, 41.75, 43, 41.28, 39.27, 38.39, 36, 39.45, 44.17, 41.88}
var known_dx = []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 28, 33, 13, 20, 9, 19, 15, 5, 11, 6, 6, 11, 12, 28, 33.08, 27.06, 24, 39.52, 44.40, 28.57, 32.77, 33.24, 33.24, 34.31, 28.78, 38.11, 47.29, 47.29}
var known_adx = []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 16.76, 17.50, 17.97, 19.51, 21.29, 21.81, 22.59, 23.35, 24.10, 24.73, 25.04, 26.00, 27.53, 28.94}
