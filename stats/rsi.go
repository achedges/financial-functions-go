package stats

import (
	"github.com/achedges/financial-core-go/pricebar"
	"github.com/shopspring/decimal"
)

// Relative Strength Index implementation

type RelativeStrengthIndex struct {
	period           int
	barsWindow       []pricebar.PriceBar
	relativeStrength decimal.Decimal
	rsi              decimal.Decimal
}

func NewRelativeStrengthIndex(period int, bars []pricebar.PriceBar) *RelativeStrengthIndex {
	rsi := RelativeStrengthIndex{
		period:     period,
		barsWindow: bars,
	}

	if len(rsi.barsWindow) == period {
		rsi.Calculate()
	}

	return &rsi
}

func (rsi *RelativeStrengthIndex) RelativeStrength() decimal.Decimal {
	return rsi.relativeStrength
}

func (rsi *RelativeStrengthIndex) RSI() decimal.Decimal {
	return rsi.rsi
}

func (rsi *RelativeStrengthIndex) Calculate() {
	// This is a bit different than the original calculation.
	// For a 14-period measurement, Wilder's book includes a 15th price before calculating.
	// Since I'm not planning to use a period of 14 it seem unnecessary to adjust the backfill/period logic to account for this.

	gains := decimal.Zero
	losses := decimal.Zero

	for i := 1; i < len(rsi.barsWindow); i++ {
		curr := rsi.barsWindow[i]
		prev := rsi.barsWindow[i-1]
		if curr.Close > prev.Close {
			gains = gains.Add(decimal.NewFromFloat(curr.Close - prev.Close))
		} else if curr.Close < prev.Close {
			losses = losses.Add(decimal.NewFromFloat(prev.Close - curr.Close))
		}
	}

	n := decimal.NewFromInt32(int32(len(rsi.barsWindow)))
	onehundred := decimal.NewFromInt32(100)
	one := decimal.NewFromInt32(1)

	avgGains := gains.Div(n)
	avgLosses := losses.Div(n)

	rsi.relativeStrength = avgGains.Div(avgLosses)
	rsi.rsi = onehundred.Sub(onehundred.Div(rsi.relativeStrength.Add(one)))
}

func (rsi *RelativeStrengthIndex) Slide(bar pricebar.PriceBar) {
	if len(rsi.barsWindow) == rsi.period {
		rsi.barsWindow = rsi.barsWindow[1:] // remove first element
	}
	rsi.barsWindow = append(rsi.barsWindow, bar)
	if len(rsi.barsWindow) == rsi.period {
		rsi.Calculate()
	}
}
