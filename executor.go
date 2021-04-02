package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func ExecuteTask(task *Task) {

	switch task.taskType {
	case Transcode:

		args := []string{
			"--input", task.fileIn,
			"--output", task.fileOut,
			"--encoder", "x265",
			"--quality", fmt.Sprintf("%d", settings.quality),
			"--two-pass", "--turbo",
			"--aencoder", "copy",
			"2>&1"}

		output, err := invoke("HandBrakeCLI", args)

		if settings.verbose {
			log.Print(string(output))
		}

		if err != nil {
			log.Fatal(err)
		}

	case FixAudio:

		args := []string{
			task.fileIn,
			task.fileOut}

		output, err := invoke("cp", args)

		if settings.verbose {
			log.Print(string(output))
		}

		if err != nil {
			log.Fatal(err)
		}

	case Concatenate:

		listfile := makeTemporaryFile(".list")

		fileHandle, err := os.Create(listfile)
		if err == nil {
			//defer fileHandle.Close()
			writer := bufio.NewWriter(fileHandle)
			for _, dependee := range task.dependsOn {
				if settings.verbose {
					log.Printf("Add file %s", dependee.fileOut)
				}
				fmt.Fprintf(writer, "file %s\n", dependee.fileOut)
			}
			writer.Flush()
			fileHandle.Close()
			if settings.verbose {
				log.Printf("Wrote concatenation list to %s", listfile)
			}
		} else {
			fatal(fmt.Sprintf("Could not open %s for the concatenation list", listfile))
		}

		args := []string{
			"-f", "concat",
			"-safe", "0",
			"-i", listfile,
			"-c", "copy",
			task.fileOut}

		output, err := invoke("ffmpeg", args)

		if settings.verbose {
			log.Printf("Stdout: %v", string(output))
		}

		if err != nil {
			log.Print("Unhappy bunnies")
			log.Fatal(err)
		}
	}
}

func invoke(command string, args []string) ([]byte, error) {

	if settings.verbose {
		log.Printf("Invoking: %s %s", command, strings.Join(args, ` `))
	}

	if settings.dryRun {
		return []byte("[dry run]"), nil

	} else {
		return exec.Command(command, args...).Output()
	}
}
