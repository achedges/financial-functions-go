package functions

import (
	"github.com/achedges/financial-core-go/pricebar"
	"github.com/shopspring/decimal"
)

// Directional Movement implementation

type DirectionalMovement struct {
	Period                    decimal.Decimal
	TrueRange                 decimal.Decimal
	PosDirectionalMovement    decimal.Decimal
	NegDirectionalMovement    decimal.Decimal
	TrueRangeSum              decimal.Decimal
	PosDirectionalMovementSum decimal.Decimal
	NegDirectionalMovementSum decimal.Decimal
	PosDirectionalIndex       decimal.Decimal
	NegDirectionalIndex       decimal.Decimal
	DirectionalIndex          decimal.Decimal
	DirectionalIndexSum       decimal.Decimal
	AvgDirectionalIndex       decimal.Decimal
	LastBar                   *pricebar.PriceBar
}

func NewDirectionalMovement(period int, bars []pricebar.PriceBar) *DirectionalMovement {
	p := decimal.NewFromInt32(int32(period))
	tr := decimal.Zero
	positiveDm := decimal.Zero
	negativeDm := decimal.Zero
	trueRangeSum := decimal.Zero
	positiveDmSum := decimal.Zero
	negativeDmSum := decimal.Zero
	positiveDi := decimal.Zero
	negativeDi := decimal.Zero
	dx := decimal.Zero
	dxSum := decimal.Zero
	lastBar := bars[0]

	onehundred := decimal.NewFromInt32(100)

	for i := 1; i < len(bars); i++ {
		tr = CalculateTrueRange(&bars[i], &lastBar)
		positiveDm = GetPositiveDm(&bars[i], &lastBar)
		negativeDm = GetNegativeDm(&bars[i], &lastBar)

		if i < period {
			// in the first segment
			trueRangeSum = trueRangeSum.Add(tr)
			positiveDmSum = positiveDmSum.Add(choosePositiveValue(positiveDm, negativeDm))
			negativeDmSum = negativeDmSum.Add(chooseNegativeValue(positiveDm, negativeDm))
		} else {
			// in the second segment
			trueRangeSum = GetSmoothedSum(trueRangeSum, tr, p)
			positiveDmSum = GetSmoothedSum(positiveDmSum, choosePositiveValue(positiveDm, negativeDm), p)
			negativeDmSum = GetSmoothedSum(negativeDmSum, chooseNegativeValue(positiveDm, negativeDm), p)

			positiveDi = positiveDmSum.Div(trueRangeSum).Mul(onehundred)
			negativeDi = negativeDmSum.Div(trueRangeSum).Mul(onehundred)
			dx = positiveDi.Sub(negativeDi).Abs().Div(positiveDi.Add(negativeDi)).Mul(onehundred)
			dxSum = dxSum.Add(dx)
		}

		lastBar = bars[i]
	}

	return &DirectionalMovement{
		Period:                    p,
		TrueRange:                 tr,
		PosDirectionalMovement:    positiveDm,
		NegDirectionalMovement:    negativeDm,
		TrueRangeSum:              trueRangeSum,
		PosDirectionalMovementSum: positiveDmSum,
		NegDirectionalMovementSum: negativeDmSum,
		PosDirectionalIndex:       positiveDi,
		NegDirectionalIndex:       negativeDi,
		DirectionalIndex:          dx,
		DirectionalIndexSum:       dxSum,
		LastBar:                   &lastBar,
	}
}

func choosePositiveValue(pos decimal.Decimal, neg decimal.Decimal) decimal.Decimal {
	if pos.GreaterThan(neg) {
		return pos
	}
	return decimal.Zero
}

func chooseNegativeValue(pos decimal.Decimal, neg decimal.Decimal) decimal.Decimal {
	if neg.GreaterThan(pos) {
		return neg
	}
	return decimal.Zero
}

func GetPositiveDm(currentBar *pricebar.PriceBar, previousBar *pricebar.PriceBar) decimal.Decimal {
	currentHigh := decimal.NewFromFloat(currentBar.High)
	previousHigh := decimal.NewFromFloat(previousBar.High)
	d := currentHigh.Sub(previousHigh)
	return decimal.Max(d, decimal.Zero)
}

func GetNegativeDm(currentBar *pricebar.PriceBar, previousBar *pricebar.PriceBar) decimal.Decimal {
	previousLow := decimal.NewFromFloat(previousBar.Low)
	currentLow := decimal.NewFromFloat(currentBar.Low)
	d := previousLow.Sub(currentLow)
	return decimal.Max(d, decimal.Zero)
}

func GetSmoothedSum(currentSum decimal.Decimal, newValue decimal.Decimal, period decimal.Decimal) decimal.Decimal {
	return currentSum.Sub(currentSum.Div(period)).Add(newValue)
}

func (dm *DirectionalMovement) Slide(newBar *pricebar.PriceBar) {
	dm.TrueRange = CalculateTrueRange(newBar, dm.LastBar)
	dm.PosDirectionalMovement = GetPositiveDm(newBar, dm.LastBar)
	dm.NegDirectionalMovement = GetNegativeDm(newBar, dm.LastBar)

	dm.TrueRangeSum = GetSmoothedSum(dm.TrueRangeSum, dm.TrueRange, dm.Period)
	dm.PosDirectionalMovementSum = GetSmoothedSum(dm.PosDirectionalMovementSum, choosePositiveValue(dm.PosDirectionalMovement, dm.NegDirectionalMovement), dm.Period)
	dm.NegDirectionalMovementSum = GetSmoothedSum(dm.NegDirectionalMovementSum, chooseNegativeValue(dm.PosDirectionalMovement, dm.NegDirectionalMovement), dm.Period)

	onehundred := decimal.NewFromInt32(100)
	dm.PosDirectionalIndex = dm.PosDirectionalMovementSum.Div(dm.TrueRangeSum).Mul(onehundred)
	dm.NegDirectionalIndex = dm.NegDirectionalMovementSum.Div(dm.TrueRangeSum).Mul(onehundred)

	dm.DirectionalIndex = dm.PosDirectionalIndex.Sub(dm.NegDirectionalIndex).Abs().Div(dm.PosDirectionalIndex.Add(dm.NegDirectionalIndex)).Mul(onehundred)
	dm.DirectionalIndexSum = GetSmoothedSum(dm.DirectionalIndexSum, dm.DirectionalIndex, dm.Period)
	dm.AvgDirectionalIndex = dm.DirectionalIndexSum.Div(dm.Period)

	dm.LastBar = newBar
}
