package classifiers

import (
	"financial/functions/pivots"

	"github.com/achedges/financial-core-go/candle"
)

type MarketStructureClassifier struct {
	bars                 []candle.Candle
	highPivots           []int
	lowPivots            []int
	trendClassification  TrendClassification
	highPivotDiff        float64
	lowPivotDiff         float64
	strongTrendThreshold float64
	thresholdQualifier   ThresholdQualifier
}

type MarketStructureConfig struct {
	StrongTrendThreshold float64
	ThresholdQualifier   ThresholdQualifier
}

var DefaultMarketStructureConfig = MarketStructureConfig{
	StrongTrendThreshold: 0.25, // 1.0 is a 45-degree slope
	ThresholdQualifier:   SlopeDouble,
}

func NewMarketStructure(config MarketStructureConfig) *MarketStructureClassifier {
	classifier := MarketStructureClassifier{
		strongTrendThreshold: config.StrongTrendThreshold,
		thresholdQualifier:   config.ThresholdQualifier,
	}
	return &classifier
}

func (msc *MarketStructureClassifier) GetStrongTrendThreshold() float64 {
	return msc.strongTrendThreshold
}

func (msc *MarketStructureClassifier) GetThresholdQualifier() ThresholdQualifier {
	return msc.thresholdQualifier
}

func (msc *MarketStructureClassifier) GetHighPivotDiff() float64 {
	return msc.highPivotDiff
}

func (msc *MarketStructureClassifier) GetLowPivotDiff() float64 {
	return msc.lowPivotDiff
}

func (msc *MarketStructureClassifier) Classify(bars []candle.Candle) TrendClassification {
	msc.bars = bars
	msc.highPivotDiff = 0.0
	msc.lowPivotDiff = 0.0

	msc.highPivots = pivots.Get(pivots.Params[candle.Candle]{
		Values:         msc.bars,
		ComparisonFunc: func(a candle.Candle, b candle.Candle) bool { return a.High >= b.High },
		DifferenceFunc: func(a candle.Candle, b candle.Candle) float64 { return a.High - b.High },
	})

	msc.lowPivots = pivots.Get(pivots.Params[candle.Candle]{
		Values:         msc.bars,
		ComparisonFunc: func(a candle.Candle, b candle.Candle) bool { return a.Low <= b.Low },
		DifferenceFunc: func(a candle.Candle, b candle.Candle) float64 { return b.Low - a.Low },
	})

	if len(msc.highPivots) <= 1 || len(msc.lowPivots) <= 1 {
		// need to set msc.trendClassification here?
		return Mixed
	}

	switch msc.thresholdQualifier {
	case SlopeSingle, SlopeDouble:
		msc.highPivotDiff = msc.getSlope(msc.bars, func(c candle.Candle) float64 { return c.High }, msc.highPivots)
		msc.lowPivotDiff = msc.getSlope(msc.bars, func(c candle.Candle) float64 { return c.Low }, msc.lowPivots)
		break
	case MagnitudeSingle, MagnitudeDouble:
		lastHighPivot := msc.getLastPivot(msc.highPivots)
		firstHighPivot := msc.getFirstPivot(msc.highPivots)
		if lastHighPivot != nil && firstHighPivot != nil {
			msc.highPivotDiff = lastHighPivot.High - firstHighPivot.High
		}

		lastLowPivot := msc.getLastPivot(msc.lowPivots)
		firstLowPivot := msc.getFirstPivot(msc.lowPivots)
		if lastLowPivot != nil && firstLowPivot != nil {
			msc.lowPivotDiff = lastLowPivot.Low - firstLowPivot.Low
		}
	}

	isStrongUp := false
	isStrongDown := false

	switch msc.thresholdQualifier {
	case SlopeSingle, MagnitudeSingle:
		isStrongUp = msc.highPivotDiff >= msc.strongTrendThreshold
		isStrongDown = msc.lowPivotDiff <= -msc.strongTrendThreshold
	case SlopeDouble, MagnitudeDouble:
		isStrongUp = msc.highPivotDiff >= msc.strongTrendThreshold && msc.lowPivotDiff >= msc.strongTrendThreshold
		isStrongDown = msc.highPivotDiff <= -msc.strongTrendThreshold && msc.lowPivotDiff <= -msc.strongTrendThreshold
	}

	if isStrongUp {
		msc.trendClassification = StrongUp
	} else if msc.highPivotDiff > 0 && msc.lowPivotDiff > 0 {
		msc.trendClassification = WeakUp
	} else if isStrongDown {
		msc.trendClassification = StrongDown
	} else if msc.highPivotDiff < 0 && msc.lowPivotDiff < 0 {
		msc.trendClassification = WeakDown
	} else {
		msc.trendClassification = Mixed
	}

	return msc.trendClassification
}

