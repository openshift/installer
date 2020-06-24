package ovirt

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeoutWithContext(t *testing.T) {
	cases := []struct {
		name              string
		timeout           time.Duration
		contextTimeout    time.Duration
		contextHasTimeout bool
		expected          time.Duration
		allowForDelta     bool
	}{
		{
			name:     "no deadline",
			timeout:  1 * time.Minute,
			expected: 1 * time.Minute,
		},
		{
			name:              "timeout sooner",
			timeout:           1 * time.Minute,
			contextTimeout:    2 * time.Minute,
			contextHasTimeout: true,
			expected:          1 * time.Minute,
		},
		{
			name:              "timeout later",
			timeout:           2 * time.Minute,
			contextTimeout:    1 * time.Minute,
			contextHasTimeout: true,
			expected:          1 * time.Minute,
			allowForDelta:     true,
		},
		{
			name:              "context done",
			timeout:           1 * time.Minute,
			contextHasTimeout: true,
			expected:          0,
			allowForDelta:     true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			startTime := time.Now()
			ctx := context.Background()
			if tc.contextHasTimeout {
				ctxWithTimeout, cancel := context.WithTimeout(ctx, tc.contextTimeout)
				defer cancel()
				ctx = ctxWithTimeout
			}
			actual := timeoutWithContext(ctx, tc.timeout)
			endTime := time.Now()
			if tc.allowForDelta {
				delta := endTime.Sub(startTime)
				assert.InDelta(t, tc.expected, actual, float64(delta))
			}
		})
	}
}
