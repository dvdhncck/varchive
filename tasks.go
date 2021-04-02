package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type TaskState int

const (
	Pending  TaskState = 1
	Running  TaskState = 2
	Complete TaskState = 3
)

type TaskType int

const (
	Transcode   TaskType = 1
	FixAudio    TaskType = 2
	Concatenate TaskType = 3
)

type Task struct {
	id int

	runTime time.Duration

	taskState TaskState
	taskType  TaskType

	fileIn  string
	fileOut string

	dependsOn []*Task
}

var taskId = 1

func (t *Task) isNotCompleted() bool {
	return t.taskState != Complete
}

func (t *Task) canRun() bool {
	if t.taskState != Pending {
		return false
	}

	for _, d := range t.dependsOn {
		if d.taskState != Complete {
			return false
		}
	}

	return true
}

func (t *Task) String() string {
	
	state := "?"
	switch t.taskState {
	case Running:
		state = "Running"
	case Pending:
		if t.canRun() { state = "Runnable" } else { state = fmt.Sprintf("Pending (%d dependees)", len(t.dependsOn)) }
	case Complete:
		state = "Complete"
	}

	switch t.taskType {
	case Transcode:
		return fmt.Sprintf("%d\n  Type: Transcode\n  State: %s\n  From: %s\n  To: %s",
			t.id, state, t.fileIn, t.fileOut)
	case FixAudio:
		return fmt.Sprintf("%d\n  Type: FixAudio\n  State: %s\n  From: %s\n  To: %s",
			t.id, state, t.fileIn, t.fileOut)
	case Concatenate:
		return fmt.Sprintf("%d\n  Type: Concatenate\n  State: %s\n  To: %s",
			t.id, state, t.fileOut)
	default:
		return "?"
	}
}

func (t *Task) addDependant(other *Task) {
	t.dependsOn = append(t.dependsOn, other)
}

func (t *Task) addDependants(others []*Task) {
	t.dependsOn = append(t.dependsOn, others...)
}


func NewTask(taskType TaskType, fileIn string, fileOut string) *Task {
	task := Task{taskId, 0, Pending, taskType, fileIn, fileOut, []*Task{}}
	taskId += 1
	return &task
}

func NewFixAudioTask(fileIn string, fileOut string) *Task {
	return NewTask(FixAudio, fileIn, fileOut)
}

func NewTranscodeTask(fileIn string, fileOut string) *Task {
	return NewTask(Transcode, fileIn, fileOut)
}

func NewConcatenateTask(fileOut string, dependsOn []*Task) *Task {
	task := NewTask(Concatenate, "", fileOut)
	task.addDependants(dependsOn)
	return task
}

func GenerateTasks() []*Task {
	
	createOutputRootIfRequired()

	tasks := []*Task{}

	paths := ScanPaths()

	for path, files := range paths {
		log.Printf("%s : %v", path, files)

		concatenateDependees := []*Task{}

		for _, fileIn1 := range files {
			fileOut1 := makeTemporaryFile(".mp4")
			task1 := NewFixAudioTask(fileIn1, fileOut1)
			tasks = append(tasks, task1)

			fileIn2 := fileOut1
			fileOut2 := makeTemporaryFile(".mp4")
			task2 := NewTranscodeTask(fileIn2, fileOut2)
			task2.addDependant(task1)
			tasks = append(tasks, task2)

			concatenateDependees = append(concatenateDependees, task2)
		}
	
		finalFileName := lastBitOfPath(path)
		finalFileOut := filepath.Join(settings.outputRoot, sanitisePath(finalFileName)+".mp4")
		//finalFileOut = makeTemporaryFile(".mp4")

		failIfConcatenationFileAlreadyExists(finalFileOut)


		finalTask := NewConcatenateTask(finalFileOut, concatenateDependees)
		tasks = append(tasks, finalTask)
	}

	return tasks
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
