package varchive

import (
	"time"
)

// Timer provices a minimal wrapper around time.* so we can unit test time-critical 
// bits of the code
//
// The rest of the code base should not import the "time" at all,
// instead it should accept a Timer interface

type Timestamp = time.Time

type Timer interface {
	Now() Timestamp
	SecondsSince(Timestamp) float64
	MilliSleep(int64)
}

// the default implementation delegates to the real time.* methods

func NewTimer() Timer {
	return timerImpl{}
}

type timerImpl struct{}

func (timerImpl) Now() Timestamp                     { return time.Now() }
func (timerImpl) SecondsSince(ago Timestamp) float64 { return time.Since(ago).Seconds() }
func (timerImpl) MilliSleep(nap int64)               { time.Sleep(time.Duration(nap * int64(time.Millisecond))) }
