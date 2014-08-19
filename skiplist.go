package gk

import (
	"math/rand"
	"time"
)

const maxHeight = 31

type skiplist struct {
	height int
	head   *node
	rnd    *rand.Rand
}

type node struct {
	value tuple
	next  []*node
	prev  []*node
}

func newSkiplist() *skiplist {
	return &skiplist{
		height: 0,
		head:   &node{next: make([]*node, maxHeight)},
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *skiplist) Insert(t tuple) *node {
	level := 0

	n := s.rnd.Int31()
	for n&1 == 1 {
		level++
		n >>= 1
	}

	if level > s.height {
		s.height++
		level = s.height
	}

	node := &node{
		value: t,
		next:  make([]*node, level+1),
		prev:  make([]*node, level+1),
	}
	curr := s.head
	for i := s.height; i >= 0; i-- {

		for curr.next[i] != nil && t.v >= curr.next[i].value.v {
			curr = curr.next[i]
		}

		if i > level {
			continue
		}

		node.next[i] = curr.next[i]
		if curr.next[i] != nil && curr.next[i].prev[i] != nil {
			curr.next[i].prev[i] = node
		}
		curr.next[i] = node
		node.prev[i] = curr
	}

	return node
}

func (s *skiplist) Remove(node *node) {

	// remove n from each level of the skiplist

	for i := range node.next {
		prev := node.prev[i]
		next := node.next[i]

		if prev != nil {
			prev.next[i] = next
		}
		if next != nil {
			next.prev[i] = prev
		}
		node.next[i] = nil
		node.prev[i] = nil
	}
}
