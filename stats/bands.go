package stats

import (
	"financial/functions/buffer"

	"github.com/shopspring/decimal"
)

// Bollinger Bands implementation

type BollingerBands struct {
	Buffer            buffer.Container
	standardDeviation *StandardDeviation
	UpperBand         decimal.Decimal
	LowerBand         decimal.Decimal
}

func NewBollingerBands(period int, prices []float64) *BollingerBands {
	values := make([]decimal.Decimal, len(prices))
	for i, v := range prices {
		values[i] = decimal.NewFromFloat(v)
	}

	stdev := NewStandardDeviation(period, prices)

	bb := &BollingerBands{
		Buffer:            buffer.NewContainer(values, period-1, period, len(values)),
		UpperBand:         decimal.Zero,
		LowerBand:         decimal.Zero,
		standardDeviation: stdev,
	}

	return bb
}

func (bb *BollingerBands) SetIndex(i int) {
	bb.standardDeviation.SetIndex(i)
	bb.Buffer.Index = i
}

func (bb *BollingerBands) Calculate() {
	mavg := bb.standardDeviation.movingAverage.GetValue()
	stdev := bb.standardDeviation.GetValue()
	mult := decimal.NewFromInt(2)
	bb.UpperBand = mavg.Add(stdev.Mul(mult))
	bb.LowerBand = mavg.Sub(stdev.Mul(mult))
}

func (bb *BollingerBands) Slide(value float64) {
	oldindex := bb.Buffer.GetLowerBound()
	bb.standardDeviation.Slide(value)
	bb.Calculate()

	if bb.Buffer.IsRing() {
		bb.Buffer.Values[oldindex] = decimal.NewFromFloat(value)
	}

	bb.Buffer.Advance()
}

func (bb *BollingerBands) GetValue() decimal.Decimal {
	return bb.standardDeviation.GetValue()
}

func (bb *BollingerBands) GetValueFloat() float64 {
	d, _ := bb.standardDeviation.GetValue().Float64()
	return d
}

func (bb *BollingerBands) GetMovingAvgValue() decimal.Decimal {
	return bb.standardDeviation.movingAverage.GetValue()
}

func (bb *BollingerBands) GetMovingAvgValueFloat() float64 {
	d, _ := bb.standardDeviation.movingAverage.GetValue().Float64()
	return d
}
