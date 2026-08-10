package functions_test

import (
	"financial/functions"
	"testing"

	"github.com/achedges/go-assertions"
)

// We're always working with a buffer that's been backfilled, so the initial index should always be the end of the buffer array.

func TestBufferContainer_NewBufferContainer(t *testing.T) {
	buffer := functions.NewBufferContainer(1, 2, 3)
	assertions.EqualInts(1, buffer.Index, t)
	assertions.EqualInts(2, buffer.Period, t)
	assertions.EqualInts(3, buffer.Length, t)
}

func TestBufferContainer_GetLowerBound(t *testing.T) {
	linearBuffer := functions.NewBufferContainer(3, 4, 10)
	var linearExpected = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}
	for _, v := range linearExpected {
		assertions.EqualInts(v, linearBuffer.GetLowerBound(), t)
		linearBuffer.Advance()
	}

	ringBuffer := functions.NewBufferContainer(3, 4, 4)
	var ringExpected = []int{0, 1, 2, 3, 0, 1, 2, 3, 0, 1}
	for _, v := range ringExpected {
		assertions.EqualInts(v, ringBuffer.GetLowerBound(), t)
		ringBuffer.Advance()
	}
}

func TestBufferContainer_GetUpperBound(t *testing.T) {
	linearBuffer := functions.NewBufferContainer(3, 4, 10)
	var linearExpected = []int{3, 4, 5, 6, 7, 8, 9, 0, 1, 2}
	for _, v := range linearExpected {
		assertions.EqualInts(v, linearBuffer.GetUpperBound(0), t)
		linearBuffer.Advance()
	}

	linearBuffer.Index = 3
	var linearExpectedPad = []int{4, 5, 6, 7, 8, 9, 10, 1, 2, 3}
	for _, v := range linearExpectedPad {
		assertions.EqualInts(v, linearBuffer.GetUpperBound(1), t)
		linearBuffer.Advance()
	}

	ringBuffer := functions.NewBufferContainer(3, 4, 4)
	var ringExpected = []int{3, 0, 1, 2}
	for _, v := range ringExpected {
		assertions.EqualInts(v, ringBuffer.GetUpperBound(0), t)
		ringBuffer.Advance()
	}

	ringBuffer.Index = 3
	var ringExpectedPad = []int{4, 1, 2, 3}
	for _, v := range ringExpectedPad {
		assertions.EqualInts(v, ringBuffer.GetUpperBound(1), t)
		ringBuffer.Advance()
	}
}

func TestBufferContainer_IsRing(t *testing.T) {
	ringBuffer := functions.NewBufferContainer(9, 10, 10)
	assertions.True(ringBuffer.IsRing(), t)

	linearBuffer := functions.NewBufferContainer(9, 10, 20)
	assertions.False(linearBuffer.IsRing(), t)
}

func TestBufferContainer_Advance(t *testing.T) {
	buffer := functions.NewBufferContainer(9, 10, 10)
	var expected = []int{9, 10, 11, 12, 13, 14, 15, 16, 17, 18}
	for _, v := range expected {
		assertions.EqualInts(v, buffer.Index, t)
		buffer.Advance()
	}
}
