package varchive

import (
	"fmt"
	"time"
)

type TaskState int

const (
	Pending  TaskState = 0
	Running  TaskState = 1
	Complete TaskState = 2
)

type TaskType int
const TaskTypeCount int = 3
const (
	// note that the ordering defines the 'priority' of a task when we schedule them
	FixAudio    TaskType = 0
	Transcode   TaskType = 1
	Concatenate TaskType = 2
)

type Task struct {
	id int

	inputSize int64

	startTimestamp   Timestamp
	runTimeInSeconds float64
	estimatedRemainingTimeInSeconds float64

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

func NewTask(taskType TaskType, fileIn string, fileOut string, inputSize int64) *Task {
	task := Task{taskId, inputSize, time.Time{}, 0, 0, Pending, taskType, fileIn, fileOut, []*Task{}}
	taskId += 1
	return &task
}

func (t *Task) EstimatedRemainingTimeInSeconds() float64 {
	return t.estimatedRemainingTimeInSeconds
}

func (t *Task) IsRunning() bool {
	return t.taskState == Running
}

func (t *Task) IsPending() bool {
	return t.taskState == Pending
}

func (t *Task) IsNotCompleted() bool {
	return t.taskState != Complete
}

func (t *Task) MarkAsCompleted()  {
	t.taskState = Complete
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
