package varchive

import (
	"math"
)

type Estimator struct {
	totalInputSize          [TaskTypeCount]float64 // cumulative for all completed tasks of this type
	totalRunTime            [TaskTypeCount]float64 // ditto
	estimatedBytesPerSecond [TaskTypeCount]float64
}

func NewEstimator() *Estimator {
	estimator := Estimator{[TaskTypeCount]float64{}, [TaskTypeCount]float64{}, [TaskTypeCount]float64{}}
	// these estimates that came from a very long encoding session (on skink in March 2021)
	// 30.6 MiB kps, Transcode 100.0 KiB
	estimator.estimatedBytesPerSecond[FixAudio] = 30.6 * 1000 * 1000
	estimator.estimatedBytesPerSecond[Transcode] = 100 * 1000
	return &estimator
}

// called when a worker completes a task (and before any new task is scheduled)
func (e *Estimator) UpdateEstimates(task *Task, workersOfThisType int) {

	// how accurate was the last estimate?
	estimated := e.EstimateRuntime(task, workersOfThisType)
	actual := task.runTimeInSeconds
	error := math.Abs(estimated-actual) / actual // bigger values are worse
	Log("Estimation error: %.2f  (e=%f, a=%f)", error, estimated, actual)

	taskType := task.taskType

	e.totalInputSize[taskType] += float64(task.inputSize)
	e.totalRunTime[taskType] += float64(task.runTimeInSeconds)

	ebpsAllWorkers := e.totalInputSize[taskType] / e.totalRunTime[taskType]

	e.estimatedBytesPerSecond[taskType] = ebpsAllWorkers * float64(workersOfThisType)
}

func (e *Estimator) EstimateBytesPerSecond(taskType TaskType) float64 {
	return e.estimatedBytesPerSecond[taskType]
}

func (e *Estimator) EstimateRuntime(task *Task, workersOfThisType int) float64 {

	bpsForThisWorker := e.EstimateBytesPerSecond(task.taskType) / float64(workersOfThisType)

	estimatedTotalTimeInSeconds := float64(task.inputSize) / bpsForThisWorker

	return estimatedTotalTimeInSeconds
}

// returns the estimate of how much longer this task will take
// or -Inf if the task is taking longer than expected
// or +Inf if there is no data available to make the estimation
func (e *Estimator) EstimateTimeRemaining(task *Task, workersOfThisType int) float64 {
	remainingTimeInSeconds := e.EstimateRuntime(task, workersOfThisType) - task.runTimeInSeconds

	if remainingTimeInSeconds < 0 {
		return math.Inf(-1)
	}

	return remainingTimeInSeconds
}

func (e *Estimator) EstimateRemainingRunTime(tasks []*Task) float64 {
	totalEstimatedTime := float64(0)
	for _, task := range(tasks) {
		if task.IsNotCompleted() {
			bps := e.EstimateBytesPerSecond(task.taskType)
			estimatedTime := float64(task.inputSize) / bps
			totalEstimatedTime += estimatedTime
		}
	}
	return float64(totalEstimatedTime)
}