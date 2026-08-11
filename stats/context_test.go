package stats_test

import "github.com/achedges/financial-core-go/pricebar"

const period int = 4

var data = []float64{303, 305.12, 303.17, 302.09, 306.17, 309.01, 310.24, 311.26, 310.35, 303.26, 300.42, 306.41, 307.4, 303.5, 305.57, 300.8}

var sma = []float64{303.3450, 304.1375, 305.1100, 306.8775, 309.1700, 310.2150, 308.7775, 306.3225, 305.1100, 304.3725, 304.4325, 305.7200, 304.3175, 303.2175, 303.6225, 303.0225}
var ema = []float64{303.3450, 304.4750, 306.2890, 307.8694, 309.2256, 309.6754, 307.1092, 304.4335, 305.2241, 306.0945, 305.0567, 305.2620, 303.4772, 303.2863, 304.0198, 303.6799}
var dev = []float64{1.1040, 1.5988, 2.7027, 3.1335, 1.9065, 0.8005, 3.2100, 4.6048, 3.6934, 2.7467, 2.7243, 1.4361, 2.4552, 1.6960, 1.8967, 1.5298}

func getPriceBarList(count int) []*pricebar.PriceBar {
	bars := make([]*pricebar.PriceBar, count)
	for i := range count {
		config := pricebar.Config{
			Symbol:     "TEST",
			BasisPrice: data[i],
		}
		bars[i] = pricebar.New(config)
	}
	return bars
}

func mapClosePricesFromBars(bars []*pricebar.PriceBar) []float64 {
	floats := make([]float64, len(bars))
	for i := range bars {
		floats[i] = bars[i].Close
	}
	return floats
}
