package varchive

import (
	"fmt"
	"log"
	"os"
)

func Log(message string, variables ...interface{}) {
	if settings.consoleOutputAllowed || settings.logToFile != "" {
		if len(variables) == 0 {
			log.Println(message)
		} else {
			log.Println(fmt.Sprintf(message, variables...))
		}
	}
}

func InitialiseLogging() {

	if settings.logToFile != "" {
		file, err := os.OpenFile(settings.logToFile, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fatal(err.Error())

		} else {
			log.SetOutput(file)
		}
	} else {
		log.SetOutput(os.Stdout)
		if settings.verbose {
			log.Println("Logging to StdOut")
		}
	}
}
