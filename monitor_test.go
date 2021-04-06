package main

import (
	"fmt"
	//"davidhancock.com/varchive/main"
	//"main"
	"testing"
)

func TestEstimation(t *testing.T) {

	task1 := NewTask(FixAudio, "", "", 3000)

	task2 := NewTask(FixAudio, "", "", 2000)

	task3 := NewTask(FixAudio, "", "", 1000)

	tasks := []*Task{task1, task2, task3}

	m := NewMonitor(TimerImpl{}, tasks)

	m.NotifyTaskBegins(task1)

	//task1.runTimeInSeconds = 6

	m.NotifyTaskEnds(task1)

	expected := 3000.0 / 6.0
	actual := m.EstimateBytesPerSecond(FixAudio)

	assertEqual(t, "simple estimate after 1 task", expected, actual)

}

func assertEqual(t *testing.T, message string, expected interface{}, actual interface{}) {
	if expected == actual {
		return
	}
	t.Fatal(fmt.Sprintf("%s : expected %v, got %v", message, expected, actual))
}
