package varchive

import (
	"fmt"
	"math"
	"sync"
)

const maxMessages = 8

type Stats struct {
	tasksRemaining        int
	tasksCompleted        int
	totalBytes            float64
	totalRunTimeInSeconds float64
}

type Estimator struct {
	totalInputSize          [TaskTypeCount]float64 // cumulative for all completed tasks of this type
	totalRunTime            [TaskTypeCount]float64 // ditto
	estimatedBytesPerSecond [TaskTypeCount]float64
}

type Monitor struct {
	timer       Timer
	lock        sync.Mutex
	activeTasks []*Task
	display     Display
	messages    [maxMessages]*string
	stats       Stats
	estimator   Estimator
}

func (m *Monitor) ShutdownCleanly() {
	Log("Clean shutdown requested")
	m.display.Close()
}

func (m *Monitor) NotifyTaskBegins(task *Task) {
	task.startTimestamp = m.timer.Now()
	m.lock.Lock()
	defer m.lock.Unlock()
	m.addMessage(fmt.Sprintf("Running task %v", task.BriefString()))
	m.activeTasks = append(m.activeTasks, task)
}

func (m *Monitor) NotifyTaskEnds(task *Task) {
	task.runTimeInSeconds = m.timer.SecondsSince(task.startTimestamp)

	m.lock.Lock()
	defer m.lock.Unlock()

	for index, t := range m.activeTasks {
		if t.id == task.id {
						
			bytesPerSecondForTask := int64(float64(task.inputSize) / float64(task.runTimeInSeconds))

			m.addMessage(fmt.Sprintf("Completed task %s in %v (%v/s)",
				task.BriefString(), niceTime(task.runTimeInSeconds), niceSize(bytesPerSecondForTask)))

			m.stats.tasksCompleted++
			m.stats.tasksRemaining--

			m.updateEstimates(task)

			// rebuild the activeTasks list with the index'th element removed
			m.activeTasks = append(m.activeTasks[:index], m.activeTasks[index+1:]...)
			return
		}
	}
}

func NewMonitor(timer Timer, tasks []*Task, display Display) *Monitor {
	activeTasks := []*Task{}
	messages := [maxMessages]*string{}
	estimator := Estimator{[TaskTypeCount]float64{}, [TaskTypeCount]float64{}, [TaskTypeCount]float64{}}
	stats := Stats{len(tasks), 0, 0, 0}

	m := &Monitor{timer, sync.Mutex{}, activeTasks, display, messages, stats, estimator}

	for i := 0; i < maxMessages; i++ {
		text := "..."
		m.messages[i] = &text
	}

	return m
}

func (m *Monitor) Start() {
	m.display.Init()
	m.display.Clear()

	go func(monitor *Monitor) {

		startTimestamp := monitor.timer.Now()

		for {
			runTime := monitor.timer.SecondsSince(startTimestamp)
			monitor.tick(runTime)
			monitor.timer.MilliSleep(900)
		}
	}(m)
}

func (m *Monitor) tick(runTimeInSecond float64) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.display.Clear()

	m.display.Write(fmt.Sprintf("%d beavers employed, %d tasks completed, %d remaining\n",
		len(m.activeTasks), m.stats.tasksCompleted, m.stats.tasksRemaining))
	m.display.Write("Task     Purpose       Size          Run time        ETA")
	m.display.Write("-------+------------+-------------+---------------+----------------")

	for _, task := range m.activeTasks {

		task.runTimeInSeconds = m.timer.SecondsSince(task.startTimestamp)

		m.display.Write(fmt.Sprintf("%4d    %-13s%11s   %-16s%-16s",
			task.id,
			task.TaskType(),
			task.Size(),
			niceTime(task.runTimeInSeconds),
			niceTime(m.EstimateTimeRemaining(task))))
	}

	m.display.Write(fmt.Sprintf("\nElapsed: %s", niceTime(runTimeInSecond)))

	m.display.Write("\nRecently:")
	for _, message := range m.messages {
		m.display.Write(fmt.Sprintf("   %s", *message))
	}

	m.display.Flush()
}

func (m *Monitor) addMessage(message string) {
	Log(message) // the permanent record
	
	for i := maxMessages - 1; i > 0; i-- {
		m.messages[i] = m.messages[i-1]
	}
	m.messages[0] = &message
	//m.messages[maxMessages - 1] = "..."
}

func (m *Monitor) countWorkersOfType(taskType TaskType) int {
	count := 0
	for _, task := range m.activeTasks {
		if task.taskType == taskType {
			count++
		}
	}
	return count
}

// called when a worker completes a task (and before any new task is scheduled)
func (m *Monitor) updateEstimates(task *Task) {
	e := &m.estimator


	// how accurate was the last estimate?
	estimated := m.EstimateRuntime(task)
	actual := task.runTimeInSeconds
	error := math.Abs(estimated - actual) / actual   // bigger values are worse
	Log("Estimation error: %.2f  (e=%f, a=%f)", error, estimated, actual)
	
	workersOfThisType := m.countWorkersOfType(task.taskType)

	taskType := task.taskType

	e.totalInputSize[taskType] += float64(task.inputSize)
	e.totalRunTime[taskType] += float64(task.runTimeInSeconds)

	ebpsAllWorkers := e.totalInputSize[taskType] / e.totalRunTime[taskType]

	e.estimatedBytesPerSecond[taskType] = ebpsAllWorkers * float64(workersOfThisType)

	m.addMessage(fmt.Sprintf("Estimates computed: FixAudio %s kps, Transcode %s kps",
		niceSize(int64(e.estimatedBytesPerSecond[FixAudio])),
		niceSize(int64(e.estimatedBytesPerSecond[Transcode]))))
}

func (m *Monitor) EstimateBytesPerSecond(taskType TaskType) float64 {
	return m.estimator.estimatedBytesPerSecond[taskType]
}

func(m *Monitor) EstimateRuntime(task *Task) float64 {
	workersOfThisType := m.countWorkersOfType(task.taskType)

	bpsForThisWorker := m.EstimateBytesPerSecond(task.taskType) / float64(workersOfThisType)

	estimatedTotalTimeInSeconds := float64(task.inputSize) / bpsForThisWorker
	
	return estimatedTotalTimeInSeconds
}

// returns the estimate of how much longer this task will take
// or -Inf if the task is taking longer than expected
// or +Inf if there is no data available to make the estimation
func (m *Monitor) EstimateTimeRemaining(task *Task) float64 {
	remainingTimeInSeconds := m.EstimateRuntime(task) - task.runTimeInSeconds

	if remainingTimeInSeconds < 0 {
		return math.Inf(-1)
	}

	return remainingTimeInSeconds
}
