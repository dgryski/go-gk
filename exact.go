package gk

import (
	"math"
	"sort"
)

// Exact is an exact (map-based) quantile summary
type Exact struct {
	summary map[uint64]int
	n       int
	// we use float64bits because uint64 is an optimized map type; it's 100x faster
	keys f64bits
}

// NewExact returns a new exact quantile summary
func NewExact() *Exact {
	return &Exact{
		summary: make(map[uint64]int),
	}
}

// Insert inserts an item into the quantile summary
func (ex *Exact) Insert(v float64) {
	ex.n++

	vbits := math.Float64bits(v)
	ex.summary[vbits]++

	// clear out the cache of sorted keys
	ex.keys = nil
}

// Query returns the qth quantile of the stream
func (ex *Exact) Query(q float64) float64 {

	if ex.keys == nil {

		ex.keys = make([]uint64, 0, len(ex.summary))

		for k := range ex.summary {
			ex.keys = append(ex.keys, k)
		}

		sort.Sort(ex.keys)

	}

	// TODO(dgryski): create prefix sum array and then binsearch to find quantile.

	total := 0

	for _, k := range ex.keys {
		total += ex.summary[k]
		p := float64(total) / float64(ex.n)
		if q <= p {
			return math.Float64frombits(k)
		}
	}

	panic("not reached")
}

// to sort IEEE Float64's as bits
type f64bits []uint64

func (s f64bits) Len() int { return len(s) }
func (s f64bits) Less(i, j int) bool {
	sgni := s[i] >> 63
	sgnj := s[j] >> 63

	// both positive, then 'less than' is correct
	if sgni == sgnj && sgni == 0 {
		return s[i] < s[j]
	}

	// Either the signs differ, or they're both negative.  If we have one
	// positive and one negative, then the negative number which sorts
	// higher than a positive number (due to the sign bit) we actually want
	// to tag as 'less than'.  If they're both negative, then the 'larger'
	// uint64 value is the 'smaller' float64

	return s[i] >= s[j]

}
func (s f64bits) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
