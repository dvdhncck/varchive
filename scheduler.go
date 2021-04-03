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

	for !completed {

		guard <- 1 // blocks if the channel is full (i.e. enough go routines are running)

		task := findFirstRunnableTask(tasks)

		if task != nil {
			log.Printf("Running task %v", task)

			waitGroup.Add(1)
			task.taskState = Running

			go func() {
				defer waitGroup.Done()

				taskStartTime := time.Now()

				ExecuteTask(task)
				
				<-guard // consume an item from the channel, allowing another go routine to start

				task.taskState = Complete

				task.runTime = time.Since(taskStartTime)

				log.Printf("Completed task %d in %v", task.id, task.runTime)
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

	runTime := time.Since(startTime)
	log.Printf("Elapsed (real) time: %v", runTime)

	totalTime := time.Duration(0)
	for _, task := range tasks {
		totalTime += task.runTime
	}
	log.Printf("Total compute time: %v", totalTime)

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

// efficient removal of item from list (does not preserve the order of the list)
// func remove(s []*Task, i int) []*Task {
// 	s[i] = s[len(s)-1]
// 	return s[:len(s)-1]
// }
