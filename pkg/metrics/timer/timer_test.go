package timer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func convertToFormat(buf bytes.Buffer) string {
	output := buf.String()

	outputText := ""

	for _, item := range strings.Split(output, "\n") {
		var dat map[string]interface{}

		if err := json.Unmarshal([]byte(item), &dat); err != nil {
			continue
		}
		outputText += dat["msg"].(string) + "\n"
	}

	return outputText
}

func TestBasicLogSummaryMultipleStages(t *testing.T) {
	timer := NewTimer()

	timer.StartTimer(TotalTimeElapsed)
	timer.StartTimer("testStage1")
	timer.StartTimer("testStage2")
	timer.StartTimer("testStage3")
	timer.StartTimer("testStage4")

	time.Sleep(5 * time.Second)

	timer.StopTimer("testStage1")
	timer.StopTimer("testStage2")
	timer.StopTimer("testStage3")
	timer.StopTimer("testStage4")
	timer.StopTimer(TotalTimeElapsed)

	timeElapsed := fmt.Sprintf("Time elapsed per stage:\n")
	time1 := fmt.Sprintf("testStage1: %s\n", timer.stageTimes["testStage1"])
	time2 := fmt.Sprintf("testStage2: %s\n", timer.stageTimes["testStage2"])
	time3 := fmt.Sprintf("testStage3: %s\n", timer.stageTimes["testStage3"])
	time4 := fmt.Sprintf("testStage4: %s\n", timer.stageTimes["testStage4"])
	timeStageElapsed := fmt.Sprintf("Time elapsed: %s\n", timer.stageTimes[TotalTimeElapsed])

	text := timeElapsed + time1 + time2 + time3 + time4 + timeStageElapsed

	textOutput := bytes.Buffer{}

	logger := logrus.New()
	logger.Out = &textOutput
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.JSONFormatter{}

	timer.LogSummary(logger)

	outputText := convertToFormat(textOutput)

	if text != outputText {
		t.Fatalf("expected message summary printed to be %s, but got %s", text, outputText)
	}
}

func TestTotalOnlyLogSummary(t *testing.T) {
	timer := NewTimer()

	timer.StartTimer(TotalTimeElapsed)
	time.Sleep(5 * time.Second)
	timer.StopTimer(TotalTimeElapsed)

	timeStageElapsed := fmt.Sprintf("Time elapsed: %s\n", timer.stageTimes[TotalTimeElapsed])

	textOutput := bytes.Buffer{}

	logger := logrus.New()
	logger.Out = &textOutput
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.JSONFormatter{}

	timer.LogSummary(logger)

	outputText := convertToFormat(textOutput)

	if timeStageElapsed != outputText {
		t.Fatalf("expected message summary printed to be %s, but got %s", timeStageElapsed, outputText)
	}
}

func TestStartAndStopTimer(t *testing.T) {
	timerTotal := NewTimer()

	timerTotal.StartTimer(TotalTimeElapsed)
	time.Sleep(5 * time.Second)

	duration := time.Since(timerTotal.startTimes[TotalTimeElapsed]).Round(time.Second)
	t.Logf("%s", duration)
	if duration < 5*time.Second {
		t.Fatalf("Slept for 5 seconds, expected start time to be 5 seconds old, got %s", duration)
	} else if duration > 10*time.Second {
		t.Fatalf("Slept for 5 seconds, expected start time to be close to 5 seconds old, got %s", duration)
	}
	timerTotal.StopTimer(TotalTimeElapsed)

	if timerTotal.stageTimes[TotalTimeElapsed] < 5*time.Second || timerTotal.stageTimes[TotalTimeElapsed] > 10*time.Second {
		t.Fatalf("Slept for 5 seconds, expected duration to be close to 5 seconds old, got %s", timerTotal.stageTimes[TotalTimeElapsed])
	}
}

func TestNewTimer(t *testing.T) {
	timer := NewTimer()
	if len(timer.listOfStages) != 0 {
		t.Fatalf("Expected empty list of stages property in the new timer created, got %d", len(timer.listOfStages))
	}

	if len(timer.startTimes) != 0 {
		t.Fatalf("Expected empty list of startTimes property in the new timer created, got %d", len(timer.startTimes))
	}

	if len(timer.stageTimes) != 0 {
		t.Fatalf("Expected empty list of startTimes property in the new timer created, got %d", len(timer.stageTimes))
	}
}
