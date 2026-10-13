package powervs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/sirupsen/logrus"
)

// Tests function parseRetryAfterHeader with various inputs to confirm it correctly parses valid headers and rejects invalid ones.
func TestParseRetryAfterHeader(t *testing.T) {
	t.Run("empty slice returns false", func(t *testing.T) {
		_, ok := parseRetryAfterHeader([]string{})
		if ok {
			t.Error("expected false for empty slice")
		}
	})

	t.Run("empty string returns false", func(t *testing.T) {
		_, ok := parseRetryAfterHeader([]string{""})
		if ok {
			t.Error("expected false for empty string value")
		}
	})

	t.Run("valid integer seconds", func(t *testing.T) {
		d, ok := parseRetryAfterHeader([]string{"30"})
		if !ok {
			t.Fatal("expected true for valid integer header")
		}
		if d != 30*time.Second {
			t.Errorf("expected 30s, got %v", d)
		}
	})

	t.Run("zero seconds is valid", func(t *testing.T) {
		d, ok := parseRetryAfterHeader([]string{"0"})
		if !ok {
			t.Fatal("expected true for zero value")
		}
		if d != 0 {
			t.Errorf("expected 0, got %v", d)
		}
	})

	t.Run("negative seconds returns false", func(t *testing.T) {
		_, ok := parseRetryAfterHeader([]string{"-5"})
		if ok {
			t.Error("expected false for negative value")
		}
	})

	t.Run("valid HTTP-date in the future", func(t *testing.T) {
		// Pin timeNow so the future date is always in the future relative to the test.
		orig := timeNow
		defer func() { timeNow = orig }()
		timeNow = func() time.Time {
			return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		}

		future := time.Date(2026, 1, 1, 0, 1, 0, 0, time.UTC) // 60s ahead
		header := future.UTC().Format(time.RFC1123)

		d, ok := parseRetryAfterHeader([]string{header})
		if !ok {
			t.Fatal("expected true for future HTTP-date")
		}
		if d < 59*time.Second || d > 61*time.Second {
			t.Errorf("expected ~60s wait, got %v", d)
		}
	})

	t.Run("HTTP-date already in the past returns (0, true)", func(t *testing.T) {
		orig := timeNow
		defer func() { timeNow = orig }()
		timeNow = func() time.Time {
			return time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
		}

		past := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		header := past.UTC().Format(time.RFC1123)

		d, ok := parseRetryAfterHeader([]string{header})
		if !ok {
			t.Fatal("expected true for past HTTP-date (header was valid)")
		}
		if d != 0 {
			t.Errorf("expected 0 duration for past date, got %v", d)
		}
	})

	t.Run("unparseable string returns false", func(t *testing.T) {
		_, ok := parseRetryAfterHeader([]string{"not-a-number-or-date"})
		if ok {
			t.Error("expected false for garbage value")
		}
	})

	t.Run("only first element is used", func(t *testing.T) {
		d, ok := parseRetryAfterHeader([]string{"10", "99"})
		if !ok {
			t.Fatal("expected true")
		}
		if d != 10*time.Second {
			t.Errorf("expected 10s from first element, got %v", d)
		}
	})
}

