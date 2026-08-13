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
