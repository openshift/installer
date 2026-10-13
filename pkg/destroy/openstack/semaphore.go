package openstack

import "sync"

type semaphore struct {
	semC chan struct{}
	wg   sync.WaitGroup
}

// newSemaphore returns a semaphore. A semaphore runs at most maxConcurrency
// functions concurrently.
func newSemaphore(maxConcurrency int) *semaphore {
	return &semaphore{
		semC: make(chan struct{}, maxConcurrency),
	}
}

// Add enqueues the function f to be run in a separate goroutine as soon as
// there is a free slot. Add returns immediately.
func (s *semaphore) Add(f func()) {
	s.wg.Add(1)
	go func() {
		s.semC <- struct{}{}
		defer func() {
			<-s.semC
			s.wg.Done()
		}()
		f()
	}()
}

// Wait returns when the queue is empty and no function is running.
func (s *semaphore) Wait() {
	s.wg.Wait()
}
