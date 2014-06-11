package gk

import "math"

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

func New(epsilon float64) *Stream {
	return &Stream{
		epsilon: epsilon,
		summary: newSkiplist(),
	}
}

func (s *Stream) Insert(v float64) {

	value := tuple{v, 1, 0}

	elt := s.summary.Insert(value)

	if elt.prev[0] != s.summary.head && elt.next[0] != nil {
		//elt.value.delta = int(math.Floor(2 * s.epsilon * float64(s.n)))
		elt.value.delta = elt.next[0].value.g + elt.next[0].value.delta - 1
	}

	s.n++
	if s.n%int(1.0/float64(2.0*s.epsilon)) == 0 {
		s.compress()
	}
}

func (s *Stream) compress() {

	for elt := s.summary.head.next[0]; elt != nil && elt.next[0] != nil; {
		next := elt.next[0]
		t := elt.value
		nt := &next.value
		if t.g+nt.g+nt.delta <= int(math.Floor(2*s.epsilon*float64(s.n))) {
			nt.g += t.g
			s.summary.Remove(elt)
		}
		elt = next
	}
}

func (s *Stream) Query(q float64) float64 {

	// convert quantile to rank

	r := int(q * float64(s.n))

	var rmin int

	for elt := s.summary.head.next[0]; elt.next[0] != nil; elt = elt.next[0] {

		t := elt.value

		rmin += t.g
		rmax := rmin + t.delta

		if r-rmin <= int(s.epsilon*float64(s.n)) && rmax-r <= int(s.epsilon*float64(s.n)) {
			return t.v
		}
	}

	// panic("not reached")

	return 0
}
