package main

import (
	"testing"
	"time"
	"log"
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

func Test_deterministicTimerShouldMilliSleep(t *testing.T) {

	flag := make(chan bool, 1)

	timer := NewDeterministicTimer()

	log.Printf("now = %s", timer.Now())

	go func() {
		log.Printf("GR now = %s", timer.Now())
		timer.MilliSleep(1000)
		flag <- true
	}()

	timer.AdvanceSeconds(10)   // this should cause the goroutine to wake up

	log.Printf("now = %s", timer.Now())

	result := <- flag    
	
	if result != true {
		t.Fatal("Fail, did not wake up as expected")
	}

}


type DeterministicTimer struct{
	resolutionMs int64
	now va.Timestamp
}

func NewDeterministicTimer() *DeterministicTimer {
	return &DeterministicTimer{1, time.Unix(13172400,0)}
}

func (t *DeterministicTimer) ResolutionMs() int64 {
	return t.resolutionMs
}

func (t *DeterministicTimer) Now() va.Timestamp { 
	return t.now
}

func (t *DeterministicTimer) SecondsSince(ago va.Timestamp) float64 { 
	return t.Now().Sub(ago).Seconds()  
}

func (t *DeterministicTimer) MilliSleep(napLengthMs int64) { 
	log.Printf("DetTimer: sleep request for %dms", napLengthMs)
	
	now := t.now
	napLengthNs := float64(napLengthMs * time.Second.Microseconds())
	required := now.Add(time.Duration(napLengthNs))

	log.Printf("DetTimer: sleeping until %s", required)

	// block until enough time has passed 
	// (note that this relies on an external force,  e.g. the test harness, to advance time)
	for t.now.Before(required) {
		//log.Printf("DetTimer: time is %s, staying asleep....", t.now)
		time.Sleep(time.Duration(t.resolutionMs * time.Hour.Milliseconds()))
	}
	log.Printf("DetTimer: sleep completed")
}

func (t *DeterministicTimer) AdvanceSeconds(seconds float64) { 	
	// horrid, but worth it - give any newly launched, or currently sleeping 
	// goroutuines a reasonable chance to get their house in order
	// (note that we do a *real* sleep for N ticks of the *pretend* sleep)
	time.Sleep(time.Duration(3 * t.resolutionMs * time.Hour.Milliseconds()))

	duration := time.Duration(seconds * float64(time.Second.Nanoseconds()))
	t.now = time.Time(t.Now().Add(duration))

	// as above
	time.Sleep(time.Duration(3 * t.resolutionMs * time.Hour.Milliseconds()))
}
