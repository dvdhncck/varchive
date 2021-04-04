package main

import (
	"fmt"
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

	startTime        time.Time
	runTimeInSeconds float64

	taskState TaskState
	taskType  TaskType

	fileIn  string
	fileOut string

	dependsOn []*Task
}

func (t *Task) addDependant(other *Task) {
	t.dependsOn = append(t.dependsOn, other)
}

func (t *Task) addDependants(others []*Task) {
	t.dependsOn = append(t.dependsOn, others...)
}

var taskId = 1

func NewTask(taskType TaskType, fileIn string, fileOut string) *Task {
	task := Task{taskId, 0, time.Time{}, 0, Pending, taskType, fileIn, fileOut, []*Task{}}
	taskId += 1
	return &task
}

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

func (t *Task) TaskType() string {
	switch t.taskType {
	case Transcode:
		return "Transcode"
	case FixAudio:
		return "FixAudio"
	case Concatenate:
		return "Concatenate"
	default:
		return "?"
	}
}

func (t *Task) Size() string {
	switch t.taskType {
	case Transcode, FixAudio:
		return niceSize(t.inputSize)
	case Concatenate:
		return fmt.Sprintf("%d items", len(t.dependsOn))
	default:
		return "?"
	}
}

func (t *Task) BriefString() string {
	return fmt.Sprintf("#%d %s %s", t.id, t.TaskType(), t.Size())
}
