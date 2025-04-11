package workflowreport

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	stageEmpty  = NewStageID("", "")
	stageS1     = NewStageID("s1", "s1 long desc")
	stageS1Sub1 = NewStageID("s1.sub1", "")
	stageS1Sub2 = NewStageID("s1.sub2", "s1.sub2 long desc")
	stageS2     = NewStageID("s2", "")
)

func TestInvalidCases(t *testing.T) {
	cases := []struct {
		name          string
		do            func(r Report) error
		expectedError string
	}{
		{
			name: "result - empty substage",
			do: func(r Report) error {
				return r.SubStageResult(stageEmpty, "some value")
			},
			expectedError: "cannot add substage result: empty stage name",
		},
		{
			name: "result - substage stage not found",
			do: func(r Report) error {
				return r.SubStageResult(stageS1Sub2, "some value")
			},
			expectedError: "cannot add substage result: s1.sub2 stage owner not found",
		},
		{
			name: "result - substage not found",
			do: func(r Report) error {
				assert.NoError(t, r.Stage(stageS1))
				return r.SubStageResult(stageS1Sub2, "some value")
			},
			expectedError: "cannot add result for substage: s1.sub2 substage not found",
		},
		{
			name: "result - stage id not valid",
			do: func(r Report) error {
				return r.StageResult(stageEmpty, "some value")
			},
			expectedError: "cannot add result for stage : empty stage name",
		},
		{
			name: "result - stage not found",
			do: func(r Report) error {
				return r.StageResult(stageS1, "some value")
			},
			expectedError: "cannot add result for stage: s1 not found",
		},
		{
			name: "sub - already active",
			do: func(r Report) error {
				assert.NoError(t, r.Stage(stageS1))
				assert.NoError(t, r.SubStage(stageS1Sub1))
				return r.SubStage(stageS1Sub1)
			},
			expectedError: "cannot add substage: substage s1.sub1 already active",
		},
		{
			name: "sub - stage not active",
			do: func(r Report) error {
				assert.NoError(t, r.Stage(stageS1))
				assert.NoError(t, r.Stage(stageS2))
				return r.SubStage(stageS1Sub1)
			},
			expectedError: "cannot add substage: s1.sub1 stage owner is not active",
		},
		{
			name: "sub - owner not found",
			do: func(r Report) error {
				assert.NoError(t, r.Stage(stageS2))
				return r.SubStage(stageS1Sub1)
			},
			expectedError: "cannot add substage: s1.sub1 stage owner not found",
		},
		{
			name: "sub - empty",
			do: func(r Report) error {
				assert.NoError(t, r.Stage(stageS1))
				return r.SubStage(stageEmpty)
			},
			expectedError: "cannot add substage: empty stage name",
		},
		{
			name: "sub - empty stage",
			do: func(r Report) error {
				assert.NoError(t, r.Stage(stageS1))
				return r.SubStage(NewStageID(".sub1", ""))
			},
			expectedError: "cannot add substage: empty stage name",
		},
		{
			name: "sub - empty sub id",
			do: func(r Report) error {
				assert.NoError(t, r.Stage(stageS1))
				return r.SubStage(NewStageID("s1.", ""))
			},
			expectedError: "cannot add substage: empty substage name: s1.",
		},
		{
			name: "empty stage",
			do: func(r Report) error {
				return r.Stage(stageEmpty)
			},
			expectedError: "cannot add stage: empty stage name",
		},
		{
			name: "sub stage",
			do: func(r Report) error {
				return r.Stage(stageS1Sub1)
			},
			expectedError: "cannot add stage: invalid stage name: s1.sub1",
		},
		{
			name: "active stage",
			do: func(r Report) error {
				assert.NoError(t, r.Stage(stageS1))
				return r.Stage(stageS1)
			},
			expectedError: "cannot add stage: s1 is already active",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var fakeTime fakeTimeProvider
			var fakeWriter fakeWriter

			report := newReport("invalid-test", &fakeWriter, &fakeTime)
			err := tc.do(report)
			assert.Equal(t, tc.expectedError, err.Error())
		})
	}
}

