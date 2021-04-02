package main

import (
	"errors"
	"log"
	"time"
)

func ScheduleTasks(tasks []*Task) {

	completed := false

	for !completed {

		task, remainingTasks, err := findFirstRunnableTask(tasks)

		if err == nil {
			tasks = remainingTasks
			executeTask(task)
		} else {
			// no incompleted tasks left, we are done
			completed = true
		}
	}
}

// if there is a task ready to run
//   return (the task, the list with the task now removed, nil error)
// else
//   return (nil task, the original list, an error)
func findFirstRunnableTask(tasks []*Task) (*Task, []*Task, error) {
	for index, task := range tasks {
		if task.canRun(){
			return task, remove(tasks, index), nil
		}
	}
	return nil, tasks, errors.New("no task available")
}

func executeTask(task *Task) {
	log.Printf("Running task %v", task)
	task.taskState = Running
	time.Sleep(2 * time.Second)
	task.taskState = Complete
	log.Printf("Completed task %v", task)
}

// efficient removal of item from list (does not preserve the order of the list)
func remove(s []*Task, i int) []*Task {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
