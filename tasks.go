package varchive

import (
	"log"
	"path/filepath"
	"sort"
)

func NewFixAudioTask(fileIn *FileWithSize, fileOut string) *Task {
	return NewTask(FixAudio, fileIn.path, fileOut, fileIn.size)
}

func NewTranscodeTask(fileIn *FileWithSize, fileOut string) *Task {
	return NewTask(Transcode, fileIn.path, fileOut, fileIn.size)
}

func NewConcatenateTask(fileOut string, dependsOn []*Task) *Task {
	task := NewTask(Concatenate, "", fileOut, 0)
	task.addDependants(dependsOn)
	return task
}

func GenerateTasks() []*Task {

	createOutputRootIfRequired()

	tasks := []*Task{}

	paths := ScanPaths()

	for path, files := range paths {
		if settings.verbose {
			log.Printf("%s : %v", path, files)
		}

		concatenateDependees := []*Task{}

		for _, file := range files {

			fileIn := file
			fileOut := makeTemporaryFile(".mp4")
			transcodeTask := NewTranscodeTask(file, fileOut)
			transcodeTask.inputSize = fileIn.size
			tasks = append(tasks, transcodeTask)

			if settings.fixAudio {
				fixAudioFileOut := makeTemporaryFile(".mp4")
				fixAudioTask := NewFixAudioTask(file, fixAudioFileOut)
				fixAudioTask.inputSize = fileIn.size

				transcodeTask.fileIn = fixAudioFileOut
				transcodeTask.addDependant(fixAudioTask)

				tasks = append(tasks, fixAudioTask)
			}

			concatenateDependees = append(concatenateDependees, transcodeTask)
		}

		finalFileName := lastBitOfPath(path)
		finalFileOut := filepath.Join(settings.outputRoot, sanitisePath(finalFileName)+".mp4")

		failIfConcatenationFileAlreadyExists(finalFileOut)

		finalTask := NewConcatenateTask(finalFileOut, concatenateDependees)
		tasks = append(tasks, finalTask)
	}

	return tasks
}

// put the FixAudio tasks at the front of the queue, ordered by the
// size of their inputs, then Transcode tasks, again, order by input size
//
func SortTasks(tasks []*Task) {
	sort.Slice(tasks, func(i1, i2 int) bool { return tasks[i1].lessThan(tasks[i2]) })
}