func TestGenerate(t *testing.T) {
	cases := []struct {
		name           string
		do             func(t *testing.T, r Report)
		verify         func(t *testing.T, r report)
		expectedReport string
	}{
		{
			name: "empty",
			verify: func(t *testing.T, r report) {
				t.Helper()
				assert.Len(t, r.Stages, 0)
				assert.True(t, r.StartTime.IsZero())
				assert.True(t, r.EndTime.IsZero())
				assert.Nil(t, r.Result)
				assert.NoError(t, r.Complete(nil))
			},
			expectedReport: `{
  "id": "report-test-202401010000",
  "start_time": "0001-01-01T00:00:00Z",
  "end_time": "0001-01-01T00:00:00Z",
  "result": {
    "exit_code": 0
  }
}`,
		},
		{
			name: "empty - not complete",
			verify: func(t *testing.T, r report) {
				t.Helper()
				assert.Len(t, r.Stages, 0)
				assert.True(t, r.StartTime.IsZero())
				assert.True(t, r.EndTime.IsZero())
				assert.Nil(t, r.Result)
				expectedPartialReport := `{
  "id": "report-test-202401010000",
  "start_time": "0001-01-01T00:00:00Z",
  "end_time": "0001-01-01T00:00:00Z",
  "result": null
}`
				assert.Equal(t, expectedPartialReport, string(r.json()))
			},
			expectedReport: "",
		},
		{
			name: "single stage with result",
			do: func(t *testing.T, r Report) {
				t.Helper()
				assert.NoError(t, r.Stage(stageS1))
				assert.NoError(t, r.Complete(fmt.Errorf("some error: %w", fmt.Errorf("specific error"))))
			},
			verify: func(t *testing.T, r report) {
				t.Helper()
				assert.Len(t, r.Stages, 1)
				assert.Equal(t, r.Stages[0].ID(), "s1")
				assert.Equal(t, r.Stages[0].Description(), "s1 long desc")
				assert.Equal(t, r.Result.ExitCode, 1)
				assert.Equal(t, r.Result.ErrorMessage, "specific error")
				assert.Equal(t, r.Result.DetailedErrorMessage, "some error: specific error")
				// report time equals to the stage time
				assert.Equal(t, r.StartTime, r.Stages[0].StartTime)
				assert.Equal(t, r.EndTime, r.Stages[0].EndTime)
			},
			expectedReport: `{
  "id": "report-test-202401010000",
  "start_time": "2024-01-01T00:00:02Z",
  "end_time": "2024-01-01T00:00:03Z",
  "stages": [
    {
      "id": "s1",
      "description": "s1 long desc",
      "start_time": "2024-01-01T00:00:02Z",
      "end_time": "2024-01-01T00:00:03Z"
    }
  ],
  "result": {
    "exit_code": 1,
    "error_message": "specific error",
    "detailed_error_message": "some error: specific error"
  }
}`,
		},
		{
			name: "two stages",
			do: func(t *testing.T, r Report) {
				t.Helper()
				assert.NoError(t, r.Stage(stageS1))
				assert.NoError(t, r.Stage(stageS2))
				assert.NoError(t, r.Complete(nil))
			},
			verify: func(t *testing.T, r report) {
				t.Helper()
				assert.Len(t, r.Stages, 2)
				assert.Equal(t, r.Stages[0].ID(), "s1")
				assert.Equal(t, r.Stages[1].ID(), "s2")
				// report time equals to all the stages time
				assert.Equal(t, r.StartTime, r.Stages[0].StartTime)
				assert.Equal(t, r.EndTime, r.Stages[1].EndTime)
			},
			expectedReport: `{
  "id": "report-test-202401010000",
  "start_time": "2024-01-01T00:00:02Z",
  "end_time": "2024-01-01T00:00:05Z",
  "stages": [
    {
      "id": "s1",
      "description": "s1 long desc",
      "start_time": "2024-01-01T00:00:02Z",
      "end_time": "2024-01-01T00:00:03Z"
    },
    {
      "id": "s2",
      "start_time": "2024-01-01T00:00:04Z",
      "end_time": "2024-01-01T00:00:05Z"
    }
  ],
  "result": {
    "exit_code": 0
  }
}`,
		},
		{
			name: "stage with substages",
			do: func(t *testing.T, r Report) {
				t.Helper()
				assert.NoError(t, r.Stage(stageS1))
				assert.NoError(t, r.SubStage(stageS1Sub1))
				assert.NoError(t, r.SubStage(stageS1Sub2))
				assert.NoError(t, r.Complete(nil))
			},
			verify: func(t *testing.T, r report) {
				t.Helper()
				s := r.Stages[0]
				assert.Len(t, s.SubStages, 2)
				assert.Equal(t, s.SubStages[0].ID(), "s1.sub1")
				assert.Equal(t, s.SubStages[1].ID(), "s1.sub2")

				assert.True(t, s.StartTime.Before(s.SubStages[0].StartTime))
				assert.True(t, s.EndTime.After(s.SubStages[1].EndTime))
			},
			expectedReport: `{
  "id": "report-test-202401010000",
  "start_time": "2024-01-01T00:00:02Z",
  "end_time": "2024-01-01T00:00:07Z",
  "stages": [
    {
      "id": "s1",
      "description": "s1 long desc",
      "start_time": "2024-01-01T00:00:02Z",
      "end_time": "2024-01-01T00:00:07Z",
      "sub_stages": [
        {
          "id": "s1.sub1",
          "start_time": "2024-01-01T00:00:03Z",
          "end_time": "2024-01-01T00:00:04Z"
        },
        {
          "id": "s1.sub2",
          "description": "s1.sub2 long desc",
          "start_time": "2024-01-01T00:00:05Z",
          "end_time": "2024-01-01T00:00:06Z"
        }
      ]
    }
  ],
  "result": {
    "exit_code": 0
  }
}`,
		},
		{
			name: "stage result",
			do: func(t *testing.T, r Report) {
				t.Helper()
				assert.NoError(t, r.Stage(stageS1))
				assert.NoError(t, r.StageResult(stageS1, "some value"))
				assert.NoError(t, r.SubStage(stageS1Sub1))
				assert.NoError(t, r.SubStageResult(stageS1Sub1, "some sub result"))
				assert.NoError(t, r.Complete(nil))
			},
			expectedReport: `{
  "id": "report-test-202401010000",
  "start_time": "2024-01-01T00:00:02Z",
  "end_time": "2024-01-01T00:00:05Z",
  "stages": [
    {
      "id": "s1",
      "description": "s1 long desc",
      "start_time": "2024-01-01T00:00:02Z",
      "end_time": "2024-01-01T00:00:05Z",
      "result": "some value",
      "sub_stages": [
        {
          "id": "s1.sub1",
          "start_time": "2024-01-01T00:00:03Z",
          "end_time": "2024-01-01T00:00:04Z",
          "result": "some sub result"
        }
      ]
    }
  ],
  "result": {
    "exit_code": 0
  }
}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var fakeTime fakeTimeProvider
			var fakeWriter fakeWriter

			rep := newReport("test", &fakeWriter, &fakeTime)
			if tc.do != nil {
				tc.do(t, rep)
			}

			r, ok := rep.(*report)
			assert.True(t, ok)
			if tc.verify != nil {
				tc.verify(t, *r)
			}
			if tc.expectedReport != "" {
				assert.Equal(t, tc.expectedReport, fakeWriter.buf.String())
			}
		})
	}
}

func TestReportContext(t *testing.T) {
	cases := []struct {
		name           string
		ctxName        string
		expectedReport string
	}{
		{
			name:           "report from context",
			ctxName:        "test",
			expectedReport: `(?s){.*"id": "report-test-\d{12}",.*"start_time":.*"end_time":.*}`,
		},
		{
			name:           "no report",
			expectedReport: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.ctxName != "" {
				ctx = Context("test", "")
			}
			r, ok := GetReport(ctx).(*report)
			assert.Equal(t, ok, tc.expectedReport != "")
			assert.Regexp(t, tc.expectedReport, string(r.json()))
		})
	}
}

type fakeTimeProvider struct {
	secs int
}

func (ftp *fakeTimeProvider) Now() time.Time {
	// Simulate one second tick for each call
	ftp.secs++
	return time.Date(2024, 1, 1, 0, 0, ftp.secs, 0, time.UTC)
}

type fakeWriter struct {
	buf bytes.Buffer
}

func (fw *fakeWriter) Write(data []byte) (n int, err error) {
	fw.buf.Reset()
	return fw.buf.Write(data)
}
