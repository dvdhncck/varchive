package main

import (
	va "davidhancock.com/varchive"
)

func main() {

	va.ParseArguments() // guarantees that arguments are acceptable
	
	va.InitialiseLogging()

	va.GetBusy()

}
