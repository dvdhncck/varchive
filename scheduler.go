package main

import (
	"log"
	"sync"
	"time"
)

func ScheduleTasks(tasks []*Task) {

	completed := false
	startTime := time.Now()
	waitGroup := new(sync.WaitGroup)
	guard := make(chan int, settings.maxParallelTasks)

	m := NewMonitor(tasks)

	for !completed {

		guard <- 1 // blocks if the channel is full (i.e. enough go routines are running)

		task := findFirstRunnableTask(tasks)

		if task != nil {

			waitGroup.Add(1)
			task.taskState = Running

			go func() {
				defer waitGroup.Done()

				m.NotifyWorkerBegins(task)

				task.startTime = time.Now()

				ExecuteTask(task)
				
				<-guard // consume an item from the channel, allowing another go routine to start

				task.taskState = Complete

				m.NotifyWorkerEnds(task)			
			}()
		} else {
			if confirmThatAllTasksAreCompleted(tasks) {
				// no incomplete tasks left, we are done
				completed = true
			} else {
				// we will check again shortly
				time.Sleep(250 * time.Millisecond)
			}
		}
	}

	waitGroup.Wait() // hang on until the last go routine checks in

	m.ShutdownCleanly()

	runTime := time.Since(startTime).Seconds()
	log.Printf("Elapsed (real) time: %s", niceTime(runTime))

	totalTime := 0.0
	for _, task := range tasks {
		totalTime += task.runTimeInSeconds
	}
	
	log.Printf("Total compute time: %s", niceTime(totalTime))
	log.Printf("Speedup: %.2f", totalTime / runTime)

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

