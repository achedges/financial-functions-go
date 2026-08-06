package util_test

import (
	"financial/functions/util"
	"testing"

	"github.com/achedges/go-assertions"
	"github.com/shopspring/decimal"
)

func TestUtil_Sum(t *testing.T) {
	ints := []int64{1, 2, 3, 4, 5}
	decimals := make([]decimal.Decimal, 5)
	for i, v := range ints {
		decimals[i] = decimal.NewFromInt(v)
	}

	assertions.EqualInts(decimal.NewFromInt(15).IntPart(), util.Sum(decimals).IntPart(), t)
}