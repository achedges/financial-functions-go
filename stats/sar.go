package stats

import (
	"github.com/achedges/financial-core-go/candle"
	"github.com/shopspring/decimal"
)

type ParabolicSAR struct {
	isLong        bool
	reverseSignal bool
	high          decimal.Decimal
	low           decimal.Decimal
	stop          decimal.Decimal
	rate          decimal.Decimal
	alpha         decimal.Decimal
	maxAlpha      decimal.Decimal
}

func NewParabolicSAR(candle candle.Candle) *ParabolicSAR {
	psar := ParabolicSAR{
		isLong:        candle.IsUp(),
		reverseSignal: false,
		high:          decimal.NewFromFloat(candle.High),
		low:           decimal.NewFromFloat(candle.Low),
		stop:          decimal.Zero,
		rate:          decimal.NewFromFloat(0.02),
		maxAlpha:      decimal.NewFromFloat(0.2),
	}

	psar.alpha = psar.rate
	if psar.isLong {
		psar.stop = psar.low
	} else {
		psar.stop = psar.high
	}

	return &psar
}

func (psar *ParabolicSAR) updateAlpha() {
	psar.alpha = psar.alpha.Add(psar.rate)
	psar.alpha = decimal.Min(psar.alpha, psar.maxAlpha)
}

func (psar *ParabolicSAR) reverse(newStop decimal.Decimal) {
	psar.isLong = !psar.isLong
	psar.stop = newStop
	psar.alpha = psar.rate
	psar.reverseSignal = true
}

func (psar *ParabolicSAR) SetIsLong(isLong bool) {
	psar.isLong = isLong
}

func (psar *ParabolicSAR) Slide(newCandle candle.Candle) {
	candleLow := decimal.NewFromFloat(newCandle.Low)
	candleHigh := decimal.NewFromFloat(newCandle.High)

	// important to follow these steps in order:
	// - check for breach
	// - check for new high/low (extreme)
	// - calculate new stop
	// - update alpha (if new high/low)

	newHighLow := false

	if psar.isLong {
		if candleLow.LessThan(psar.stop) {
			psar.reverse(decimal.Max(psar.high, candleHigh))
			psar.low = candleLow
		} else {
			if candleHigh.GreaterThan(psar.high) {
				psar.high = candleHigh
				newHighLow = true
			}
			psar.stop = psar.stop.Add(psar.alpha.Mul(psar.high.Sub(psar.stop)))
			psar.reverseSignal = false
		}
	} else {
		if candleHigh.GreaterThan(psar.stop) {
			psar.reverse(decimal.Min(psar.low, candleLow))
			psar.high = candleHigh
		} else {
			if candleLow.LessThan(psar.low) {
				psar.low = candleLow
				newHighLow = true
			}
			psar.stop = psar.stop.Sub(psar.alpha.Mul(psar.stop.Sub(psar.low)))
			psar.reverseSignal = false
		}
	}

	if newHighLow {
		psar.updateAlpha()
	}
}
