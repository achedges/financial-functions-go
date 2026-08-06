package window

type BufferContainer struct {
	Index  int
	Period int
	Length int
}

func NewBufferContainer(index int, period int, length int) BufferContainer {
	return BufferContainer{
		Index:  index,
		Period: period,
		Length: length,
	}
}

func (b *BufferContainer) GetLowerBound() int {
	return (b.Index - b.Period + 1) % b.Length
}

func (b *BufferContainer) GetUpperBound(pad int) int {
	// pad is used to generate an exclusive upper bound, useful for feeding directly into a slice
	return (b.Index % b.Length) + pad
}

func (b *BufferContainer) IsRing() bool {
	return b.Period == b.Length
}

func (b *BufferContainer) Advance() {
	b.Index++
}
