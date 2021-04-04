package main

import (
	"fmt"
	"log"
	"goncurses"
	"sync"
	"time"
)

const maxMessages = 8

type Stats struct {
	isReady                 bool
	estimatedBytesPerSecond float64
	totalBytes              float64
	totalComputeTimeSeconds float64
}

type WorkerInfo struct {
	workerId 		 int
	task     		 *Task
	runTimeInSeconds float64
}

type Monitor struct {
	lock           sync.Mutex
	bytesPerSecond float64
	workerInfo     []*WorkerInfo
	window         *goncurses.Window
	messages       [maxMessages]*string
	stats          Stats
}

func (m *Monitor) ShutdownCleanly() {
	log.Printf("Clean shutdown requested")
	m.closeTerminal()
}

func (m *Monitor) NotifyWorkerBegins(task *Task) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.addMessage(fmt.Sprintf("Running task %v", task.BriefString()))
	m.workerInfo = append(m.workerInfo, &WorkerInfo{1, task, 0})
}

func (m *Monitor) NotifyWorkerEnds(task *Task) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for index, workerInfo := range m.workerInfo {
		if workerInfo.task.id == task.id {

			runTimeInSeconds := time.Since(task.startTime).Seconds()
			bytesPerSecond := float64(task.inputSize) / float64(runTimeInSeconds)

			if task.taskType == Transcode {
				m.stats.totalBytes += float64(task.inputSize)
				m.stats.totalComputeTimeSeconds += runTimeInSeconds
				m.stats.estimatedBytesPerSecond = m.stats.totalBytes / m.stats.totalComputeTimeSeconds
				m.stats.isReady = true
				m.addMessage(fmt.Sprintf("Recomputed speed as %v/s", niceSize(int64(bytesPerSecond))))
			}

			m.addMessage(fmt.Sprintf("Completed task %s in %v (%v/s)",
				task.BriefString(), niceTime(runTimeInSeconds), niceSize(int64(bytesPerSecond))))

			// rebuild the slice with the index'th element removed
			m.workerInfo = append(m.workerInfo[:index], m.workerInfo[index+1:]...)
			return
		}
	}
}

func NewMonitor() *Monitor {
	m := &Monitor{sync.Mutex{},
		1204,
		[]*WorkerInfo{},
		&goncurses.Window{},
		[maxMessages]*string{},
		Stats{false, 0, 0, 0}}

	for i := 0; i < maxMessages; i++ {
		text := "..."
		m.messages[i] = &text
	}

	m.initTerminal()
	defer m.closeTerminal()

	go func(m *Monitor) {

		startTime := time.Now()

		for {
			m.lock.Lock()

			m.clearTerminal()
			m.writeToTerminal(fmt.Sprintf("%d beavers employed\n", len(m.workerInfo)))
			m.writeToTerminal("Task    Purpose      Size       Run time        ETA")
			m.writeToTerminal("-------+------------+-----------+---------------+----------------")
			for _, workerInfo := range m.workerInfo {
				task := workerInfo.task
				workerInfo.runTimeInSeconds = time.Since(task.startTime).Seconds()
				remaining := m.estimateTimeRemaining(workerInfo.runTimeInSeconds, task)
				m.writeToTerminal(fmt.Sprintf("%4d    %-13s%-12s%-16s%-16s",
					task.id, task.TaskType(), task.Size(), 
					niceTime(workerInfo.runTimeInSeconds), remaining))
			}

			runTime := time.Since(startTime).Seconds()
			m.writeToTerminal(fmt.Sprintf("\nElapsed: %s\nCompute: 000:00:00\nSpeedUp: 000", 
				niceTime(runTime)))

			m.writeToTerminal("\nRecently:")
			for _, message := range m.messages {
				m.writeToTerminal(fmt.Sprintf("   %s", *message))
			}
			m.flush()

			time.Sleep(1 * time.Second)
			m.lock.Unlock()
		}
	}(m)

	return m
}

func niceTime(seconds float64) string {
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

func (m *Monitor) estimateTimeRemaining(timeSoFarInSeconds float64, task *Task) string {
	if m.stats.isReady {
		estimatedTotalTimeInSeconds := float64(task.inputSize) / m.stats.estimatedBytesPerSecond
		remainingTimeInSeconds := estimatedTotalTimeInSeconds - timeSoFarInSeconds
		if remainingTimeInSeconds < 0 {
			remainingTimeInSeconds = 0
		}
		return niceTime(remainingTimeInSeconds)
	} else {
		return "---:--:--"
	}
}

func (m *Monitor) initTerminal() {
	window, err := goncurses.Init()

	if err != nil {
		fatal(err.Error())
	}

	m.window = window
}

func (m *Monitor) clearTerminal() {
	m.window.Erase()
}

func (m *Monitor) writeToTerminal(message string) {
	m.window.Println(message)
	//log.Println(message)
}

func (m *Monitor) closeTerminal() {
	goncurses.End()
}

func (m *Monitor) flush() {
	m.window.Refresh()
}

func (m *Monitor) addMessage(message string) {
	for i := maxMessages - 1; i > 0; i-- {
		m.messages[i] = m.messages[i-1]
	}
	m.messages[0] = &message
	//m.messages[maxMessages - 1] = "..."
}
