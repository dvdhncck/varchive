package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

const maxMessages = 8

type Stats struct {
	isReady        bool
	tasksCompleted int
	tasksRemaining int
	totalBytes     float64
	totalRunTimeInSeconds float64
}

type Estimator struct {
	totalInputSize          [TaskTypeCount]float64 // cumulative for all completed tasks of this type
	totalRunTime            [TaskTypeCount]float64 // ditto
	estimatedBytesPerSecond [TaskTypeCount]float64
}

type WorkerInfo struct {
	task             *Task
}

type Monitor struct {
	lock      sync.Mutex
	workers   []*WorkerInfo
	display   *Display
	messages  [maxMessages]*string
	stats     Stats
	estimator Estimator
}

func (m *Monitor) ShutdownCleanly() {
	log.Printf("Clean shutdown requested")
	m.display.Close()
}

func (m *Monitor) NotifyTaskBegins(task *Task) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.addMessage(fmt.Sprintf("Running task %v", task.BriefString()))
	m.workers = append(m.workers, &WorkerInfo{task})
}

func (m *Monitor) NotifyTaskEnds(task *Task) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for index, workerInfo := range m.workers {
		if workerInfo.task.id == task.id {
			task := workerInfo.task
			task.runTimeInSeconds = time.Since(task.startTime).Seconds()
			
			bytesPerSecond := float64(task.inputSize) / float64(task.runTimeInSeconds)

			m.stats.tasksCompleted++
			m.stats.tasksRemaining--

			m.updateEstimates(workerInfo)

			m.addMessage(fmt.Sprintf("Completed task %s in %v (%v/s)",
				task.BriefString(), niceTime(task.runTimeInSeconds), niceSize(int64(bytesPerSecond))))

			// rebuild the slice with the index'th element removed
			m.workers = append(m.workers[:index], m.workers[index+1:]...)
			return
		}
	}
}

func NewMonitor(tasks []*Task) *Monitor {
	workers := []*WorkerInfo{}
	messages := [maxMessages]*string{}
	estimator := Estimator{[TaskTypeCount]float64{}, [TaskTypeCount]float64{}, [TaskTypeCount]float64{}}
	stats := Stats{false, len(tasks), 0, 0, 0}

	m := &Monitor{sync.Mutex{}, workers, NewDisplay(), messages, stats, estimator}

	for i := 0; i < maxMessages; i++ {
		text := "..."
		m.messages[i] = &text
	}

	return m
}

func (m *Monitor) Start() {
	m.display.Init()
	defer m.display.Close()

	go func(m *Monitor) {

		startTime := time.Now()

		for {
			m.lock.Lock()

			m.display.Clear()
			m.display.Write(fmt.Sprintf("%d beavers employed, %d tasks completed, %d remaining\n",
				len(m.workers), m.stats.tasksCompleted, m.stats.tasksRemaining))
			m.display.Write("Task     Purpose       Size          Run time        ETA")
			m.display.Write("-------+------------+-------------+---------------+----------------")

			for _, workerInfo := range m.workers {
				task := workerInfo.task
				task.runTimeInSeconds = time.Since(task.startTime).Seconds()
				remaining := m.EstimateTimeRemaining(workerInfo)
				m.display.Write(fmt.Sprintf("%4d    %-13s%11s   %-16s%-16s",
					task.id, 
					task.TaskType(), 
					task.Size(),
					niceTime(task.runTimeInSeconds), 
					niceTime(remaining)))
			}

			runTime := time.Since(startTime).Seconds()
			m.display.Write(fmt.Sprintf("\nElapsed: %s", niceTime(runTime)))

			m.display.Write("\nRecently:")
			for _, message := range m.messages {
				m.display.Write(fmt.Sprintf("   %s", *message))
			}
			m.display.Flush()
			m.lock.Unlock()

			time.Sleep(1 * time.Second)
		}
	}(m)
}

func niceTime(seconds float64) string {
	if seconds <= 0 {
		return "---:--:--"
	}

	const spm = 60
	const sph = 60 * 60
	h, m, s := 0, 0, int64(seconds)
	for s > sph {
		h++
		s -= sph
	}
	for s > spm {
		m++
		s -= spm
	}
	return fmt.Sprintf("%03d:%02d:%02d", h, m, s)
}

func (m *Monitor) addMessage(message string) {
	for i := maxMessages - 1; i > 0; i-- {
		m.messages[i] = m.messages[i-1]
	}
	m.messages[0] = &message
	//m.messages[maxMessages - 1] = "..."
}

func (m *Monitor) countWorkersOfType(taskType TaskType) int {
	count := 0
	for _, workerInfo := range m.workers {
		if workerInfo.task.taskType == taskType {
			count++
		}
	}
	return count
}

// called when a worker completes a task (and before any new task is scheduled)
func (m *Monitor) updateEstimates(workerInfo *WorkerInfo) {
	e := m.estimator
	
	workersOfThisType := m.countWorkersOfType(workerInfo.task.taskType)

	taskType := workerInfo.task.taskType

	e.totalInputSize[taskType] += float64(workerInfo.task.inputSize)
	e.totalRunTime[taskType] += float64(workerInfo.task.runTimeInSeconds)	
	
	ebpsAllWorkers := e.totalInputSize[taskType] / e.totalRunTime[taskType]

	e.estimatedBytesPerSecond[taskType] = ebpsAllWorkers / float64(workersOfThisType)
}

func (m *Monitor) EstimateBytesPerSecond(taskType TaskType) float64 {
	return m.estimator.estimatedBytesPerSecond[taskType]
}

func (m *Monitor) EstimateTimeRemaining(workerInfo *WorkerInfo) float64 {
	workersOfThisType := m.countWorkersOfType(workerInfo.task.taskType)

	bpsForThisWorker := m.EstimateBytesPerSecond(workerInfo.task.taskType) / float64(workersOfThisType)

	estimatedTotalTimeInSeconds := float64(workerInfo.task.inputSize) / bpsForThisWorker
	remainingTimeInSeconds := estimatedTotalTimeInSeconds - workerInfo.task.runTimeInSeconds
	
	if remainingTimeInSeconds < 0 {
	 	remainingTimeInSeconds = 0
	}

	return remainingTimeInSeconds
}
