package main

import (
	"testing"
	"time"
	va "davidhancock.com/varchive"
)

func Test_deterministicTimerShouldAdvance(t *testing.T) {

	timer := NewDeterministicTimer()
	
	start := timer.Now();

	timer.AdvanceSeconds(7.0)

	delta := timer.SecondsSince(start)

	if delta != 7.0 {
		t.Fatal("Fail, did not advance as expected")
	}
}

func Test_deterministicTimerShouldSleep(t *testing.T) {

	timer := NewDeterministicTimer()

	before := timer.Now()

	timer.MilliSleep(500);

	after := timer.Now()

	delta := after.Sub(before).Milliseconds()

	if delta != 500 {
		t.Fatal("Fail, did not sleep as expected")
	}
}


type DeterministicTimer struct{
	now va.Timestamp
}

func NewDeterministicTimer() *DeterministicTimer {
	return &DeterministicTimer{time.Now()}
}

func (t *DeterministicTimer) Now() va.Timestamp { 
	return t.now
}

func (t *DeterministicTimer) SecondsSince(ago va.Timestamp) float64 { 
	return t.Now().Sub(ago).Seconds()  
}

func (t *DeterministicTimer) MilliSleep(nap int64) { 
	t.AdvanceSeconds(float64(nap)* 0.001)
}

func (t *DeterministicTimer) AdvanceSeconds(seconds float64) { 
	duration := time.Duration(seconds * float64(time.Second.Nanoseconds()))
	t.now = time.Time(t.Now().Add(duration))
}
