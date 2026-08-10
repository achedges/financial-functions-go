package pivots_test

import (
	"financial/functions/pivots"
	"testing"

	"github.com/achedges/go-assertions"
)

func TestPivots_FindAllIndexes(t *testing.T) {
	highPivots := pivots.FindAllIndexes(pivotPrices1, func(l float64, r float64) bool { return l >= r })
	assertions.False(len(highPivots) == 0, t)
	assertions.EqualInts(len(expectedHighPivots1), len(highPivots), t)
	for i, exp := range expectedHighPivots1 {
		assertions.EqualInts(exp, highPivots[i], t)
	}

	lowPivots := pivots.FindAllIndexes(pivotPrices1, func(l float64, r float64) bool { return l <= r })
	assertions.False(len(lowPivots) == 0, t)
	assertions.EqualInts(len(expectedLowPivots1), len(lowPivots), t)
	for i, exp := range expectedLowPivots1 {
		assertions.EqualInts(exp, lowPivots[i], t)
	}

	// second example from 20260428 0246-0345 (highs only)
	highPivots2 := pivots.FindAllIndexes(pivotPrices2, func(l float64, r float64) bool { return l >= r })
	assertions.False(len(highPivots2) == 0, t)
	assertions.EqualInts(len(expectedHighPivots2), len(highPivots2), t)
	for i, exp := range expectedHighPivots2 {
		assertions.EqualInts(exp, highPivots2[i], t)
	}
}

func TestPivots_CalculateDispersions(t *testing.T) {
	span := len(pivotPrices1) / pivots.DefaultDispersionSpan

	highPivots := pivots.FindAllIndexes(pivotPrices1, func(l float64, r float64) bool { return l >= r })
	highDispersions := pivots.CalculateDispersions(pivotPrices1, highPivots, span, func(l float64, r float64) float64 { return l - r })
	assertions.False(len(highDispersions) == 0, t)
	assertions.EqualInts(len(expectedHighDispersions1), len(highDispersions), t)
	for i, exp := range expectedHighDispersions1 {
		assertions.CloseEnough(exp, highDispersions[i], 0.00001, t)
	}

	lowPivots := pivots.FindAllIndexes(pivotPrices1, func(l float64, r float64) bool { return l <= r })
	lowDispersions := pivots.CalculateDispersions(pivotPrices1, lowPivots, span, func(l float64, r float64) float64 { return r - l })
	assertions.False(len(lowDispersions) == 0, t)
	assertions.EqualInts(len(expectedLowDispersions1), len(lowDispersions), t)
	for i, exp := range expectedLowDispersions1 {
		assertions.CloseEnough(exp, lowDispersions[i], 0.00001, t)
	}

	// second example from 20260428 0246-0345 (highs only)
	highPivots2 := pivots.FindAllIndexes(pivotPrices2, func(l float64, r float64) bool { return l >= r })
	highDispersions2 := pivots.CalculateDispersions(pivotPrices2, highPivots2, span, func(l float64, r float64) float64 { return l - r })
	assertions.False(len(highDispersions2) == 0, t)
	assertions.EqualInts(len(expectedHighDispersions2), len(highDispersions2), t)
	for i, exp := range expectedHighDispersions2 {
		assertions.CloseEnough(exp, highDispersions2[i], 0.00001, t)
	}
}

func TestPivots_Consolidate(t *testing.T) {
	dSpan := len(pivotPrices1) / pivots.DefaultDispersionSpan
	cSpan := len(pivotPrices1) / pivots.DefaultConsolidationSpan

	highPivots := pivots.FindAllIndexes(pivotPrices1, func(l float64, r float64) bool { return l >= r })
	highDispersions := pivots.CalculateDispersions(pivotPrices1, highPivots, dSpan, func(l float64, r float64) float64 { return l - r })

	lowPivots := pivots.FindAllIndexes(pivotPrices1, func(l float64, r float64) bool { return l <= r })
	lowDispersions := pivots.CalculateDispersions(pivotPrices1, lowPivots, dSpan, func(l float64, r float64) float64 { return r - l })

	consolidatedHighs := pivots.Consolidate(highPivots, highDispersions, cSpan)
	assertions.EqualInts(len(expectedHighConsolidations1), len(consolidatedHighs), t)
	for i, exp := range expectedHighConsolidations1 {
		assertions.EqualInts(exp, consolidatedHighs[i], t)
	}

	consolidatedLows := pivots.Consolidate(lowPivots, lowDispersions, cSpan)
	assertions.EqualInts(len(expectedLowConsolidations1), len(consolidatedLows), t)
	for i, exp := range expectedLowConsolidations1 {
		assertions.EqualInts(exp, (consolidatedLows)[i], t)
	}

	// second example from 20260428 0246-0345 (highs only)
	highPivots2 := pivots.FindAllIndexes(pivotPrices2, func(l float64, r float64) bool { return l >= r })
	highDispersions2 := pivots.CalculateDispersions(pivotPrices2, highPivots2, dSpan, func(l float64, r float64) float64 { return l - r })
	consolidatedHighs2 := pivots.Consolidate(highPivots2, highDispersions2, cSpan)
	assertions.EqualInts(len(expectedHighConsolidations2), len(consolidatedHighs2), t)
	for i, exp := range expectedHighConsolidations2 {
		assertions.EqualInts(exp, consolidatedHighs2[i], t)
	}
}

func TestPivots_Get(t *testing.T) {
	expectedPivots := []int{1, 7, 12}

	defaultPivots := pivots.Get(pivots.Params[float64]{
		Values:         data,
		ComparisonFunc: func(l float64, r float64) bool { return l > r },
		DifferenceFunc: func(l float64, r float64) float64 { return l - r },
	})
	assertions.EqualInts(len(expectedPivots), len(defaultPivots), t)
	for i, exp := range expectedPivots {
		assertions.EqualInts(exp, defaultPivots[i], t)
	}

	overridePivots := pivots.Get(pivots.Params[float64]{
		Values:            data,
		ComparisonFunc:    func(l float64, r float64) bool { return l > r },
		DifferenceFunc:    func(l float64, r float64) float64 { return l - r },
		DispersionSpan:    2,
		ConsolidationSpan: 3,
	})
	assertions.EqualInts(len(expectedPivots), len(overridePivots), t)
	for i, exp := range expectedPivots {
		assertions.EqualInts(exp, overridePivots[i], t)
	}
}
