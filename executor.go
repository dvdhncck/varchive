package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func doTranscode(task *Task) {
	args := []string{
		"--input", task.fileIn,
		"--output", task.fileOut,
		"--encoder", "x265",
		"--quality", fmt.Sprintf("%d", settings.quality),
		"--two-pass", "--turbo",
		"--aencoder", "copy",
		"2>&1"}

	invoke("HandBrakeCLI", args)
}

func doFixAudio(task *Task) {
	audioStream := makeTemporaryFile(".mp3")
	videoStream := makeTemporaryFile(".mp2")

	// demux the video stream (leave encoding as is)
	args := []string{
		"-i", task.fileIn,
		"-map", "0:0",
		"-codec", "copy",
		videoStream}

	invoke("ffmpeg", args)

	// demux the audio stream and transcode to mp3
	args = []string{
		"-i", task.fileIn,
		"-map", "0:1",
		"-codec", "mp3",
		audioStream}

	invoke("ffmpeg", args)

	// remux the audio and video streams into a new container
	args = []string{
		"-i", videoStream,
		"-i", audioStream,
		"-map", "0:v:0",
		"-map", "1:a:0",
		"-shortest",
		task.fileOut}

	invoke("ffmpeg", args)

	removeTemporaryFile(audioStream)
	removeTemporaryFile(videoStream)
}

func doConcatenate(task *Task) {
	listFile := makeTemporaryFile(".list")

	fileHandle, err := os.Create(listFile)
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
			log.Printf("Wrote concatenation list to %s", listFile)
		}
	} else {
		fatal(fmt.Sprintf("Could not open %s for the concatenation list", listFile))
	}

	args := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", listFile,
		"-c", "copy",
		task.fileOut}

	invoke("ffmpeg", args)

	// TODO: remove all of the fileOuts from the dependees

	removeTemporaryFile(listFile)
}

func ExecuteTask(task *Task) {
	switch task.taskType {
	case Transcode:
		doTranscode(task)

	case FixAudio:
		doFixAudio(task)

	case Concatenate:
		doConcatenate(task)
	}
}

func invoke(command string, args []string) {

	if settings.verbose {
		log.Printf("Invoking: %s %s", command, strings.Join(args, ` `))
	}

	if settings.dryRun {
		return
	} else {
		output, err := exec.Command(command, args...).Output()

		if err == nil {
			if settings.verbose {
				log.Printf("Return ok, stdout: %v", string(output))
			}
		} else {
			log.Fatal(fmt.Sprintf("Failed: %s %s\nErr: %v\nStdout: %v",
				command, strings.Join(args, ` `), err, string(output)))
		}
	}
}
