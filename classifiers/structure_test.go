package classifiers_test

import (
	"financial/functions/classifiers"
	"math/rand"
	"testing"

	"github.com/achedges/financial-core-go/pricebar"
	"github.com/achedges/go-assertions"
)

func getRandomBars() []pricebar.PriceBar {
	random := rand.New(rand.NewSource(1234))
	bars := make([]pricebar.PriceBar, 20)
	for i := range 20 {
		bars[i] = *pricebar.New(pricebar.Config{
			Symbol:     "TEST",
			BasisPrice: 20.0 + random.Float64(),
		})
	}
	return bars
}

func TestMarketStructureClassifier_init(t *testing.T) {
	classifier1 := classifiers.NewMarketStructure(classifiers.DefaultMarketStructureConfig)
	assertions.EqualFloats(0.25, classifier1.GetStrongTrendThreshold(), t)

	classifier2 := classifiers.NewMarketStructure(classifiers.MarketStructureConfig{
		StrongTrendThreshold: 0.5,
	})
	assertions.EqualFloats(0.5, classifier2.GetStrongTrendThreshold(), t)
}

func TestMarketStructureClassifier_GetLastPivots(t *testing.T) {
	cls := classifiers.NewMarketStructure(classifiers.DefaultMarketStructureConfig)

	assertions.True(cls.GetLastHighPivot() == nil, t)
	assertions.True(cls.GetLastLowPivot() == nil, t)

	bars := getRandomBars()
	// generate an ambiguous structure
	bars[4].High = 24.0
	bars[14].High = 22.0
	bars[8].Low = 16.0
	bars[18].Low = 18.0
	cls.Classify(bars)

	lastHighPivot := cls.GetLastHighPivot()
	lastLowPivot := cls.GetLastLowPivot()

	if lastHighPivot == nil || lastLowPivot == nil {
		t.FailNow()
		return
	}

	assertions.EqualFloats(22.0, cls.GetLastHighPivot().High, t)
	assertions.EqualFloats(18.0, cls.GetLastLowPivot().Low, t)
	assertions.EqualInts(14, cls.GetLastHighPivotIndex(), t)
	assertions.EqualInts(18, cls.GetLastLowPivotIndex(), t)
}

func TestMarketStructureClassifier_Mixed(t *testing.T) {
	cls := classifiers.NewMarketStructure(classifiers.MarketStructureConfig{
		StrongTrendThreshold: 1.0,
	})

	bars := getRandomBars()

	// base case, not enough pivots identified
	trendClass := cls.Classify(bars)
	assertions.True(trendClass == classifiers.Mixed, t)

	// generate a perfectly flat structure
	bars[4].High = 22.0
	bars[14].High = 22.0
	bars[8].Low = 18.0
	bars[18].Low = 18.0
	trendClass = cls.Classify(bars)
	assertions.True(trendClass == classifiers.Mixed, t)
	assertions.EqualFloats(0.0, cls.GetHighPivotDiff(), t)
	assertions.EqualFloats(0.0, cls.GetLowPivotDiff(), t)

	// generate an ambiguous structure
	bars[4].High = 24.0
	bars[14].High = 22.0
	bars[8].Low = 16.0
	bars[18].Low = 18.0
	trendClass = cls.Classify(bars)
	assertions.True(trendClass == classifiers.Mixed, t)
	assertions.EqualFloats(-2.0, cls.GetHighPivotDiff(), t)
	assertions.EqualFloats(2.0, cls.GetLowPivotDiff(), t)
}

func TestMarketStructureClassifier_WeakUp(t *testing.T) {
	cls := classifiers.NewMarketStructure(classifiers.MarketStructureConfig{
		StrongTrendThreshold: 1.0,
	})

	bars := getRandomBars()
	bars[4].High = 22.0
	bars[14].High = 22.95 // just under 1.0
	bars[8].Low = 18.0
	bars[18].Low = 18.95 // just under 1.0
	trendClass := cls.Classify(bars)
	assertions.True(trendClass == classifiers.WeakUp, t)
	assertions.CloseEnough(0.95, cls.GetHighPivotDiff(), 0.001, t)
	assertions.CloseEnough(0.95, cls.GetLowPivotDiff(), 0.001, t)
}

func TestMarketStructureClassifier_WeakDown(t *testing.T) {
	cls := classifiers.NewMarketStructure(classifiers.MarketStructureConfig{
		StrongTrendThreshold: 1.0,
	})

	bars := getRandomBars()
	bars[4].High = 22.0
	bars[14].High = 21.05 // just under 1.0
	bars[8].Low = 18.0
	bars[18].Low = 17.05 // just under 1.0
	trendClass := cls.Classify(bars)
	assertions.True(trendClass == classifiers.WeakDown, t)
	assertions.CloseEnough(-0.95, cls.GetHighPivotDiff(), 0.001, t)
	assertions.CloseEnough(-0.95, cls.GetLowPivotDiff(), 0.001, t)
}

func TestMarketStructureClassifier_StrongUp(t *testing.T) {
	cls := classifiers.NewMarketStructure(classifiers.MarketStructureConfig{
		StrongTrendThreshold: 0.1,
	})

	bars := getRandomBars()
	bars[4].High = 22.0
	bars[14].High = 23.05 // just over 1.0
	bars[8].Low = 18.0
	bars[18].Low = 19.05 // just over 1.0
	trendClass := cls.Classify(bars)
	assertions.True(trendClass == classifiers.StrongUp, t)
	assertions.CloseEnough(1.05, cls.GetHighPivotDiff(), 0.001, t)
	assertions.CloseEnough(1.05, cls.GetLowPivotDiff(), 0.001, t)
}

func TestMarketStructureClassifier_StrongDown(t *testing.T) {
	cls := classifiers.NewMarketStructure(classifiers.MarketStructureConfig{
		StrongTrendThreshold: 0.1,
	})

	bars := getRandomBars()
	bars[4].High = 22.0
	bars[14].High = 20.95 // just over 1.0
	bars[8].Low = 18.0
	bars[18].Low = 16.95 // just over 1.0
	trendClass := cls.Classify(bars)
	assertions.True(trendClass == classifiers.StrongDown, t)
	assertions.CloseEnough(-1.05, cls.GetHighPivotDiff(), 0.001, t)
	assertions.CloseEnough(-1.05, cls.GetLowPivotDiff(), 0.001, t)
}
