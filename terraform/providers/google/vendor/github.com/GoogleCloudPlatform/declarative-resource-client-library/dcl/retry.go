// Copyright 2023 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dcl

import (
	"context"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
	glog "github.com/golang/glog"
)

// Stop is a value that indicates that no more retries should be attempted.
const Stop time.Duration = -1

// BackoffInitialInterval is the default InitialInterval value for Backoff.
const BackoffInitialInterval = 500 * time.Millisecond

// BackoffMaxInterval is the default MaxInterval value for Backoff.
const BackoffMaxInterval = 30 * time.Second

// RetryDetails provides information about an operation that a Retry implementation
// can use to make decisions about when or if to perform further requests.
type RetryDetails struct {
	Request  *http.Request
	Response *http.Response
}

// Operation is a retryable function. Implementations should return nil to indicate
// that the operation has concluded successfully, OperationNotDone to indicate
// that the operation should be retried, and any other error to indicate that a
// non-retryable error has occurred.
type Operation func(ctx context.Context) (*RetryDetails, error)

// Retry provides an interface for handling retryable operations in a flexible manner.
type Retry interface {
	// RetryAfter returns the amount of time that should elapse before an operation is re-run. Returning
	// Stop indicates that no more retries should occur, and returning zero indicates that the operation
	// should be immediately retried.
	RetryAfter(details *RetryDetails) time.Duration
}

// RetryProvider allows callers to provide custom retry behavior.
type RetryProvider interface {
	// New returns an initialized Retry.
	New() Retry
}

// NoRetry is a Retry implementation that will never retry.
type NoRetry struct{}

// RetryAfter implementation that never retries.
func (n *NoRetry) RetryAfter(_ *RetryDetails) time.Duration {
	return Stop
}

// Reset is a no-op.
func (n *NoRetry) Reset() {}

// Backoff is a Retry implementation that uses exponential backoff with jitter.
type Backoff struct {
	// InitialInterval sets the time interval for the first retry delay.
	InitialInterval time.Duration
	// MaxInterval is the largest amount of time that should elapse between retries.
	MaxInterval time.Duration

	bo *backoff.ExponentialBackOff
}

// NewBackoff returns a Backoff with sensible defaults set.
func NewBackoff() *Backoff {
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = BackoffMaxInterval
	bo.InitialInterval = BackoffInitialInterval
	bo.MaxElapsedTime = 0
	return &Backoff{
		bo: bo,
	}
}

// NewBackoffWithOptions returns a Backoff with caller-supplied parameters.
func NewBackoffWithOptions(initialInterval, maxInterval time.Duration) *Backoff {
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = maxInterval
	bo.InitialInterval = initialInterval
	bo.MaxElapsedTime = 0
	return &Backoff{
		bo: bo,
	}
}

// RetryAfter implementation that uses exponential backoff.
func (n *Backoff) RetryAfter(_ *RetryDetails) time.Duration {
	if next := n.bo.NextBackOff(); next != backoff.Stop {
		return next
	}
	return Stop
}

// BackoffRetryProvider is a default RetryProvider that returns a Backoff.
type BackoffRetryProvider struct{}

// New returns an initialized Retry.
func (r *BackoffRetryProvider) New() Retry {
	return NewBackoff()
}

// Do performs op as a retryable operation, using retry to determine when and if to retry.
// Do will only continue if a OperationNotDone{} is returned. If op() returns another error
// or no error, Do will finish.
// OperationNotDone{} may have an error inside of it, indicating that it's a retryable error.
func Do(ctx context.Context, op Operation, retryProvider RetryProvider) error {
	retry := retryProvider.New()
	for {
		details, err := op(ctx)
		// Responsible for returning nil error too.
		if _, ok := err.(OperationNotDone); !ok {
			return err
		}

		w := retry.RetryAfter(details)
		if w == Stop {
			if e, ok := err.(OperationNotDone); ok {
				if e.Err != nil {
					return e.Err
				}
			}
			return OperationNotDone{}
		}

		t := time.NewTimer(w)
		select {
		case <-ctx.Done():
			t.Stop()
			glog.Info("retryable operation canceled by context")
			return OperationNotDone{}
		case <-t.C:
		}
	}
}
