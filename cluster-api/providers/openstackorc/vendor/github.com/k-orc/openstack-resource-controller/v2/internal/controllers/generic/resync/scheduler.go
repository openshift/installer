/*
Copyright The ORC Authors.

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

package resync

import (
	"errors"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
)

const (
	// jitterFactor is the maximum fraction of extra time added to the base
	// duration. A value of 0.2 produces durations in [base, base*1.2).
	// Positive-only jitter ensures the requeue always fires after
	// resyncPeriod has elapsed, so shouldReconcile always returns true
	// when the requeue fires.
	jitterFactor = 0.2
)

// CalculateJitteredDuration returns a duration in the range [base, base*1.2)
// using uniform random positive-only jitter. Jitter prevents thundering-herd
// problems when many resources share the same resync period.
func CalculateJitteredDuration(base time.Duration) time.Duration {
	return wait.Jitter(base, jitterFactor)
}

// ShouldScheduleResync reports whether a periodic resync should be scheduled
// based on the effective resync period and the current reconcile status.
//
// It returns false (do not schedule) when:
//   - resyncPeriod <= 0: periodic resync is disabled.
//   - reconcileStatus contains a terminal error: the resource is in a
//     non-retryable error state; resync would be pointless.
//   - reconcileStatus already requests a requeue: another reconcile is
//     already pending so a resync requeue would be redundant.
//
// When it returns true, the caller should schedule a requeue after
// CalculateJitteredDuration(resyncPeriod).
func ShouldScheduleResync(resyncPeriod time.Duration, reconcileStatus progress.ReconcileStatus) bool {
	// Resync disabled.
	if resyncPeriod <= 0 {
		return false
	}

	// Terminal error: no further reconciles will help.
	if err := reconcileStatus.GetError(); err != nil {
		var terminalError *orcerrors.TerminalError
		if errors.As(err, &terminalError) {
			return false
		}
	}

	// Another requeue is already pending; avoid adding a redundant one.
	if reconcileStatus.GetRequeue() > 0 {
		return false
	}

	return true
}
