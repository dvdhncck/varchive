package main

import (
	"fmt"
	"log"
	"os"
)

func fatal(message string) {

	fmt.Println(message)
	os.Exit(1)
}

var settings *Settings

func main() {

	initialiseLogging()

	settings = ParseArguments() // guarantees that arguments are acceptable

	tasks := GenerateTasks()

	if settings.verbose {
		for _, task := range tasks {
			log.Printf("%v\n\n", task)
		}
	}

	ScheduleTasks(tasks)
}

func initialiseLogging() {
	// file, err := os.OpenFile("Scanner.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.SetOutput(os.Stdout)
}
