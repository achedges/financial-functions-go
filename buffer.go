package functions

import (
	"financial/functions/util"

	"github.com/shopspring/decimal"
)

type BufferContainer struct {
	Values []decimal.Decimal
	Index  int
	Period int
	Length int
}

func NewBufferContainer(values []decimal.Decimal, index int, period int, length int) BufferContainer {
	return BufferContainer{
		Values: values,
		Index:  index,
		Period: period,
		Length: length,
	}
}

func (b *BufferContainer) GetLowerBound() int {
	return (b.Index - b.Period + 1) % b.Length
}

func (b *BufferContainer) GetUpperBound(pad int) int {
	// pad is used to generate an exclusive upper bound, useful for feeding directly into a slice
	return (b.Index % b.Length) + pad
}

func (b *BufferContainer) IsRing() bool {
	return b.Period == b.Length
}

func (b *BufferContainer) Advance() {
	b.Index++
}

func (b *BufferContainer) GetLastValue() decimal.Decimal {
	return b.Values[b.GetUpperBound(0)]
}

func (b *BufferContainer) GetSum() decimal.Decimal {
	return util.Sum(b.Values[b.GetLowerBound():b.GetUpperBound(1)])
}

func (b *BufferContainer) GetSumSq() decimal.Decimal {
	squaredValues := make([]decimal.Decimal, len(b.Values))
	for i, v := range b.Values[b.GetLowerBound():b.GetUpperBound(1)] {
		squaredValues[i] = v.Mul(v)
	}
	return util.Sum(squaredValues)
}
