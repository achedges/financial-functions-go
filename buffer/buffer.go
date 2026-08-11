package buffer

import (
	"financial/functions/util"

	"github.com/shopspring/decimal"
)

type Container struct {
	Values []decimal.Decimal
	Index  int
	Period int
	Length int
}

func NewContainer(values []decimal.Decimal, index int, period int, length int) Container {
	return Container{
		Values: values,
		Index:  index,
		Period: period,
		Length: length,
	}
}

func (b *Container) GetLowerBound() int {
	return (b.Index - b.Period + 1) % b.Length
}

func (b *Container) GetUpperBound(pad int) int {
	// pad is used to generate an exclusive upper bound, useful for feeding directly into a slice
	return (b.Index % b.Length) + pad
}

func (b *Container) IsRing() bool {
	return b.Period == b.Length
}

func (b *Container) Advance() {
	b.Index++
}

func (b *Container) GetLastValue() decimal.Decimal {
	return b.Values[b.GetUpperBound(0)]
}

func (b *Container) GetSum() decimal.Decimal {
	return util.Sum(b.Values[b.GetLowerBound():b.GetUpperBound(1)])
}

func (b *Container) GetSumSq() decimal.Decimal {
	squaredValues := make([]decimal.Decimal, len(b.Values))
	for i, v := range b.Values[b.GetLowerBound():b.GetUpperBound(1)] {
		squaredValues[i] = v.Mul(v)
	}
	return util.Sum(squaredValues)
}
