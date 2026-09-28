package stats

import (
	"testing"

	"github.com/achedges/financial-core-go/candle"
	"github.com/achedges/go-assertions"
)

func getCandle(open, close float64) candle.Candle {
	c := candle.New(candle.Config{Symbol: "", BasisPrice: open})
	c.Close = close // set this first so IsUp() works
	if c.IsUp() {
		c.High = close
		c.Low = open
	} else {
		c.High = open
		c.Low = close
	}
	return c
}

func TestParaboliSAR_Init(t *testing.T) {
	c := getCandle(10.0, 11.0)
	psar := NewParabolicSAR(c)
	assertions.True(psar.isLong, t)
	assertions.CloseEnough(10.0, psar.stop.InexactFloat64(), 0.001, t)
	assertions.CloseEnough(10.0, psar.low.InexactFloat64(), 0.001, t)
	assertions.CloseEnough(11.0, psar.high.InexactFloat64(), 0.001, t)
}

func TestParabolicSAR_Slide_Basic_Long(t *testing.T) {
	prices := [][]float64{
		{51.5, 52.5},
		{52, 53},
		{52.5, 53.5},
		{53, 54},
		{53.5, 54.5},
		{54, 55},
		{54.5, 55.5},
		{55, 56},
		{55.5, 56.5},
		{56, 57},
		{56.5, 57.5},
		{57, 58},
		{57.5, 58.5},
		{58, 59},
	}
	expectedSar := []float64{50.05, 50.16, 50.36, 50.66, 51.04, 51.52, 52.08, 52.71, 53.38, 54.11, 54.78, 55.43, 56.04, 56.63, 56.75}

	psar := NewParabolicSAR(getCandle(50, 51))
	assertions.CloseEnough(50.0, psar.stop.InexactFloat64(), 0.01, t)

	for i, p := range prices {
		psar.Slide(getCandle(p[0], p[1]))
		assertions.CloseEnough(expectedSar[i], psar.stop.InexactFloat64(), 0.01, t)
	}

	assertions.CloseEnough(0.2, psar.alpha.InexactFloat64(), 0.01, t)
}

func TestParabolicSAR_Slide_Basic_Short(t *testing.T) {
	prices := [][]float64{
		{59, 58},
		{58.5, 57.5},
		{58, 57},
		{57.5, 56.5},
		{57, 56},
		{56.5, 55.5},
		{56, 55},
		{55.5, 54.5},
		{55, 54},
		{54.5, 53.5},
		{54, 53},
		{53.5, 52.5},
		{53, 52},
		{52.5, 51.5},
	}
	expectedSar := []float64{60.45, 60.33, 60.13, 59.84, 59.45, 58.98, 58.42, 57.79, 57.11, 56.39, 55.71, 55.07, 54.45, 53.86, 53.75}

	psar := NewParabolicSAR(getCandle(60.5, 59.5))
	assertions.CloseEnough(60.5, psar.stop.InexactFloat64(), 0.01, t)

	for i, p := range prices {
		psar.Slide(getCandle(p[0], p[1]))
		assertions.CloseEnough(expectedSar[i], psar.stop.InexactFloat64(), 0.01, t)
	}

	assertions.CloseEnough(0.2, psar.alpha.InexactFloat64(), 0.01, t)
}

func TestParabolicSAR_Slide_Reversal(t *testing.T) {
	prices := [][]float64{
		{51.5, 52.5},
		{52, 53},
		{52.5, 53.5}, // max high should become new stop
		{53, 49.99},
		{50.5, 49.5},
		{50, 49},
	}
	expectedSar := []float64{50.05, 50.16, 50.36, 53.50, 53.42, 53.24}

	psar := NewParabolicSAR(getCandle(50, 51))
	assertions.CloseEnough(50.0, psar.stop.InexactFloat64(), 0.01, t)
	assertions.True(psar.isLong, t) // should start long

	for i, p := range prices {
		psar.Slide(getCandle(p[0], p[1]))
		assertions.CloseEnough(expectedSar[i], psar.stop.InexactFloat64(), 0.01, t)

		if i == 3 {
			assertions.True(psar.reverseSignal, t)
		} else {
			assertions.False(psar.reverseSignal, t)
		}
	}

	assertions.False(psar.isLong, t) // should end short
}

