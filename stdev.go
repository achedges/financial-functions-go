package functions

import "github.com/shopspring/decimal"

// Standard Deviation implementation

type StandardDeviation struct {
	baseFunction
	valuesSqSum   decimal.Decimal
	deviation     decimal.Decimal
	movingAverage *SimpleMovingAvg
}

func NewStandardDeviation(period int, prices []float64) *StandardDeviation {
	buffer := NewBufferContainer(period-1, period, len(prices))
	values := make([]decimal.Decimal, len(prices))

	for i, v := range prices {
		values[i] = decimal.NewFromFloat(v)
	}

	movingAverage := NewSimpleMovingAvg(period, prices)

	stdev := &StandardDeviation{
		baseFunction: baseFunction{
			Values: values,
			Buffer: buffer,
		},
		valuesSqSum:   decimal.Zero,
		deviation:     decimal.Zero,
		movingAverage: movingAverage,
	}

	stdev.SetIndex(stdev.Buffer.Index)

	return stdev
}

func (stdev *StandardDeviation) Calculate() {
	mavg := stdev.movingAverage.GetValue()
	p := decimal.NewFromInt32(int32(stdev.Buffer.Period))
	variance := stdev.valuesSqSum.Sub(mavg.Pow(decimal.NewFromFloat(2.0)).Mul(p))
	stdev.deviation = variance.Abs().Div(p).Pow(decimal.NewFromFloat32(0.5))
}

func (stdev *StandardDeviation) SetIndex(i int) {
	stdev.movingAverage.SetIndex(i)
	stdev.Buffer.Index = i
	stdev.valuesSqSum = stdev.GetBufferSumSq()
	stdev.Calculate()
}

func (stdev *StandardDeviation) Slide(value float64) {
	stdev.movingAverage.Slide(value)
	oldindex := stdev.Buffer.GetLowerBound()

	if stdev.valuesSqSum == decimal.Zero {
		stdev.SetIndex(stdev.Buffer.Index)
	}

	valued := decimal.NewFromFloat(value)
	two := decimal.NewFromInt32(2)
	stdev.valuesSqSum = stdev.valuesSqSum.Add(valued.Pow(two).Sub(stdev.Values[oldindex].Pow(two)))
	stdev.Calculate()

	if stdev.Buffer.IsRing() {
		stdev.Values[oldindex] = valued
	}

	stdev.Buffer.Advance()
}

func (stdev *StandardDeviation) GetValue() decimal.Decimal {
	return stdev.deviation
}

func (stdev *StandardDeviation) GetValueFloat() float64 {
	d, _ := stdev.deviation.Float64()
	return d
}