func (msc *MarketStructureClassifier) getLastPivot(pivots []int) *candle.Candle {
	if pivots == nil || len(pivots) == 0 {
		return nil
	}

	lastPivot := pivots[len(pivots)-1]
	if msc.bars == nil || lastPivot >= len(msc.bars) {
		return nil
	}

	return &msc.bars[lastPivot]
}

func (msc *MarketStructureClassifier) getFirstPivot(pivots []int) *candle.Candle {
	if pivots == nil || len(pivots) == 0 {
		return nil
	}

	firstPivot := pivots[0]
	if msc.bars == nil || firstPivot >= len(msc.bars) {
		return nil
	}

	return &msc.bars[firstPivot]
}

func (msc *MarketStructureClassifier) getSlope(candles []candle.Candle, valueAccessor func(candle.Candle) float64, pivots []int) float64 {
	highestValue := valueAccessor(candles[0])
	lowestValue := valueAccessor(candles[0])

	for _, c := range candles {
		if valueAccessor(c) > highestValue {
			highestValue = valueAccessor(c)
		}
		if valueAccessor(c) < lowestValue {
			lowestValue = valueAccessor(c)
		}
	}

	valueScaleDenom := highestValue - lowestValue
	firstPivotNum := valueAccessor(candles[pivots[0]]) - lowestValue
	lastPivotNum := valueAccessor(candles[pivots[len(pivots)-1]]) - lowestValue
	pivotDiff := lastPivotNum - firstPivotNum
	pivotDiffScaled := pivotDiff / valueScaleDenom
	timeScaled := float64(pivots[len(pivots)-1]-pivots[0]) / float64(len(candles)-1)

	return pivotDiffScaled / timeScaled
}

func (msc *MarketStructureClassifier) GetLastHighPivot() *candle.Candle {
	return msc.getLastPivot(msc.highPivots)
}

func (msc *MarketStructureClassifier) GetFirstHighPivot() *candle.Candle {
	return msc.getFirstPivot(msc.highPivots)
}

func (msc *MarketStructureClassifier) GetLastLowPivot() *candle.Candle {
	return msc.getLastPivot(msc.lowPivots)
}

func (msc *MarketStructureClassifier) GetFirstLowPivot() *candle.Candle {
	return msc.getFirstPivot(msc.lowPivots)
}

func (msc *MarketStructureClassifier) GetLastHighPivotIndex() int {
	if msc.highPivots == nil || len(msc.highPivots) == 0 {
		return -1
	}
	return msc.highPivots[len(msc.highPivots)-1]
}

func (msc *MarketStructureClassifier) GetFirstHighPivotIndex() int {
	if msc.lowPivots == nil || len(msc.lowPivots) == 0 {
		return -1
	}
	return msc.highPivots[0]
}

func (msc *MarketStructureClassifier) GetLastLowPivotIndex() int {
	if msc.lowPivots == nil || len(msc.lowPivots) == 0 {
		return -1
	}
	return msc.lowPivots[len(msc.lowPivots)-1]
}

func (msc *MarketStructureClassifier) GetFirstLowPivotIndex() int {
	if msc.lowPivots == nil || len(msc.lowPivots) == 0 {
		return -1
	}
	return msc.lowPivots[0]
}
