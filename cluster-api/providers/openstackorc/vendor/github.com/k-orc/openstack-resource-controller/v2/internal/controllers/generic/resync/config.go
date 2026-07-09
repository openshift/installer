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

// Package resync provides helpers for determining the effective resync period
// for ORC controllers, implementing the configuration resolution hierarchy.
package resync

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DetermineResyncPeriod resolves the effective resync period using the
// following hierarchy:
//
//  1. If specValue is non-nil and non-zero, return its duration (per-resource
//     override takes precedence).
//  2. If specValue is explicitly zero (0s), return 0 (resync is disabled
//     regardless of the global default).
//  3. If specValue is nil, return globalDefault.
//
// A return value of 0 means periodic resync is disabled.
func DetermineResyncPeriod(specValue *metav1.Duration, globalDefault time.Duration) time.Duration {
	if specValue != nil {
		// Explicit spec value: use it unconditionally (zero means disabled).
		return specValue.Duration
	}
	// No per-resource override: fall back to the global default.
	return globalDefault
}

// RemainingUntilNextSync returns how long remains before the next periodic
// resync is due. It returns 0 when periodic resync is disabled, lastSyncTime is
// unset, or the period has already elapsed.
func RemainingUntilNextSync(lastSyncTime *metav1.Time, resyncPeriod time.Duration) time.Duration {
	if resyncPeriod <= 0 || lastSyncTime == nil {
		return 0
	}

	remaining := resyncPeriod - time.Since(lastSyncTime.Time)
	if remaining <= 0 {
		return 0
	}
	return remaining
}
