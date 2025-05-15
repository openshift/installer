// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package lockedrand

import (
	"math/rand"
	"sync"
)

// LockedSource is a simple locked random source, similar to the one implemented but
// not exported in math/rand.
type LockedSource struct {
	lock sync.Mutex
	src  rand.Source64
}

var (
	_ rand.Source   = &LockedSource{}
	_ rand.Source64 = &LockedSource{}
)

func (l *LockedSource) Int63() int64 {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.src.Int63()
}

func (l *LockedSource) Uint64() uint64 {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.src.Uint64()
}

func (l *LockedSource) Seed(seed int64) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.src.Seed(seed)
}

// NewSource returns a new LockedSource
func NewSource(seed int64) rand.Source {
	return &LockedSource{
		src: rand.NewSource(seed).(rand.Source64),
	}
}
