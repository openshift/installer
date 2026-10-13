// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package randextensions

import (
	"math"
	"math/rand"
	"time"
)

// Jitter returns the provided duration, scaled by a jitter factor.
// jitter must be between 0 and 1.
// After jitter is applied the result will be in the range: [(1-jitter)t, (1+jitter)t],
// UNLESS (1+jitter)t is larger than MaxInt64, in which case the range will be [(1-jitter)t, math.MaxInt64]
func Jitter(r *rand.Rand, t time.Duration, jitter float64) time.Duration {
	if jitter <= 0 || jitter > 1 {
		panic("jitter must be in the range (0, 1]")
	}
	halfRange := int64(float64(t) * jitter)
	min := int64(t) - halfRange
	fullRange := 2 * halfRange
	// check for possible overflow and limit range if so
	if int64(t) > math.MaxInt64-halfRange {
		fullRange = math.MaxInt64 - int64(t) + halfRange
	}

	val := min + int64(r.Float64()*float64(fullRange))
	return time.Duration(val)
}
