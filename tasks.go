package main

import (
	"fmt"
	"log"
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
	taskState TaskState
	taskType TaskType

	input    []string
	output   []string

	dependsOn []*Task
}

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
	common := fmt.Sprintf("  Completed? %v\n  Can run? %v\n  Depends on %v others",
		t.taskState == Complete, t.canRun(), len(t.dependsOn))

	switch t.taskType {
	case Transcode:
		return fmt.Sprintf("Transcode\n%v\n  From %v\n  To %v",
			common, t.input[0], t.output[0])
	case FixAudio:
		return fmt.Sprintf("FixAudio\n%v\n  From %v\n  To %v",
			common, t.input[0], t.output[0])
	case Concatenate:
		return fmt.Sprintf("Concatenate\n%v",
			common)
	default:
		return "?"
	}
}

func (t *Task) addDependant(other *Task) {
	t.dependsOn = append(t.dependsOn, other)
}

func NewFixAudioTask(fileIn string, fileOut string) *Task {
	return &Task{Pending, FixAudio, []string{fileIn}, []string{fileOut}, []*Task{}}
}

func NewTranscodeTask(fileIn string, fileOut string) *Task {
	return &Task{Pending, Transcode, []string{fileIn}, []string{fileOut}, []*Task{}}
}

func NewConcatenateTask(fileOut string, dependsOn []*Task) *Task {
	return &Task{Pending, Concatenate, nil, []string{fileOut}, dependsOn}
}

func GenerateTasks() []*Task {

	tasks := []*Task{}

	paths := ScanPaths()

	for path, files := range paths {
		log.Printf("%s : %v", path, files)

		concatenateDependees := []*Task{}

		for _, fileIn1 := range files {
			fileOut1 := fileIn1 + "_fixedaudio"
			task1 := NewFixAudioTask(fileIn1, fileOut1)
			tasks = append(tasks, task1)

			fileIn2 := fileOut1
			fileOut2 := fileIn2 + "_transcoded"
			task2 := NewTranscodeTask(fileIn2, fileOut2)
			task2.addDependant(task1)
			tasks = append(tasks, task2)

			concatenateDependees = append(concatenateDependees, task2)
		}

		finalFileOut := path + "_concatenated"
		finalTask := NewConcatenateTask(finalFileOut, concatenateDependees)
		tasks = append(tasks, finalTask)
	}

	return tasks
}
