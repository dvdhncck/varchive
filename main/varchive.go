package main

import (
	"log"
	"os"
	va "davidhancock.com/varchive"
)

func main() {
	
	initialiseLogging()

	va.ParseArguments() // guarantees that arguments are acceptable

	tasks := va.GenerateTasks()

	va.SortTasks(tasks)

	if va.IsVerbose() {
		for _, task := range tasks {
			log.Printf("%v\n\n", task)
		}
	}

	va.ScheduleTasks(va.NewTimer(), tasks)
}

func initialiseLogging() {
	// file, err := os.OpenFile("Scanner.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.SetOutput(os.Stdout)
}

