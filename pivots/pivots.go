package pivots

import (
	"math"
)

const DefaultDispersionSpan = 4
const DefaultConsolidationSpan = 4

type Params[T any] struct {
	Values            *[]T
	ComparisonFunc    func(l T, r T) bool
	DifferenceFunc    func(l T, r T) float64
	ConsolidationSpan int
	DispersionSpan    int
}

func FindAllIndexes[T any](values *[]T, comparison func(a T, b T) bool) *[]int {
	// Find all elements where the provided comparison lambda is true for all non-null neighbors, and return their indices
	n := len(*values)
	pivots := make([]int, 0, n/2) // start with half the size of the input

	var c *T
	var l *T
	var r *T

	for i := 0; i < n; i++ {
		c = &(*values)[i]

		if i == 0 {
			l = nil
		} else {
			l = &(*values)[i-1]
		}

		if i == n-1 {
			r = nil
		} else {
			r = &(*values)[i+1]
		}

		if l == nil && r == nil {
			return &pivots
		} else if l != nil && r != nil {
			if comparison(*c, *l) && comparison(*c, *r) {
				pivots = append(pivots, i)
			}
		} else if l == nil && comparison(*c, *r) {
			pivots = append(pivots, i)
		} else if r == nil && comparison(*c, *l) {
			pivots = append(pivots, i)
		}
	}

	return &pivots
}

func CalculateDispersions[T any](values *[]T, pivotIndexes *[]int, span int, difference func(a T, b T) float64) *[]float64 {
	// Given a list of values, and a corresponding list of pivots (from getPivots()), find the dispersion of each pivot point
	// by calculating the maximum difference between that pivot point and the neighboring elements, up to some maximum span.
	n := len(*values)
	dispersions := make([]float64, len(*pivotIndexes))

	for i, pivot := range *pivotIndexes {
		dispersion := 0.0
		peak := (*values)[pivot]
		for s := 1; s <= span; s++ {
			leftspan := 0.0
			rightspan := 0.0

			if pivot-s >= 0 {
				leftspan = difference(peak, (*values)[pivot-s])
			}
			if pivot+s < n {
				rightspan = difference(peak, (*values)[pivot+s])
			}

			maxspan := math.Max(leftspan, rightspan)
			if maxspan > dispersion {
				dispersion = maxspan
			}
		}

		dispersions[i] = dispersion
	}

	return &dispersions
}

func Consolidate(pivotIndexes *[]int, dispersions *[]float64, span int) *[]int {
	// Given a list of pivot indices (from getPivots()) and dispersions (from getDispersions()), produce a consolidated
	// list of pivot indices by collapsing neighboring pivot points into the pivot point with the highest dispersion,
	// up to some maximum span.
	n := len(*pivotIndexes)
	results := make([]int, 0, n)
	if len(*pivotIndexes) != len(*dispersions) {
		return &results
	}

	consolidated := make([]int, n)
	copy(consolidated, *pivotIndexes)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if (*pivotIndexes)[j]-(*pivotIndexes)[i] > span {
				break
			}

			if consolidated[j] == -1 {
				continue
			}

			if (*dispersions)[i] > (*dispersions)[j] && consolidated[i] != -1 {
				consolidated[j] = -1
			} else if consolidated[j] != -1 {
				consolidated[i] = -1
				break
			}
		}
	}

	// filter results to return all indexes >= 0
	for _, v := range consolidated {
		if v >= 0 {
			results = append(results, v)
		}
	}

	return &results
}

//func Get[T any](params Params[T]) *[]int {
//	//if params.Values == nil {
//	//	return make(*[]int, 0)
//	//}
//
//	if params.DispersionSpan == 0 {
//		params.DispersionSpan = len(*params.Values) / DefaultDispersionSpan
//	}
//	if params.ConsolidationSpan == 0 {
//		params.ConsolidationSpan = len(*params.Values) / DefaultConsolidationSpan
//	}
//
//	rawPivots := FindAllIndexes(params.Values, params.ComparisonFunc)
//	dispersions := CalculateDispersions(params.Values, rawPivots, params.DispersionSpan, params.DifferenceFunc)
//
//}
