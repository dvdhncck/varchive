package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"io/ioutil"
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

	SortTasks(tasks)

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

func createOutputRootIfRequired() {
	if _, err := os.Stat(settings.outputRoot); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(settings.outputRoot, 0755)
		}
	}
}

func failIfConcatenationFileAlreadyExists(path string) {
	if _, err := os.Stat(path); err == nil {
		fatal(fmt.Sprintf("%s exists, will not overwrite", path))
	}
}

func lastBitOfPath(path string) string {
	return filepath.Base(path)
}

// do we actually need this? if we use os.Command, daft filenames should not be a problem
func sanitisePath(path string) string {
	// re := regexp.MustCompile(`\s+`)
	// return re.ReplaceAllString(path, ``)
	return path
}

func makeTemporaryFile(extension string) string {
	file, err := ioutil.TempFile("", "varchive.*"+extension)
	if err != nil {
		fatal(err.Error())
	}
	defer os.Remove(file.Name())
	return file.Name()
}

func removeTemporaryFile(path string) {
	// ignore any errors (which will probably be "file not found")
	os.Remove(path)
}

func niceSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(bytes)/float64(div), "KMGTPE"[exp])
}
