package main

import (
	"errors"
	"log"
	"sync"
	"time"
	//"math/rand"
)

func ScheduleTasks(tasks []*Task) {

	completed := false

	waitGroup := new(sync.WaitGroup)
	guard := make(chan int, settings.maxParallelTasks)

	for !completed {

		guard <- 1 // blocks if the channel is full (i.e. enough go routines are running)

		task, remainingTasks, err := findFirstRunnableTask(tasks)

		if err == nil {
			tasks = remainingTasks

			log.Printf("Running task %v", task)

			waitGroup.Add(1)
			task.taskState = Running

			go func() {
				defer waitGroup.Done()

				start := time.Now()

				ExecuteTask(task)
				
				<-guard // consume an item from the channel, allowing another go routine to start

				task.taskState = Complete

				elapsed := time.Since(start)

				log.Printf("Completed task %d in %v", task.id, time.Duration(elapsed))
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
}

// if there is a task ready to run
//   return (the task, the list with the task now removed, nil error)
// else
//   return (nil task, the unmodified list, an error)
func findFirstRunnableTask(tasks []*Task) (*Task, []*Task, error) {
	for index, task := range tasks {
		if task.canRun() {
			return task, remove(tasks, index), nil
		}
	}
	return nil, tasks, errors.New("no task available")
}

func confirmThatAllTasksAreCompleted(tasks []*Task) bool {
	for _, task := range tasks {
		if task.isNotCompleted() {
			return false
		}
	}
	return true
}

// func executeTask(waitGroup *sync.WaitGroup, task *Task) {
// 	rand.Seed(time.Now().UnixNano())
//     min := 1
//     max := 5
//     runTime := (rand.Intn(max - min + 1) + min)

// 	time.Sleep(time.Duration(runTime) * time.Second)
// }

// efficient removal of item from list (does not preserve the order of the list)
func remove(s []*Task, i int) []*Task {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
