package varchive

import (
	"log"
	"os"
)

func InitialiseLogging() {

	if settings.logToFile != "" {
		file, err := os.OpenFile("Scanner.log", os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fatal(err.Error())

		} else {
			log.SetOutput(file)
		}
	} else {
		log.SetOutput(os.Stdout)
	}
}
