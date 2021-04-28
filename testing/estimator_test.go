package main

import (
	"fmt"
	"math"
	"testing"

	"davidhancock.com/varchive"
)

const TASK = 1
const MIBIBYTES = 1000 * 1000

func Test_estimationAfterOneTask(t *testing.T) {

	task1 := varchive.NewTask(TASK, "", "", 3000)
	tasks := []*varchive.Task{task1}

	timer := NewDeterministicTimer()
	e := varchive.NewEstimator()
	m := varchive.NewMonitor(timer, e, tasks, varchive.NewNoOpDisplay())
	m.Start()

	m.NotifyTaskBegins(task1)
	timer.AdvanceSeconds(6)
	m.NotifyTaskEnds(task1)

	expected := 3000.0 / 6.0
	actual := e.EstimateBytesPerSecond(TASK)

	assertEqual(t, "estimate after 1 task", expected, actual)
}

func Test_estimationAfterTwoSerialTasks(t *testing.T) {

	expected, actual := float64(0), float64(0)

	task1 := varchive.NewTask(TASK, "", "", 3000)
	task2 := varchive.NewTask(TASK, "", "", 2000)

	tasks := []*varchive.Task{task1, task2}

	timer := NewDeterministicTimer()

	e := varchive.NewEstimator()
	m := varchive.NewMonitor(timer, e, tasks, varchive.NewNoOpDisplay())
	m.Start()

	m.NotifyTaskBegins(task1)
	timer.AdvanceSeconds(6)	
	m.NotifyTaskEnds(task1)
	m.NotifyTaskBegins(task2)

	// we have a estimate of 500/s based on task1, task 2 is 2000, so should take 4s
	expected = 2000.0 / (3000.0 / 6.0) 
	actual = m.EstimateTimeRemaining(task2)
	assertEqual(t, "1st estimate using valid data", expected, actual)

	timer.AdvanceSeconds(1)

	// we have a estimate of 500/s based on task1, task 2 is 2000, estimate after 1s...
	expected = (2000.0 / (3000.0 / 6.0)) - 1.0
	actual = m.EstimateTimeRemaining(task2)
	assertEqual(t, "2nd estimate using valid data", expected, actual)

	timer.AdvanceSeconds(2)

	// now, after 3s, we should forecast 1s remaining
	expected = (2000.0 / (3000.0 / 6.0)) - 3.0 
	actual = m.EstimateTimeRemaining(task2)
	assertEqual(t, "3rd estimate using valid data", expected, actual)

	// our estimate was wrong, the task actually completed after 3s
	m.NotifyTaskEnds(task2)

	// after task 2 completes, we update the estimate to include the new data 
	expected = (3000.0 + 2000.0) / (6.0 + 3.0)
	actual = e.EstimateBytesPerSecond(TASK)

	assertEqual(t, "estimate after 2 tasks", expected, actual)
}

func Test_estimationAfterTwoParallelTasks(t *testing.T) {

	task1 := varchive.NewTask(TASK, "", "", 3000)
	task2 := varchive.NewTask(TASK, "", "", 2000)

	tasks := []*varchive.Task{task1, task2}

	timer := NewDeterministicTimer()

	e := varchive.NewEstimator()
	m := varchive.NewMonitor(timer, e, tasks, varchive.NewNoOpDisplay())
	m.Start()

	m.NotifyTaskBegins(task1)
	m.NotifyTaskBegins(task2)

	timer.AdvanceSeconds(6)

	m.NotifyTaskEnds(task1)

	// there are 2 tasks running, so we estimate the performance 
	// of one task on its own would be double that of that seen with 2 in parallel

	expected := (3000.0 / 6.0) * 2  
	actual := e.EstimateBytesPerSecond(TASK)

	assertEqual(t, "estimate after 2 tasks", expected, actual)
}

func Test_estimationWhenTaskIsOverrunning(t *testing.T) {
	task1 := varchive.NewTask(TASK, "", "", 3000)
	task2 := varchive.NewTask(TASK, "", "", 3000)
	tasks := []*varchive.Task{task1, task2}
	
	timer := NewDeterministicTimer()
	e := varchive.NewEstimator()
	m := varchive.NewMonitor(timer, e, tasks, varchive.NewNoOpDisplay())
	m.Start() 

	m.NotifyTaskBegins(task1)
	timer.AdvanceSeconds(5)
	m.NotifyTaskEnds(task1)

	m.NotifyTaskBegins(task2)
	timer.AdvanceSeconds(20)
		
	// we have an estimate, but task2 has overrun the predicated 5 s
	// seconds and still not finished, so we get a signal saying "I dunno"

	expected := math.Inf(-1)
	actual := m.EstimateTimeRemaining(task2)
	assertEqual(t, "estimate when task is overrunning", expected, actual)
}

func Test_estimationOfRemainingRunTime(t *testing.T) {
	task1 := varchive.NewTask(TASK, "", "", 1 * MIBIBYTES)
	task2 := varchive.NewTask(TASK, "", "", 2 * MIBIBYTES)
	task3 := varchive.NewTask(TASK, "", "", 3 * MIBIBYTES)
	tasks := []*varchive.Task{task1, task2, task3}
	
	timer := NewDeterministicTimer()
	e := varchive.NewEstimator()
	m := varchive.NewMonitor(timer, e, tasks, varchive.NewNoOpDisplay())
	m.Start() 


	runTimeSeconds := e.EstimateRemainingRunTime(tasks)
	// cost_per_MB is initially 10, so estimate is: (1+2+3) * 10
	assertEqual(t, "estimate total time for all tasks", true, runTimeSeconds == 60)

	m.NotifyTaskBegins(task1)
	timer.AdvanceSeconds(20)
	task1.MarkAsCompleted()
	m.NotifyTaskEnds(task1)    
	runTimeSeconds = e.EstimateRemainingRunTime(tasks)
	// first task was slow, cost_per_MB is now 20/1, estimate is: (2+3) * (20/1)
	assertEqual(t, "estimate total time for all tasks", true, runTimeSeconds == 100)


	m.NotifyTaskBegins(task2)
	timer.AdvanceSeconds(5)
	task2.MarkAsCompleted()
	m.NotifyTaskEnds(task2)    
	runTimeSeconds = e.EstimateRemainingRunTime(tasks)

	// cost_per_MB is now  25_seconds/3_mb, estimate is: 3_mb * cost_per_MB, i.e. 3 * (25/3)
	assertEqual(t, "estimate total time for all tasks", true, runTimeSeconds == 25)

}

func assertEqual(t *testing.T, message string, expected interface{}, actual interface{}) {
	if expected == actual {
		return
	}
	t.Fatal(fmt.Sprintf("%s : expected %v, got %v", message, expected, actual))
}
