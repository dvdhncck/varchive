package varchive

import (
	"fmt"
	"strings"
)

type Histo struct {
	counts   map[string]int64
}

func NewHisto() *Histo {
	return &Histo{make(map[string]int64)}
}

func (h *Histo) Add(key string) {
	h.counts[key] = h.counts[key] + 1
}

func (h *Histo) Get(key string) int64 {
	return h.counts[key]
}

func (h *Histo) String() string {
	var b strings.Builder

	for key, count := range h.counts {
		fmt.Fprintf(&b, " %5.d @ %v\n", count, key)
	}

	return b.String()
}
