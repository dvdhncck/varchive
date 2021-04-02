package main

import (
	"errors"
	"log"
	"sync"
	"time"
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
			
			waitGroup.Add(1)
			go func() {
				executeTask(waitGroup, task)
				<-guard // consumes an item from the channel
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
//   return (nil task, the original list, an error)
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

func executeTask(waitGroup *sync.WaitGroup, task *Task) {
	defer waitGroup.Done()

	log.Printf("Running task %v", task)
	task.taskState = Running

	switch task.taskType {
	case Transcode:
		time.Sleep(3 * time.Second)
	case FixAudio:
		time.Sleep(1 * time.Second)
	case Concatenate:
		time.Sleep(2 * time.Second)
	}

	task.taskState = Complete
	log.Printf("Completed task %v", task)
}

// efficient removal of item from list (does not preserve the order of the list)
func remove(s []*Task, i int) []*Task {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
