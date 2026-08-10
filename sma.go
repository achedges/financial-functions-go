package functions

import "github.com/shopspring/decimal"

// Simple Moving Average implementation

type SimpleMovingAvg struct {
	baseFunction
	valuesum decimal.Decimal
	average  decimal.Decimal
}

func NewSimpleMovingAvg(period int, prices []float64) *SimpleMovingAvg {
	buffer := NewBufferContainer(period-1, period, len(prices))
	values := make([]decimal.Decimal, len(prices))

	for i, v := range prices {
		values[i] = decimal.NewFromFloat(v)
	}

	return &SimpleMovingAvg{
		baseFunction: baseFunction{
			Values: values,
			Buffer: buffer,
		},
		valuesum: decimal.Zero,
		average:  decimal.Zero,
	}
}

func (avg *SimpleMovingAvg) Calculate() {
	avg.average = avg.valuesum.Div(decimal.NewFromInt(int64(avg.Buffer.Period)))
}

func (avg *SimpleMovingAvg) SetIndex(i int) {
	avg.Buffer.Index = i
	avg.valuesum = avg.GetBufferSum()
	avg.Calculate()
}

func (avg *SimpleMovingAvg) Slide(value float64) {
	valuedec := decimal.NewFromFloat(value)
	oldindex := avg.Buffer.GetLowerBound()
	if avg.valuesum == decimal.Zero {
		avg.SetIndex(avg.Buffer.Index)
	}

	avg.valuesum = avg.valuesum.Add(valuedec.Sub(avg.Values[oldindex]))
	avg.Calculate()

	if avg.Buffer.IsRing() {
		avg.Values[oldindex] = valuedec
	}

	avg.Buffer.Advance()
}

func (avg *SimpleMovingAvg) GetValue() decimal.Decimal {
	return avg.average
}

func (avg *SimpleMovingAvg) GetValueFloat() float64 {
	d, _ := avg.average.Float64()
	return d
}
