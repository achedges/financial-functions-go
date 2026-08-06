package util

import "github.com/shopspring/decimal"

func Sum(values []decimal.Decimal) decimal.Decimal {
	sum := decimal.Zero
	for _, v := range values {
		sum = sum.Add(v)
	}
	return sum
}
