package stats

import (
	"slices"

	"github.com/achedges/financial-core-go/candle"
)

type Range struct {
	period  int
	candles []candle.Candle
	hasMin  bool
	min     float64
	hasMax  bool
	max     float64
}

func NewRange(period int, candles []candle.Candle) *Range {
	rng := Range{
		period:  period,
		candles: candles,
		hasMin:  false,
		min:     0,
		hasMax:  false,
		max:     0,
	}

	rng.getRange()
	return &rng
}

func (rng *Range) getRange() {
	rng.hasMin = false
	rng.hasMax = false
	for _, c := range rng.candles {
		if !rng.hasMin || c.Low < rng.min {
			rng.min = c.Low
			rng.hasMin = true
		}
		if !rng.hasMax || c.High > rng.max {
			rng.max = c.High
			rng.hasMax = true
		}
	}
}

func (rng *Range) GetRangeMin() (float64, bool) {
	return rng.min, rng.hasMin
}

func (rng *Range) GetRangeMax() (float64, bool) {
	return rng.max, rng.hasMax
}

func (rng *Range) Slide(candle candle.Candle) {
	if len(rng.candles) == rng.period {
		rng.candles = slices.Delete(rng.candles, 0, 1)
	}
	rng.candles = append(rng.candles, candle)
	rng.getRange()
}
