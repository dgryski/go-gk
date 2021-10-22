package gk

import (
	"math"
	"math/rand"
	"sort"
	"testing"
)

var sampleData = makeSampleData()
var sampleScratch = make([]uint64, len(sampleData))

func makeSampleData() []uint64 {
	sd := make([]uint64, 1000000)
	for i := range sd {
		sd[i] = math.Float64bits(rand.Float64()*float64(rand.Uint64()) - rand.Float64()*float64(rand.Uint64()))
	}
	return sd
}

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

func BenchmarkSortFBits(b *testing.B) {
	for l := 0; l < b.N; l++ {
		b.StopTimer()
		copy(sampleScratch, sampleData)
		b.StartTimer()
		sort.Sort(f64bits(sampleScratch))
	}
}
