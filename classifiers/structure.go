package classifiers

import (
	"financial/functions/pivots"
	"math"

	"github.com/achedges/financial-core-go/pricebar"
)

type MarketStructureClassifier struct {
	bars                 []pricebar.PriceBar
	highPivots           []int
	lowPivots            []int
	trendClassification  TrendClassification
	highPivotDiff        float64
	highPivotSlope       float64
	lowPivotDiff         float64
	lowPivotSlope        float64
	strongTrendThreshold float64
}

type MarketStructureConfig struct {
	StrongTrendThreshold float64
}

var DefaultMarketStructureConfig MarketStructureConfig = MarketStructureConfig{
	StrongTrendThreshold: 0.25,
}

func NewMarketStructure(config MarketStructureConfig) *MarketStructureClassifier {
	classifier := MarketStructureClassifier{
		strongTrendThreshold: config.StrongTrendThreshold,
	}
	return &classifier
}

func (msc *MarketStructureClassifier) GetStrongTrendThreshold() float64 {
	return msc.strongTrendThreshold
}

func (msc *MarketStructureClassifier) GetHighPivotDiff() float64 {
	return msc.highPivotDiff
}

func (msc *MarketStructureClassifier) GetLowPivotDiff() float64 {
	return msc.lowPivotDiff
}

func (msc *MarketStructureClassifier) Classify(bars []pricebar.PriceBar) TrendClassification {
	msc.bars = bars
	msc.highPivotDiff = 0.0
	msc.lowPivotDiff = 0.0

	msc.highPivots = pivots.Get(pivots.Params[pricebar.PriceBar]{
		Values:         msc.bars,
		ComparisonFunc: func(a pricebar.PriceBar, b pricebar.PriceBar) bool { return a.High >= b.High },
		DifferenceFunc: func(a pricebar.PriceBar, b pricebar.PriceBar) float64 { return a.High - b.High },
	})

	msc.lowPivots = pivots.Get(pivots.Params[pricebar.PriceBar]{
		Values:         msc.bars,
		ComparisonFunc: func(a pricebar.PriceBar, b pricebar.PriceBar) bool { return a.Low <= b.Low },
		DifferenceFunc: func(a pricebar.PriceBar, b pricebar.PriceBar) float64 { return b.Low - a.Low },
	})

	if len(msc.highPivots) <= 1 || len(msc.lowPivots) <= 1 {
		// need to set msc.trendClassification here?
		return Mixed
	}

	highestHigh := bars[0].High
	lowestHigh := highestHigh
	highestLow := bars[0].Low
	lowestLow := highestLow

	for _, b := range msc.bars {
		highestHigh = math.Max(b.High, highestHigh)
		lowestHigh = math.Min(b.High, lowestHigh)
		highestLow = math.Max(b.Low, highestLow)
		lowestLow = math.Min(b.Low, lowestLow)
	}

	highScaleDenom := highestHigh - lowestHigh
	lowScaleDenom := highestLow - lowestLow

	firstHighPivotNum := msc.bars[msc.highPivots[0]].High - lowestHigh
	lastHighPivotNum := msc.bars[msc.highPivots[len(msc.highPivots)-1]].High - lowestHigh
	firstLowPivotNum := msc.bars[msc.lowPivots[0]].Low - lowestLow
	lastLowPivotNum := msc.bars[msc.lowPivots[len(msc.lowPivots)-1]].Low - lowestLow

	msc.highPivotDiff = lastHighPivotNum - firstHighPivotNum
	msc.lowPivotDiff = lastLowPivotNum - firstLowPivotNum

	highPivotDiffScaled := msc.highPivotDiff / highScaleDenom
	lowPivotDiffScaled := msc.lowPivotDiff / lowScaleDenom

	highTimeScaled := float64(msc.highPivots[len(msc.highPivots)-1]-msc.highPivots[0]) / float64(len(msc.bars)-1)
	lowTimeScaled := float64(msc.lowPivots[len(msc.lowPivots)-1]-msc.lowPivots[0]) / float64(len(msc.bars)-1)

	msc.highPivotSlope = highPivotDiffScaled / highTimeScaled
	msc.lowPivotSlope = lowPivotDiffScaled / lowTimeScaled

	if msc.highPivotSlope >= msc.strongTrendThreshold && msc.lowPivotSlope >= msc.strongTrendThreshold {
		msc.trendClassification = StrongUp
	} else if msc.highPivotSlope > 0.0 && msc.lowPivotSlope > 0.0 {
		msc.trendClassification = WeakUp
	} else if msc.highPivotSlope <= -msc.strongTrendThreshold && msc.lowPivotSlope <= -msc.strongTrendThreshold {
		msc.trendClassification = StrongDown
	} else if msc.highPivotSlope < 0.0 && msc.lowPivotSlope < 0.0 {
		msc.trendClassification = WeakDown
	} else {
		msc.trendClassification = Mixed
	}

	return msc.trendClassification
}

func (msc *MarketStructureClassifier) getLastPivot(pivots []int) *pricebar.PriceBar {
	if pivots == nil || len(pivots) == 0 {
		return nil
	}

	lastPivot := pivots[len(pivots)-1]
	if msc.bars == nil || lastPivot >= len(msc.bars) {
		return nil
	}

	return &msc.bars[lastPivot]
}

func (msc *MarketStructureClassifier) GetLastHighPivot() *pricebar.PriceBar {
	return msc.getLastPivot(msc.highPivots)
}

func (msc *MarketStructureClassifier) GetLastLowPivot() *pricebar.PriceBar {
	return msc.getLastPivot(msc.lowPivots)
}

func (msc *MarketStructureClassifier) GetLastHighPivotIndex() int {
	if msc.highPivots == nil || len(msc.highPivots) == 0 {
		return -1
	}
	return msc.highPivots[len(msc.highPivots)-1]
}

func (msc *MarketStructureClassifier) GetLastLowPivotIndex() int {
	if msc.lowPivots == nil || len(msc.lowPivots) == 0 {
		return -1
	}
	return msc.lowPivots[len(msc.lowPivots)-1]
}
