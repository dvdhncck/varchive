package main

import (
	"log"
	va "davidhancock.com/varchive"
)

func main() {
	
	va.ParseArguments() // guarantees that arguments are acceptable

	va.InitialiseLogging()

	tasks := va.GenerateTasks()

	va.SortTasks(tasks)

	if va.IsVerbose() {
		for _, task := range tasks {
			log.Printf("%v\n\n", task)
		}
	}

	va.ScheduleTasks(va.NewTimer(), tasks)
}

