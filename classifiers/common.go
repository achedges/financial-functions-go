package classifiers

type TrendClassification int

const (
	StrongUp TrendClassification = iota
	WeakUp
	Mixed
	WeakDown
	StrongDown
)

type PriceAction int

const (
	Hammer PriceAction = iota
	BullishEngulfing
	Piercing
	TweezerBottom
	MorningStar
	ShootingStar
	BearishEngulfing
	DarkCloudCover
	TweezerTop
	EveningStar
)

type ThresholdQualifier string

const (
	SlopeSingle     ThresholdQualifier = "SlopeSingle"
	SlopeDouble     ThresholdQualifier = "SlopeDouble"
	MagnitudeSingle ThresholdQualifier = "MagnitudeSingle"
	MagnitudeDouble ThresholdQualifier = "MagnitudeDouble"
)
