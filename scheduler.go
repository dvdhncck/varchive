package varchive

import (
	"log"
	"sync"
)

func ScheduleTasks(timer Timer, tasks []*Task) {

	if settings.verbose {
		log.Println("Scheduling tasks")
	}

	completed := false
	startTime := timer.Now()
	waitGroup := new(sync.WaitGroup)
	guard := make(chan int, settings.maxParallelTasks)

	m := NewMonitor(timer, tasks, getSuitableDisplay())

	if ! settings.dryRun {
		m.Start()   // dont bother starting the monitor, we don't need the output
	}

	for !completed {

		guard <- 1 // blocks if the channel is full (i.e. enough go routines are running)

		task := findFirstRunnableTask(tasks)

		if task != nil {

			waitGroup.Add(1)
			task.taskState = Running

			go func() {
				defer waitGroup.Done()

				m.NotifyTaskBegins(task)

				ExecuteTask(task)

				task.taskState = Complete

				m.NotifyTaskEnds(task)

				<-guard // consume an item from the channel, allowing another go routine to start
			}()
		} else {
			if confirmThatAllTasksAreCompleted(tasks) {
				// no incomplete tasks left, we are done
				completed = true
			} else {
				// we will check again shortly
				timer.MilliSleep(250)
			}
		}
	}

	waitGroup.Wait() // hang on until the last go routine checks in

	m.ShutdownCleanly()

	runTime := timer.SecondsSince(startTime)
	log.Printf("Elapsed (real) time: %s", niceTime(runTime))

	totalTime := 0.0
	for _, task := range tasks {
		totalTime += task.runTimeInSeconds
	}

	log.Printf("Total compute time: %s", niceTime(totalTime))
	log.Printf("Efficiency: %.2f", totalTime/runTime/float64(settings.maxParallelTasks))
}

func findFirstRunnableTask(tasks []*Task) *Task {
	for _, task := range tasks {
		if task.canRun() {
			return task
		}
	}
	return nil
}

func confirmThatAllTasksAreCompleted(tasks []*Task) bool {
	for _, task := range tasks {
		if task.isNotCompleted() {
			return false
		}
	}
	return true
}
