package functions

import (
	"financial/functions/util"

	"github.com/shopspring/decimal"
)

type WindowFunction interface {
	GetLastValue() decimal.Decimal
	GetBufferSum() decimal.Decimal
	GetBufferSumSq() decimal.Decimal
	SetIndex(i int)
	Slide(value float64)
	Calculate()
	GetValue() decimal.Decimal
	GetValueFloat() float64
}

type baseFunction struct {
	Values []decimal.Decimal
	Buffer BufferContainer
}

func (w *baseFunction) GetLastValue() decimal.Decimal {
	return w.Values[w.Buffer.GetUpperBound(0)]
}

func (w *baseFunction) GetBufferSum() decimal.Decimal {
	return util.Sum(w.Values[w.Buffer.GetLowerBound():w.Buffer.GetUpperBound(1)])
}

func (w *baseFunction) GetBufferSumSq() decimal.Decimal {
	squaredValues := make([]decimal.Decimal, len(w.Values))
	for i, v := range w.Values[w.Buffer.GetLowerBound():w.Buffer.GetUpperBound(1)] {
		squaredValues[i] = v.Mul(v)
	}
	return util.Sum(squaredValues)
}
