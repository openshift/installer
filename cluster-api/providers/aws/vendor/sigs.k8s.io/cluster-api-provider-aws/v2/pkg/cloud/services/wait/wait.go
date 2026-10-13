/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package wait provides a set of utilities for polling and waiting.
package wait

import (
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
)

/*
 Ideally, this entire file would be replaced with returning a retryable
 error and letting the actuator requeue deletion. Unfortunately, since
 the retry behaviour is not tunable, with a max retry limit of 10, we
 implement waits manually here.
*/

// NewBackoff creates a new API Machinery backoff parameter set suitable
// for use with AWS services.
func NewBackoff() wait.Backoff {
	// Return a exponential backoff configuration which
	// returns durations for a total time of ~5m.
	// Example: 1s, 1.7s, 2.9s, 5s, 8.6s, 14.6s, 25s, 42.8s, 73.1s ... 125s
	// Jitter is added as a random fraction of the duration multiplied by the
	// jitter factor.
	return wait.Backoff{
		Duration: time.Second,
		Factor:   1.71,
		Steps:    10,
		Jitter:   0.4,
	}
}

// WaitForWithRetryable repeats a condition check with exponential backoff.
func WaitForWithRetryable(backoff wait.Backoff, condition wait.ConditionFunc, retryableErrors ...string) error {
	var errToReturn error
	waitErr := wait.ExponentialBackoff(backoff, func() (bool, error) {
		// clear errToReturn value from previous iteration
		errToReturn = nil

		ok, err := condition()
		if ok {
			// All done!
			return true, nil
		}
		if err == nil {
			// Not done, but no error, so keep waiting.
			return false, nil
		}

		// If the returned error isn't empty, check if the error is a retryable one,
		// or return immediately.
		// Also check for smithy errors
		var code string
		smithyErr := awserrors.ParseSmithyError(err)
		if smithyErr != nil {
			code = smithyErr.ErrorCode()
		} else {
			code, ok = awserrors.Code(errors.Cause(err))
			if !ok {
				return false, err
			}
		}

		for _, r := range retryableErrors {
			if code == r {
				// We should retry.
				errToReturn = err
				return false, nil
			}
		}

		// Got an error that we can't retry, so return it.
		return false, err
	})

	// If the waitError is not a timeout error (nil or a non-retryable error), return it
	if !errors.Is(waitErr, wait.ErrorInterrupted(waitErr)) {
		return waitErr
	}

	// A retryable error occurred, return the actual error
	if errToReturn != nil {
		return errToReturn
	}

	// The error was timeout error
	return waitErr
}
