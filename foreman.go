package varchive

import (
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

func GetBusy() {

	if settings.reportSizes {
		ReportSizes()
	} else {
		tasks := GenerateTasks()

		SortTasks(tasks)

		if settings.verbose {
			for _, task := range tasks {
				Log("%v\n", task)
			}
		}
		ScheduleTasks(NewTimer(), tasks)
	}
}

func ReportSizes() {
	paths := ScanPaths()

	widths := NewHisto()
	heights := NewHisto()

	for path, files := range paths {
		if settings.verbose {
			Log("Path: %s has %d files", path, len(files))
		}

		for _, file := range files {
			width, height, err := GetVideoInfo(file.path)
			if err == nil {
				Log("%s %dx%d", file.path, width, height)
				widths.Add(width)
				heights.Add(height)
			} else {
				Log("%s %s", file.path, err.Error())
			}
		}
	}

	Log("\n\n  Widths:\n%v\n  Heights:\n%v", widths, heights)
}

func GenerateTasks() []*Task {

	if settings.verbose {
		Log("Generating tasks")
	}

	createOutputRootIfRequired()

	tasks := []*Task{}

	paths := ScanPaths()

	for path, files := range paths {
		if settings.verbose {
			Log("Path: %s has %d file(s)", path, len(files))
		}

		concatenateDependees := []*Task{}

		for _, file := range files {

			fileIn := file
			fileOut := makeTemporaryFile(".mp4")
			transcodeTask := NewTranscodeTask(file, fileOut)
			transcodeTask.inputSize = fileIn.size
			tasks = append(tasks, transcodeTask)

			if settings.fixAudio {
				existingExtension := getFileExtension(fileIn.path)
				fixAudioFileOut := makeTemporaryFile(existingExtension)
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

	Log("%d tasks to be scheduled", len(tasks))

	return tasks
}

// put the FixAudio tasks at the front of the queue, ordered by the
// size of their inputs, then Transcode tasks, again, order by input size
//
func SortTasks(tasks []*Task) {
	sort.Slice(tasks, func(i1, i2 int) bool { return tasks[i1].lessThan(tasks[i2]) })
}
