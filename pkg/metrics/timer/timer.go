package timer

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// Timer is the struct that keeps track of each of the sections.
type Timer struct {
	listOfStages []string
	stageTimes   map[string]time.Duration
	startTimes   map[string]time.Time
}

const (

	// TotalTimeElapsed is a constant string value to denote total time elapsed.
	TotalTimeElapsed = "Total"
)

var timer = NewTimer()

// StartTimer initiailzes the timer object with the current timestamp information.
func StartTimer(key string) {
	timer.StartTimer(key)
}

// StopTimer records the duration for the current stage sent as the key parameter and stores the information.
func StopTimer(key string) {
	timer.StopTimer(key)
}

// LogSummary prints the summary of all the times collected so far into the INFO section.
func LogSummary() {
	timer.LogSummary(logrus.StandardLogger())
}

// NewTimer returns a new timer that can be used to track sections and
func NewTimer() Timer {
	return Timer{
		listOfStages: []string{},
		stageTimes:   make(map[string]time.Duration),
		startTimes:   make(map[string]time.Time),
	}
}

// StartTimer initializes the timer object with the current timestamp information.
func (t *Timer) StartTimer(key string) {
	t.listOfStages = append(t.listOfStages, key)
	t.startTimes[key] = time.Now().Round(time.Second)
}

// StopTimer records the duration for the current stage sent as the key parameter and stores the information.
func (t *Timer) StopTimer(key string) time.Duration {
	if item, found := t.startTimes[key]; found {
		duration := time.Since(item).Round(time.Second)
		t.stageTimes[key] = duration
	}
	return time.Since(time.Now())
}

// LogSummary prints the summary of all the times collected so far into the INFO section.
// The format of printing will be the following:
// If there are no stages except the total time stage, then it only prints the following
// Time elapsed: <x>m<yy>s
// If there are multiple stages, it prints the following:
// Time elapsed for each section
// Stage1: <x>m<yy>s
// Stage2: <x>m<yy>s
// .
// .
// .
// StageN: <x>m<yy>s
// Time elapsed: <x>m<yy>s
// All durations printed are rounded up to the next second value and printed in the format mentioned above.
func (t *Timer) LogSummary(logger *logrus.Logger) {
	maxLen := 0
	count := 0
	for _, item := range t.listOfStages {
		if len(item) > maxLen && item != TotalTimeElapsed {
			maxLen = len(item)
		}
		if t.stageTimes[item] > 0 {
			count++
		}
	}

	if maxLen != 0 && count > 0 {
		logger.Debugf("Time elapsed per stage:")
	}

	for _, item := range t.listOfStages {
		if item != TotalTimeElapsed && t.stageTimes[item] > 0 {
			logger.Debugf(fmt.Sprintf("%*s: %s", maxLen, item, t.stageTimes[item]))
		}
	}
	logger.Infof("Time elapsed: %s", t.stageTimes[TotalTimeElapsed])
}
