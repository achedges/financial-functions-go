package functions

import "github.com/shopspring/decimal"

// Exponential Moving Average implementation

type ExponentialMovingAvg struct {
	baseFunction
	ema    decimal.Decimal
	price  decimal.Decimal
	weight decimal.Decimal
}

func NewExponentialMovingAvg(period int, prices []float64) *ExponentialMovingAvg {
	buffer := NewBufferContainer(period-1, period, len(prices))
	values := make([]decimal.Decimal, len(prices))

	for i, v := range prices {
		values[i] = decimal.NewFromFloat(v)
	}

	ema := &ExponentialMovingAvg{
		baseFunction: baseFunction{
			Values: values,
			Buffer: buffer,
		},
		ema:    decimal.Zero,
		price:  decimal.Zero,
		weight: decimal.NewFromFloat(2.0).DivRound(decimal.NewFromInt32(int32(period+1)), 6),
	}

	ema.SetIndex(ema.Buffer.Index)

	return ema
}

func (avg *ExponentialMovingAvg) Calculate() {
	if avg.ema == decimal.Zero {
		sum := avg.GetBufferSum()
		avg.ema = sum.DivRound(decimal.NewFromInt32(int32(avg.Buffer.Period)), 6)
	} else {
		avg.ema = avg.price.Sub(avg.ema).Mul(avg.weight).Add(avg.ema)
	}
}

func (avg *ExponentialMovingAvg) SetIndex(i int) {
	avg.Buffer.Index = i
	avg.ema = decimal.Zero
	avg.Calculate()
}

func (avg *ExponentialMovingAvg) Slide(value float64) {
	avg.price = decimal.NewFromFloat(value)
	if avg.Buffer.IsRing() {
		avg.Values[avg.Buffer.GetLowerBound()] = avg.price
	}
	avg.Calculate()
	avg.Buffer.Advance()
}

func (avg *ExponentialMovingAvg) GetValue() decimal.Decimal {
	return avg.ema
}

func (avg *ExponentialMovingAvg) GetValueFloat() float64 {
	d, _ := avg.ema.Float64()
	return d
}