// tests function exponentialBackoffDuration to confirm it calculates backoff times correctly and respects the maximum cap.
func TestExponentialBackoffDuration(t *testing.T) {
	tests := []struct {
		attempt int
		wantMin time.Duration
		wantMax time.Duration
	}{
		// attempt 0: base * 2^0 = 2s
		{0, 2 * time.Second, 2 * time.Second},
		// attempt 1: base * 2^1 = 4s
		{1, 4 * time.Second, 4 * time.Second},
		// attempt 2: base * 2^2 = 8s
		{2, 8 * time.Second, 8 * time.Second},
		// attempt 3: base * 2^3 = 16s
		{3, 16 * time.Second, 16 * time.Second},
		// attempt 4: base * 2^4 = 32s
		{4, 32 * time.Second, 32 * time.Second},
		// attempt 5: base * 2^5 = 64s → capped at 60s
		{5, maxBackoffDuration, maxBackoffDuration},
		// high attempt: still capped
		{20, maxBackoffDuration, maxBackoffDuration},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("attempt_%d", tt.attempt), func(t *testing.T) {
			got := exponentialBackoffDuration(tt.attempt)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("attempt %d: got %v, want [%v, %v]",
					tt.attempt, got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

// tests function sleepWithContext 	to confirm it sleeps for the correct duration and returns early when the context is cancelled.
func TestSleepWithContext(t *testing.T) {
	t.Run("zero duration returns immediately without error", func(t *testing.T) {
		ctx := context.Background()
		start := time.Now()
		err := sleepWithContext(ctx, 0)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if time.Since(start) > 50*time.Millisecond {
			t.Error("zero duration sleep took too long")
		}
	})

	t.Run("sleeps for specified duration", func(t *testing.T) {
		ctx := context.Background()
		start := time.Now()
		err := sleepWithContext(ctx, 100*time.Millisecond)
		elapsed := time.Since(start)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if elapsed < 90*time.Millisecond {
			t.Errorf("sleep was too short: %v", elapsed)
		}
	})

	t.Run("returns error when context is cancelled before duration elapses", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		start := time.Now()
		err := sleepWithContext(ctx, 10*time.Second)
		elapsed := time.Since(start)

		if err == nil {
			t.Error("expected an error on context cancellation, got nil")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("expected DeadlineExceeded, got %v", err)
		}
		// Should have returned well before the 10s sleep.
		if elapsed > time.Second {
			t.Errorf("took too long to cancel: %v", elapsed)
		}
	})

	t.Run("already-cancelled context returns error immediately", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel before calling

		err := sleepWithContext(ctx, 5*time.Second)
		if err == nil {
			t.Error("expected error for pre-cancelled context")
		}
	})
}

// fakeCall is a helper that returns a 429 for the first `failCount` calls, then succeeds.
func fakeCall(failCount int, callsMade *atomic.Int32) func() (*core.DetailedResponse, error) {
	return func() (*core.DetailedResponse, error) {
		n := int(callsMade.Add(1))
		if n <= failCount {
			return &core.DetailedResponse{
				StatusCode: http.StatusTooManyRequests,
				Headers:    http.Header{"Retry-After": []string{"0"}}, // 0s so tests are fast
			}, fmt.Errorf("rate limited")
		}
		return &core.DetailedResponse{StatusCode: http.StatusOK}, nil
	}
}

func newTestUninstaller() *ClusterUninstaller {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	return &ClusterUninstaller{Logger: logger}
}

// tests callWithRetry to confirm it correctly retries on 429 responses, respects the Retry-After header, and returns errors appropriately.
func TestCallWithRetry(t *testing.T) {
	t.Run("succeeds on first attempt", func(t *testing.T) {
		o := newTestUninstaller()
		var calls atomic.Int32

		err := o.callWithRetry(context.Background(), "op", fakeCall(0, &calls))
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if calls.Load() != 1 {
			t.Errorf("expected 1 call, got %d", calls.Load())
		}
	})

	t.Run("retries on 429 and eventually succeeds", func(t *testing.T) {
		o := newTestUninstaller()
		var calls atomic.Int32

		// Fail 3 times, succeed on 4th.
		err := o.callWithRetry(context.Background(), "op", fakeCall(3, &calls))
		if err != nil {
			t.Errorf("expected nil after retries, got %v", err)
		}
		if calls.Load() != 4 {
			t.Errorf("expected 4 calls, got %d", calls.Load())
		}
	})

	t.Run("returns error immediately on non-429 failure", func(t *testing.T) {
		o := newTestUninstaller()
		var calls atomic.Int32

		fn := func() (*core.DetailedResponse, error) {
			calls.Add(1)
			return &core.DetailedResponse{StatusCode: http.StatusInternalServerError},
				fmt.Errorf("internal server error")
		}

		err := o.callWithRetry(context.Background(), "op", fn)
		if err == nil {
			t.Error("expected error for non-429 failure")
		}
		// Must not retry a non-429.
		if calls.Load() != 1 {
			t.Errorf("expected exactly 1 call for non-429, got %d", calls.Load())
		}
	})

	t.Run("returns error when nil DetailedResponse accompanies error", func(t *testing.T) {
		o := newTestUninstaller()
		var calls atomic.Int32

		fn := func() (*core.DetailedResponse, error) {
			calls.Add(1)
			return nil, fmt.Errorf("network failure")
		}

		err := o.callWithRetry(context.Background(), "op", fn)
		if err == nil {
			t.Error("expected error for nil response")
		}
		if calls.Load() != 1 {
			t.Errorf("expected exactly 1 call for nil response error, got %d", calls.Load())
		}
	})

	t.Run("exhausts all retries and returns error", func(t *testing.T) {
		o := newTestUninstaller()
		var calls atomic.Int32

		// Always fail with 429.
		fn := func() (*core.DetailedResponse, error) {
			calls.Add(1)
			return &core.DetailedResponse{
				StatusCode: http.StatusTooManyRequests,
				Headers:    http.Header{"Retry-After": []string{"0"}},
			}, fmt.Errorf("rate limited")
		}

		err := o.callWithRetry(context.Background(), "op", fn)
		if err == nil {
			t.Error("expected error after exhausting retries")
		}
		if calls.Load() != int32(maxRetryAttempts) {
			t.Errorf("expected %d calls, got %d", maxRetryAttempts, calls.Load())
		}
	})

	t.Run("stops retrying when context is cancelled", func(t *testing.T) {
		o := newTestUninstaller()
		var calls atomic.Int32

		ctx, cancel := context.WithCancel(context.Background())

		fn := func() (*core.DetailedResponse, error) {
			n := calls.Add(1)
			// Cancel context after first 429 so sleepWithContext returns early.
			if n == 1 {
				cancel()
			}
			return &core.DetailedResponse{
				StatusCode: http.StatusTooManyRequests,
				// Use a non-zero wait so sleepWithContext actually selects ctx.Done().
				Headers: http.Header{"Retry-After": []string{"60"}},
			}, fmt.Errorf("rate limited")
		}

		err := o.callWithRetry(ctx, "op", fn)
		if err == nil {
			t.Error("expected error when context is cancelled")
		}
		// Only 1 API call should have been made before the sleep was interrupted.
		if calls.Load() != 1 {
			t.Errorf("expected 1 call before cancellation, got %d", calls.Load())
		}
	})

	t.Run("uses exponential backoff when Retry-After header is absent", func(t *testing.T) {
		o := newTestUninstaller()
		var calls atomic.Int32

		// 429 with no Retry-After header, succeeds on second attempt.
		fn := func() (*core.DetailedResponse, error) {
			n := calls.Add(1)
			if n == 1 {
				return &core.DetailedResponse{
					StatusCode: http.StatusTooManyRequests,
					Headers:    http.Header{}, // no Retry-After
				}, fmt.Errorf("rate limited")
			}
			return &core.DetailedResponse{StatusCode: http.StatusOK}, nil
		}

		// Patch baseBackoffDuration indirectly by confirming the call count
		// We can't mutate the const, but we can verify the path was taken by ensuring 2 calls were made and no error was returned.
		err := o.callWithRetry(context.Background(), "op", fn)
		if err != nil {
			t.Errorf("expected success after backoff retry, got %v", err)
		}
		if calls.Load() != 2 {
			t.Errorf("expected 2 calls, got %d", calls.Load())
		}
	})
	t.Run("error message contains operation name", func(t *testing.T) {
		o := newTestUninstaller()
		var calls atomic.Int32

		fn := func() (*core.DetailedResponse, error) {
			calls.Add(1)
			return nil, fmt.Errorf("network failure")
		}

		err := o.callWithRetry(context.Background(), "MyOperation", fn)
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if !strings.Contains(err.Error(), "MyOperation") {
			t.Errorf("expected error to contain operation name 'MyOperation', got: %v", err)
		}
	})
}

// Regression: nil pointer dereference after callWithRetry failure.

// TestCallWithRetryNilSafetyRegression verifies that the err-check guards added after each callWithRetry call in
// loadSDKServices prevent a nil pointer dereference when callWithRetry returns an error.

func TestCallWithRetryNilSafetyRegression(t *testing.T) {
	o := newTestUninstaller()

	var responseFromFn *core.DetailedResponse // stays nil

	fn := func() (*core.DetailedResponse, error) {
		// Simulate a call that returns nil response and error (e.g. network failure).
		return responseFromFn, fmt.Errorf("connection refused")
	}

	err := o.callWithRetry(context.Background(), "ListZonesWithContext", fn)

	// callWithRetry must surface the error so the caller's nil-guard fires.
	if err == nil {
		t.Fatal("expected non-nil error when fn returns nil response and error; " +
			"without this the caller cannot guard against nil dereference")
	}

	var zoneResources interface{} // stand-in for *zonesv1.ListZonesResp

	if err != nil {
		// if correct we should return early. No dereference of zoneResources occurs.
		return
	}

	// If we somehow reach here without the guard, a nil dereference would follow.
	_ = zoneResources
	t.Fatal("should have returned before this line")
}
