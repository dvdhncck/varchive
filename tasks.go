package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
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
	// note that the ordering defines the 'priority' of a task when we schedule them
	FixAudio    TaskType = 1
	Transcode   TaskType = 2
	Concatenate TaskType = 3
)

type Task struct {
	id int

	inputSize int64

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

func (t1 *Task) lessThan(t2 *Task) bool {
	if t1.taskType == t2.taskType {
		return t1.inputSize > t2.inputSize // biggest jobs ones come first
	} else {
		return t1.taskType < t2.taskType // higher priority comes first, e.g. FixAudio before Transcode before Concatenate
	}
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
		if t.canRun() {
			state = "Runnable"
		} else {
			state = fmt.Sprintf("Pending (%d dependees)", len(t.dependsOn))
		}
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

func (t *Task) BriefString() string {

	switch t.taskType {
	case Transcode:
		return fmt.Sprintf("#%d Transcode %v", t.id, niceSize(t.inputSize))
	case FixAudio:
		return fmt.Sprintf("#%d FixAudio %v", t.id, niceSize(t.inputSize))
	case Concatenate:
		return fmt.Sprintf("#%d Concatenate %d items", t.id, len(t.dependsOn))
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
	task := Task{taskId, 0, 0, Pending, taskType, fileIn, fileOut, []*Task{}}
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

		for _, file := range files {

			fileIn := file
			fileOut := makeTemporaryFile(".mp4")
			transcodeTask := NewTranscodeTask(fileIn.path, fileOut)
			transcodeTask.inputSize = fileIn.size
			tasks = append(tasks, transcodeTask)

			if settings.fixAudio {
				fixAudioFileOut := makeTemporaryFile(".mp4")
				fixAudioTask := NewFixAudioTask(fileIn.path, fixAudioFileOut)
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
