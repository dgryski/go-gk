package gk

import (
	"math"
	"sort"
	"testing"
)

func TestSortFBits(t *testing.T) {

	data := []float64{-2, -2, -8, 6, 8, 10, -10, -6, -4, 0, 4}

	var fb f64bits

	for _, f := range data {
		fb = append(fb, math.Float64bits(f))
	}

	sort.Sort(fb)

	var floats []float64
	for _, f := range fb {
		floats = append(floats, math.Float64frombits(f))
	}

	if !sort.Float64sAreSorted(floats) {
		t.Error("f64bits sort order failure: ", floats)
	}

}
