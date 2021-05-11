package varchive

import (
	"fmt"
	"sync"
)

const maxMessages = 16

var ellipsis = string(`...`)

type Stats struct {
	tasksRemaining        int
	tasksCompleted        int
	totalBytes            float64
	totalRunTimeInSeconds float64
}

type Monitor struct {
	timer       Timer
	lock        sync.Mutex
	allTasks    []*Task
	activeTasks []*Task
	display     Display
	messages    [maxMessages]*string
	stats       Stats
	estimator   *Estimator
}

func (m *Monitor) ShutdownCleanly() {
	Log("Clean shutdown requested")
	if settings.liveDisplay {
		m.display.Close()
	}
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

func NewMonitor(timer Timer, estimator *Estimator, allTasks []*Task, display Display) *Monitor {
	m := &Monitor{
		timer:       timer,
		lock:        sync.Mutex{},
		allTasks:    allTasks,
		activeTasks: []*Task{},
		display:     display,
		messages:    [maxMessages]*string{},
		stats:       Stats{len(allTasks), 0, 0, 0},
		estimator:   estimator,
	}

	for i := 0; i < maxMessages; i++ {
		text := "..."
		m.messages[i] = &text
	}

	return m
}

func (m *Monitor) Start() {
	if settings.liveDisplay {
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
}

func (m *Monitor) tick(runTimeInSecond float64) {
	m.lock.Lock()
	defer m.lock.Unlock()

	totalRemainingTimeInSeconds := m.estimator.EstimateRemainingRunTime(m.allTasks)

	for _, task := range m.activeTasks {
		workersOfThisType := m.countWorkersOfType(task.taskType)
		task.runTimeInSeconds = m.timer.SecondsSince(task.startTimestamp)
		task.estimatedRemainingTimeInSeconds = m.estimator.EstimateTimeRemaining(task, workersOfThisType)

		totalRemainingTimeInSeconds += task.estimatedRemainingTimeInSeconds
	}

	m.display.Clear()

	m.display.Write(fmt.Sprintf("Elapsed: %s, remaining: %s\n", niceTime(runTimeInSecond), niceTime(totalRemainingTimeInSeconds)))
	m.display.Write(fmt.Sprintf("%d beavers employed, %d tasks completed, %d remaining\n", len(m.activeTasks), m.stats.tasksCompleted, m.stats.tasksRemaining))

	m.display.Write("Task     Purpose       Size          Run time        ETA")
	m.display.Write("-------+------------+-------------+---------------+----------------")

	for _, task := range m.activeTasks {
		m.display.Write(fmt.Sprintf("%4d    %-13s%11s   %-16s%-16s",
			task.id,
			task.TaskType(),
			task.Size(),
			niceTime(task.runTimeInSeconds),
			niceTime(task.estimatedRemainingTimeInSeconds)))
	}

	m.display.Write("\nRecently:")
	for _, message := range m.messages {
		m.display.Write(fmt.Sprintf("   %s", *message))
	}

	m.display.Flush()

}

func (m *Monitor) addMessage(message string) {
	Log(message) // the permanent record

	last := maxMessages - 1
	for i := last - 1; i > 0; i-- {
		m.messages[i] = m.messages[i-1]
	}
	m.messages[last] = &ellipsis
	m.messages[0] = &message
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

	workersOfThisType := m.countWorkersOfType(task.taskType)
	m.estimator.UpdateEstimates(task, workersOfThisType)

	m.addMessage(fmt.Sprintf("Estimates computed: FixAudio %s/s, Transcode %s/s",
		niceSize(int64(m.estimator.EstimateBytesPerSecond(FixAudio))),
		niceSize(int64(m.estimator.EstimateBytesPerSecond(Transcode)))))
}

// returns the estimate of how much longer this task will take
// or -Inf if the task is taking longer than expected
// or +Inf if there is no data available to make the estimation

func (m *Monitor) EstimateTimeRemaining(task *Task) float64 {
	workersOfThisType := m.countWorkersOfType(task.taskType)
	return m.estimator.EstimateTimeRemaining(task, workersOfThisType)
}
