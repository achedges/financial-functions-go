package stats_test

import (
	"testing"

	"github.com/achedges/financial-core-go/candle"
	"github.com/achedges/financial-functions-go/stats"
	"github.com/achedges/go-assertions"
)

func TestRange_MinMax(t *testing.T) {
	period := 4
	rng := stats.NewRange(period, getPriceBarList(period))

	mx, hasMax := rng.GetRangeMax()
	mn, hasMin := rng.GetRangeMin()

	assertions.True(hasMax, t)
	assertions.True(hasMin, t)
	assertions.EqualFloats(305.12, mx, t)
	assertions.EqualFloats(302.09, mn, t)

	expectedMax := []float64{306.17, 309.01, 310.24, 311.26, 311.26, 311.26, 311.26, 310.35, 307.40, 307.40, 307.40, 307.40}
	expectedMin := []float64{302.09, 302.09, 302.09, 306.17, 309.01, 303.26, 300.42, 300.42, 300.42, 300.42, 303.50, 300.80}

	for i := 0; i < len(data)-period; i++ {
		c := candle.New(candle.Config{Symbol: "TEST", BasisPrice: data[i+period]})
		rng.Slide(c)
		mx, hasMax = rng.GetRangeMax()
		mn, hasMin = rng.GetRangeMin()
		assertions.True(hasMax, t)
		assertions.True(hasMin, t)
		assertions.EqualFloats(expectedMax[i], mx, t)
		assertions.EqualFloats(expectedMin[i], mn, t)
	}
}
