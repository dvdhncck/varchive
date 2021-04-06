package main

import (
	"time"
)

// Timer wraps the required time.* functionality so we can unit test time-critical bits of the code
// The rest of the code base should not import the "time" at all - but use a Timer when required

type Timestamp = time.Time

type Timer interface {
	Now() Timestamp
	SecondsSince(Timestamp) float64
	MilliSleep(int64)
}

type TimerImpl struct{}

func (c TimerImpl) Now() Timestamp                     { return time.Now() }
func (c TimerImpl) SecondsSince(ago Timestamp) float64 { return time.Since(ago).Seconds() }
func (c TimerImpl) MilliSleep(nap int64)               { time.Sleep(time.Duration(nap * time.Hour.Milliseconds())) }
