package agent

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

// ValidationHistory test helpers
const (
	testKey                string = "test-cluster"
	testID                 string = "test"
	testLogPrefix          string = "Test: "
	validationsInfoSuccess string = "{\"test\":[{\"id\":\"test\",\"status\":\"success\",\"message\":\"The validation succeeded\"}]}"
	validationsInfoFailure string = "{\"test\":[{\"id\":\"test\",\"status\":\"failure\",\"message\":\"The validation failed\"}]}"
)

func TestUpdateValidationHistory(t *testing.T) {
	tests := []struct {
		name            string
		validationsInfo string
		inputHistory    map[string]*validationResultHistory
		expectedResult  map[string]*validationResultHistory
	}{
		{
			name:            "success-first-seen",
			validationsInfo: validationsInfoSuccess,
			inputHistory: map[string]*validationResultHistory{
				"test": {
					numFailures:     0,
					seen:            false,
					currentStatus:   "",
					currentMessage:  "",
					previousStatus:  "",
					previousMessage: "",
				},
			},
			expectedResult: map[string]*validationResultHistory{
				"test": {
					numFailures:     0,
					seen:            true,
					currentStatus:   "success",
					currentMessage:  "The validation succeeded",
					previousStatus:  "",
					previousMessage: "",
				},
			},
		},
		{
			name:            "success-no-change",
			validationsInfo: validationsInfoSuccess,
			inputHistory: map[string]*validationResultHistory{
				"test": {
					numFailures:     0,
					seen:            true,
					currentStatus:   "success",
					currentMessage:  "The validation succeeded",
					previousStatus:  "",
					previousMessage: "",
				},
			},
			expectedResult: map[string]*validationResultHistory{
				"test": {
					numFailures:     0,
					seen:            true,
					currentStatus:   "success",
					currentMessage:  "The validation succeeded",
					previousStatus:  "success",
					previousMessage: "The validation succeeded",
				},
			},
		},
		{
			name:            "failure-first-seen",
			validationsInfo: validationsInfoFailure,
			inputHistory: map[string]*validationResultHistory{
				"test": {
					numFailures:     0,
					seen:            false,
					currentStatus:   "",
					currentMessage:  "",
					previousStatus:  "",
					previousMessage: "",
				},
			},
			expectedResult: map[string]*validationResultHistory{
				"test": {
					numFailures:     1,
					seen:            true,
					currentStatus:   "failure",
					currentMessage:  "The validation failed",
					previousStatus:  "",
					previousMessage: "",
				},
			},
		},
		{
			name:            "failure-no-change",
			validationsInfo: validationsInfoFailure,
			inputHistory: map[string]*validationResultHistory{
				"test": {
					numFailures:     0,
					seen:            true,
					currentStatus:   "failure",
					currentMessage:  "The validation failed",
					previousStatus:  "",
					previousMessage: "",
				},
			},
			expectedResult: map[string]*validationResultHistory{
				"test": {
					numFailures:     1,
					seen:            true,
					currentStatus:   "failure",
					currentMessage:  "The validation failed",
					previousStatus:  "failure",
					previousMessage: "The validation failed",
				},
			},
		},
		{
			name:            "success-after-failure",
			validationsInfo: validationsInfoSuccess,
			inputHistory: map[string]*validationResultHistory{
				"test": {
					numFailures:     1,
					seen:            true,
					currentStatus:   "failure",
					currentMessage:  "The validation failed",
					previousStatus:  "failure",
					previousMessage: "The validation failed",
				},
			},
			expectedResult: map[string]*validationResultHistory{
				"test": {
					numFailures:     1,
					seen:            true,
					currentStatus:   "success",
					currentMessage:  "The validation succeeded",
					previousStatus:  "failure",
					previousMessage: "The validation failed",
				},
			},
		},
		{
			name:            "failure-after-success",
			validationsInfo: validationsInfoFailure,
			inputHistory: map[string]*validationResultHistory{
				"test": {
					numFailures:     0,
					seen:            true,
					currentStatus:   "success",
					currentMessage:  "The validation succeeded",
					previousStatus:  "success",
					previousMessage: "The validation succeeded",
				},
			},
			expectedResult: map[string]*validationResultHistory{
				"test": {
					numFailures:     1,
					seen:            true,
					currentStatus:   "failure",
					currentMessage:  "The validation failed",
					previousStatus:  "success",
					previousMessage: "The validation succeeded",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test logger
			logger := &logrus.Logger{
				Out:       io.Discard,
				Formatter: new(logrus.TextFormatter),
				Hooks:     make(logrus.LevelHooks),
				Level:     logrus.DebugLevel,
			}
			actualResult, err := updateValidationResultHistory(testLogPrefix, tt.validationsInfo, tt.inputHistory, logger, nil)
			assert.Equal(t, tt.expectedResult, actualResult)
			assert.Nil(t, err)
		})
	}
}

func TestLogValidationHistory(t *testing.T) {
	tests := []struct {
		name           string
		inputHistory   *validationResultHistory
		expectedResult *logrus.Entry
	}{
		{
			name: "default-case-not-seen",
			inputHistory: &validationResultHistory{
				numFailures:     0,
				seen:            false,
				currentStatus:   "test",
				currentMessage:  "This is the default validation case",
				previousStatus:  "",
				previousMessage: "",
			},
			expectedResult: &logrus.Entry{Level: logrus.TraceLevel, Message: "Test: This is the default validation case"},
		},
		{
			name: "default-case-seen",
			inputHistory: &validationResultHistory{
				numFailures:     0,
				seen:            true,
				currentStatus:   "test",
				currentMessage:  "This is the default validation case",
				previousStatus:  "",
				previousMessage: "",
			},
			expectedResult: &logrus.Entry{Level: logrus.TraceLevel, Message: "Test: This is the default validation case"},
		},
		{
			name: "success-first-seen",
			inputHistory: &validationResultHistory{
				numFailures:     0,
				seen:            false,
				currentStatus:   "success",
				currentMessage:  "The validation succeeded",
				previousStatus:  "",
				previousMessage: "",
			},
			expectedResult: &logrus.Entry{Level: logrus.DebugLevel, Message: "Test: The validation succeeded"},
		},
		{
			name: "success-no-change",
			inputHistory: &validationResultHistory{
				numFailures:     0,
				seen:            true,
				currentStatus:   "success",
				currentMessage:  "The validation succeeded",
				previousStatus:  "success",
				previousMessage: "The validation succeeded",
			},
			expectedResult: nil,
		},
		{
			name: "failure-first-seen",
			inputHistory: &validationResultHistory{
				numFailures:     0,
				seen:            false,
				currentStatus:   "failure",
				currentMessage:  "The validation failed",
				previousStatus:  "",
				previousMessage: "",
			},
			expectedResult: &logrus.Entry{Level: logrus.WarnLevel, Message: "Test: The validation failed"},
		},
		{
			name: "failure-no-change",
			inputHistory: &validationResultHistory{
				numFailures:     1,
				seen:            true,
				currentStatus:   "failure",
				currentMessage:  "The validation failed",
				previousStatus:  "failure",
				previousMessage: "The validation failed",
			},
			expectedResult: nil,
		},
		{
			name: "success-after-failure",
			inputHistory: &validationResultHistory{
				numFailures:     1,
				seen:            true,
				currentStatus:   "success",
				currentMessage:  "The validation succeeded",
				previousStatus:  "failure",
				previousMessage: "The validation failed",
			},
			expectedResult: &logrus.Entry{Level: logrus.InfoLevel, Message: "Test: The validation succeeded"},
		},
		{
			name: "failure-after-success",
			inputHistory: &validationResultHistory{
				numFailures:     1,
				seen:            true,
				currentStatus:   "failure",
				currentMessage:  "The validation failed",
				previousStatus:  "success",
				previousMessage: "The validation succeeded",
			},
			expectedResult: &logrus.Entry{Level: logrus.WarnLevel, Message: "Test: The validation failed"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test logger
			logger := &logrus.Logger{
				Out:       io.Discard,
				Formatter: new(logrus.TextFormatter),
				Hooks:     make(logrus.LevelHooks),
				Level:     logrus.TraceLevel,
			}
			var hook = test.NewLocal(logger)

			logValidationHistory(testLogPrefix, tt.inputHistory, logger, nil)
			actualResult := hook.LastEntry()

			if actualResult == nil {
				assert.Equal(t, tt.expectedResult, actualResult)
			} else {
				assert.Equal(t, tt.expectedResult.Level, actualResult.Level)
				assert.Equal(t, tt.expectedResult.Message, actualResult.Message)
			}
		})
	}
}