func TestParabolicSAR_MarketData(t *testing.T) {
	// EUR_USD, H1, 20260902 @ 1200 UTC -> 20260903 @ 2200 UTC
	initCandle := candle.New(candle.Config{Symbol: "EUR_USD", Date: 20260902, Time: 700})
	initCandle.Open = 1.15679
	initCandle.High = 1.15790
	initCandle.Low = 1.15665
	initCandle.Close = 1.15779

	priceData := [][]float64{
		{1.157740, 1.158570, 1.157580, 1.157840},
		{1.157840, 1.160940, 1.157790, 1.160440},
		{1.160440, 1.160440, 1.158870, 1.159340},
		{1.159360, 1.159700, 1.158900, 1.159210},
		{1.159200, 1.159400, 1.158740, 1.158740},
		{1.158740, 1.158990, 1.158340, 1.158660},
		{1.158670, 1.158980, 1.158440, 1.158580},
		{1.158580, 1.158770, 1.158260, 1.158740},
		{1.158740, 1.159000, 1.158720, 1.158850},
		{1.158740, 1.158840, 1.158700, 1.158720},
		{1.158760, 1.158980, 1.158580, 1.158700},
		{1.158700, 1.158740, 1.158440, 1.158520},
		{1.158500, 1.159100, 1.158350, 1.158950},
		{1.158960, 1.159360, 1.158510, 1.158620},
		{1.158610, 1.159500, 1.158520, 1.159120},
		{1.159120, 1.159800, 1.159100, 1.159670},
		{1.159690, 1.160040, 1.159460, 1.159560},
		{1.159550, 1.160060, 1.159360, 1.160020},
		{1.160020, 1.160830, 1.159630, 1.160460},
		{1.160460, 1.161480, 1.160180, 1.160400},
		{1.160410, 1.161150, 1.159980, 1.160910},
		{1.160920, 1.161250, 1.160140, 1.160580},
		{1.160590, 1.160710, 1.159790, 1.160400},
		{1.160400, 1.161460, 1.160190, 1.161260},
		{1.161280, 1.163020, 1.160840, 1.162580},
		{1.162590, 1.162880, 1.161740, 1.162740},
		{1.162740, 1.162760, 1.161330, 1.161830},
		{1.161860, 1.162760, 1.161740, 1.162300},
		{1.162290, 1.164140, 1.162130, 1.163680},
		{1.163680, 1.163860, 1.163280, 1.163280},
		{1.163280, 1.163660, 1.163180, 1.163480},
		{1.163470, 1.163520, 1.162530, 1.162870},
		{1.162880, 1.163010, 1.162500, 1.162560},
		{1.162620, 1.162900, 1.162520, 1.162560},
		{1.162600, 1.162780, 1.162290, 1.162310},
	}

	candles := make([]candle.Candle, 0, len(priceData))
	for _, ohlc := range priceData {
		c := candle.New(candle.Config{Symbol: "EUR_USD", BasisPrice: ohlc[0]})
		c.High = ohlc[1]
		c.Low = ohlc[2]
		c.Close = ohlc[3]
		candles = append(candles, c)
	}

	expectedStops := []float64{
		1.15665, 1.15674, 1.15686, 1.15702, 1.15718, 1.15733, 1.15747, 1.15761, 1.15775, 1.15787,
		1.15800, 1.15811, 1.15823, 1.15834, 1.15847, 1.15857, 1.15866, 1.15875, 1.15884, 1.15895,
		1.15910, 1.15924, 1.15937, 1.15950, 1.15971, 1.15998, 1.16022, 1.16044, 1.16074, 1.16108,
		1.16139, 1.16166, 1.16190, 1.16213, 1.16233,
	}

	psar := NewParabolicSAR(initCandle)
	psar.SetIsLong(false) // force to match actual setup

	for i, c := range candles {
		psar.Slide(c)
		assertions.CloseEnough(expectedStops[i], psar.stop.InexactFloat64(), 0.0001, t)
	}
}
