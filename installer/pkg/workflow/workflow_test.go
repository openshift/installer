package workflow

import (
	"errors"
	"testing"
)

func test1Step(m *metadata) error {
	return nil
}

func test2Step(m *metadata) error {
	return nil
}

func test3Step(m *metadata) error {
	return errors.New("Step failed!")
}

func TestWorkflowTypeExecute(t *testing.T) {
	m := metadata{}

	testCases := []struct {
		test          string
		steps         []Step
		m             metadata
		expectedError bool
	}{
		{
			test:          "All steps succeed",
			steps:         []Step{test1Step, test2Step},
			m:             m,
			expectedError: false,
		},
		{
			test:          "At least one step fails",
			steps:         []Step{test1Step, test2Step, test3Step},
			m:             m,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		wf := Workflow{
			metadata: tc.m,
			steps:    tc.steps,
		}
		err := wf.Execute()
		if (err != nil) != tc.expectedError {
			t.Errorf("Test case %s: WorkflowType.Execute() expected error: %v, got: %v", tc.test, tc.expectedError, (err != nil))
		}
	}
}
