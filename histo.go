package varchive

import (
	"fmt"
	"strings"
)

type Histo struct {
	min, max int64
	counts     map[int64]int64
}

func NewHisto() *Histo {
	return &Histo{0, 0, make(map[int64]int64)}
}

func (h *Histo) Min() int64 {
	return h.min
}

func (h *Histo) Max() int64 {
	return h.max
}

func (h *Histo) Add(key int64) {
	if len(h.counts) == 0 || key < h.min {
		h.min = key
	}
	if len(h.counts) == 0 || key > h.max {
		h.max = key
	}
	h.counts[key] = h.counts[key] + 1
}

func (h *Histo) Get(key int64) int64 {
	return h.counts[key]
}

func (h *Histo) String() string {
	if h.min == h.max {
		return fmt.Sprintf("%d instances of %d\n", h.counts[h.min], h.min)	
	}

	var b strings.Builder

	for key, count := range h.counts {
		fmt.Fprintf(&b, " %5.d @ %5.d\n", count, key)
	}

	return b.String()
}
