package clusterapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/compute/v1"
)

func TestWaitForOperation(t *testing.T) {
	cases := []struct {
		name string
		// ops are returned in order, one per call to the wait function.
		ops           []*compute.Operation
		waitErr       error
		expectedCalls int
		expectedError string
	}{
		{
			name:          "operation already done",
			ops:           []*compute.Operation{{Status: operationStatusDone}},
			expectedCalls: 1,
		},
		{
			// The compute Wait methods return after ~2 minutes whether or not
			// the operation finished, so a RUNNING result must be re-polled
			// rather than treated as success.
			name: "operation still running is polled until done",
			ops: []*compute.Operation{
				{Status: "RUNNING"},
				{Status: "RUNNING"},
				{Status: operationStatusDone},
			},
			expectedCalls: 3,
		},
		{
			name:          "wait error is surfaced",
			waitErr:       fmt.Errorf("connection reset"),
			expectedCalls: 1,
			expectedError: "failed to wait for operation test-op: connection reset",
		},
		{
			// A DONE operation can still carry a failure.
			name: "done operation carrying an error fails",
			ops: []*compute.Operation{{
				Status: operationStatusDone,
				Error: &compute.OperationError{
					Errors: []*compute.OperationErrorErrors{
						{Code: "RESOURCE_NOT_FOUND", Message: "the source object does not exist"},
					},
				},
			}},
			expectedCalls: 1,
			expectedError: "operation test-op failed: RESOURCE_NOT_FOUND: the source object does not exist",
		},
		{
			name:          "done operation with an empty error is not a failure",
			ops:           []*compute.Operation{{Status: operationStatusDone, Error: &compute.OperationError{}}},
			expectedCalls: 1,
		},
	}

	// Keep the re-poll cases from sleeping for the production interval.
	original := operationPollInterval
	operationPollInterval = time.Millisecond
	defer func() { operationPollInterval = original }()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			wait := func() (*compute.Operation, error) {
				calls++
				if tc.waitErr != nil {
					return nil, tc.waitErr
				}
				return tc.ops[calls-1], nil
			}

			err := waitForOperation(t.Context(), "test-op", wait)
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
			assert.Equal(t, tc.expectedCalls, calls)
		})
	}
}

func TestWaitForOperationHonorsContextDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// An operation that never reaches DONE must be bounded by the caller's
	// context rather than polling forever.
	err := waitForOperation(ctx, "stuck-op", func() (*compute.Operation, error) {
		return &compute.Operation{Status: "RUNNING"}, nil
	})

	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.ErrorContains(t, err, "timed out waiting for operation stuck-op to complete")
}
