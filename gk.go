// Package gk implmements Greenwald/Khanna's streaming quantiles
/*

"Space-Efficient Online Computation of Quantile Summaries" (Greenwald, Khanna 2001)

http://infolab.stanford.edu/~datar/courses/cs361a/papers/quantiles.pdf

This implementation is backed by a skiplist to make inserting elements into the
summary faster.  Querying is still O(n).

*/
package gk

// Stream is a quantile summary
type Stream struct {
	summary *skiplist
	epsilon float64
	n       int
}

type tuple struct {
	v     float64
	g     int
	delta int
}

// New returns a new stream with accuracy epsilon (0 <= epsilon <= 1)
func New(epsilon float64) *Stream {
	return &Stream{
		epsilon: epsilon,
		summary: newSkiplist(),
	}
}

// Insert inserts an item into the quantile summary
func (s *Stream) Insert(v float64) {

	value := tuple{v, 1, 0}

	elt := s.summary.Insert(value)

	s.n++

	if elt.prev[0] != s.summary.head && elt.next[0] != nil {
		elt.value.delta = int(2 * s.epsilon * float64(s.n))
	}

	if s.n%int(1.0/float64(2.0*s.epsilon)) == 0 {
		s.compress()
	}
}

func (s *Stream) compress() {

	var missing int

	epsN := int(2 * s.epsilon * float64(s.n))

	for elt := s.summary.head.next[0]; elt != nil && elt.next[0] != nil; {
		next := elt.next[0]
		t := elt.value
		nt := &next.value

		// value merging
		if t.v == nt.v {
			missing += nt.g
			nt.delta += missing
			nt.g = t.g
			s.summary.Remove(elt)
		} else if t.g+nt.g+missing+nt.delta < epsN {
			nt.g += t.g + missing
			missing = 0
			s.summary.Remove(elt)
		} else {
			nt.g += missing
			missing = 0
		}
		elt = next
	}
}

// Query returns an epsilon estimate of the element at quantile 'q' (0 <= q <= 1)
func (s *Stream) Query(q float64) float64 {

	// convert quantile to rank

	r := int(q*float64(s.n) + 0.5)

	var rmin int

	epsN := int(s.epsilon * float64(s.n))

	for elt := s.summary.head.next[0]; elt != nil; elt = elt.next[0] {

		t := elt.value

		rmin += t.g

		n := elt.next[0]

		if n == nil {
			return t.v
		}

		if r+epsN < rmin+n.value.g+n.value.delta {

			if r+epsN < rmin+n.value.g {
				return t.v
			}

			return n.value.v
		}
	}

	panic("not reached")
}
