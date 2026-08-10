package functions

import (
	"github.com/achedges/financial-core-go/pricebar"
	"github.com/shopspring/decimal"
)

// Average True Range implementation

type AverageTrueRange struct {
	n                decimal.Decimal
	averageTrueRange decimal.Decimal
	currentBar       pricebar.PriceBar
	previousBar      pricebar.PriceBar
}

func NewAverageTrueRange(bars []pricebar.PriceBar) *AverageTrueRange {
	atr := AverageTrueRange{
		n: decimal.NewFromInt(int64(len(bars))),
	}

	numbars := atr.n.IntPart()
	for i := range numbars {
		var prev *pricebar.PriceBar = nil
		if i > 0 {
			prev = &bars[i-1]
		}
		atr.averageTrueRange = atr.averageTrueRange.Add(CalculateTrueRange(&bars[i], prev))
	}

	atr.averageTrueRange = atr.averageTrueRange.Div(atr.n)
	atr.currentBar = bars[numbars-1]
	atr.previousBar = bars[numbars-2]

	return &atr
}

func (atr *AverageTrueRange) Slide(newBar *pricebar.PriceBar) {
	atr.previousBar = atr.currentBar
	atr.currentBar = *newBar
	trueRange := CalculateTrueRange(&atr.currentBar, &atr.previousBar)
	decay := atr.averageTrueRange.Mul(atr.n.Sub(decimal.NewFromInt(1)))
	atr.averageTrueRange = decay.Add(trueRange).Div(atr.n)
}

func CalculateTrueRange(currentBar *pricebar.PriceBar, previousBar *pricebar.PriceBar) decimal.Decimal {
	currentHigh := decimal.NewFromFloat(currentBar.High)
	currentLow := decimal.NewFromFloat(currentBar.Low)
	currentDiff := currentHigh.Sub(currentLow)

	if previousBar == nil {
		return currentDiff
	}

	previousClose := decimal.NewFromFloat(previousBar.Close)
	previousHighDiff := currentHigh.Sub(previousClose).Abs()
	previousLowDiff := currentLow.Sub(previousClose).Abs()

	// return max of currentDiff, previousHighDiff, and previousLowDiff
	maxPreviousDiff := decimal.Zero
	if previousHighDiff.GreaterThan(previousLowDiff) {
		maxPreviousDiff = previousHighDiff
	} else {
		maxPreviousDiff = previousLowDiff
	}

	if currentDiff.GreaterThan(maxPreviousDiff) {
		return currentDiff
	}

	return maxPreviousDiff
}

func (atr *AverageTrueRange) GetValue() decimal.Decimal {
	return atr.averageTrueRange
}

func (atr *AverageTrueRange) GetValueFloat() float64 {
	d, _ := atr.averageTrueRange.Float64()
	return d
}
