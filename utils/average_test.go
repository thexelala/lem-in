package utils

import (
	"testing"
)

func TestFeedAverage(t *testing.T) {
	var avg Average

	checkFloatAndRational := func(numerator int) {
		if avg.Get() != float64(numerator) {
			t.Error("failed to get float64 average")
		}

		if avg.GetR() != New(int64(numerator), 1) {
			t.Error("failed to get Rarional average")
		}
	}

	avg.Feed(float64(10))
	checkFloatAndRational(10)

	avg.Feed(int64(8))
	checkFloatAndRational(9)

	avg.Feed(float64(6))
	checkFloatAndRational(8)

	avg.Feed(4)
	checkFloatAndRational(7)

	var avgR Average

	avgR.Feed(New(1, 2000))
	if avgR.GetR() != New(1, 2000) {
		t.Error("failed to get Rarional average", avgR.GetR())
	}

	avgR.Feed(New(1, 6000))
	if avgR.GetR() != New(1, 3000) {
		t.Error("failed to get Rarional average", avgR.GetR())
	}
}
