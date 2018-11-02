package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/watch"
	watchtools "k8s.io/client-go/tools/watch"
)

// WatcherFunc is from https://github.com/kubernetes/kubernetes/pull/50102.
type WatcherFunc func(sinceResourceVersion string) (watch.Interface, error)

type resourceVersionGetter interface {
	GetResourceVersion() string
}

// RetryWatcher is from https://github.com/kubernetes/kubernetes/pull/50102.
type RetryWatcher struct {
	lastResourceVersion string
	watcherFunc         WatcherFunc
	resultChan          chan watch.Event
	stopChan            chan struct{}
	doneChan            chan struct{}
}

// Until is from https://github.com/kubernetes/kubernetes/pull/50102.
func Until(ctx context.Context, initialResourceVersion string, watcherFunc WatcherFunc, conditions ...watchtools.ConditionFunc) (*watch.Event, error) {
	return watchtools.UntilWithoutRetry(ctx, NewRetryWatcher(initialResourceVersion, watcherFunc), conditions...)
}

// NewRetryWatcher is from https://github.com/kubernetes/kubernetes/pull/50102.
func NewRetryWatcher(initialResourceVersion string, watcherFunc WatcherFunc) *RetryWatcher {
	rw := &RetryWatcher{
		lastResourceVersion: initialResourceVersion,
		watcherFunc:         watcherFunc,
		stopChan:            make(chan struct{}),
		doneChan:            make(chan struct{}),
		resultChan:          make(chan watch.Event, 0),
	}
	go rw.receive()
	return rw
}

func (rw *RetryWatcher) send(event watch.Event) bool {
	// Writing to an unbuffered channel is blocking and we need to check if we need to be able to stop while doing so!
	select {
	case rw.resultChan <- event:
		return true
	case <-rw.stopChan:
		return false
	}
}

func (rw *RetryWatcher) doReceive() bool {
	watcher, err := rw.watcherFunc(rw.lastResourceVersion)
	if err != nil {
		status := apierrors.NewInternalError(fmt.Errorf("retry watcher: watcherFunc failed: %v", err)).Status()
		_ = rw.send(watch.Event{
			Type:   watch.Error,
			Object: &status,
		})
		// Stop the watcher
		return true
	}
	ch := watcher.ResultChan()
	defer watcher.Stop()

	for {
		select {
		case <-rw.stopChan:
			logrus.Debug("Stopping RetryWatcher.")
			return true
		case event, ok := <-ch:
			if !ok {
				logrus.Warningf("RetryWatcher - getting event failed! Re-creating the watcher. Last RV: %s", rw.lastResourceVersion)
				return false
			}

			// We need to inspect the event and get ResourceVersion out of it
			switch event.Type {
			case watch.Added, watch.Modified, watch.Deleted:
				metaObject, ok := event.Object.(resourceVersionGetter)
				if !ok {
					status := apierrors.NewInternalError(errors.New("__internal__: RetryWatcher: doesn't support resourceVersion")).Status()
					_ = rw.send(watch.Event{
						Type:   watch.Error,
						Object: &status,
					})
					// We have to abort here because this might cause lastResourceVersion inconsistency by skipping a potential RV with valid data!
					return true
				}

				resourceVersion := metaObject.GetResourceVersion()
				if resourceVersion == "" {
					status := apierrors.NewInternalError(fmt.Errorf("__internal__: RetryWatcher: object %#v doesn't support resourceVersion", event.Object)).Status()
					_ = rw.send(watch.Event{
						Type:   watch.Error,
						Object: &status,
					})
					// We have to abort here because this might cause lastResourceVersion inconsistency by skipping a potential RV with valid data!
					return true
				}

				// All is fine; send the event and update lastResourceVersion
				ok = rw.send(event)
				if !ok {
					return true
				}
				rw.lastResourceVersion = resourceVersion

				continue

			case watch.Error:
				_ = rw.send(event)
				return true

			default:
				logrus.Errorf("RetryWatcher failed to recognize Event type %q", event.Type)
				status := apierrors.NewInternalError(fmt.Errorf("__internal__: RetryWatcher failed to recognize Event type %q", event.Type)).Status()
				_ = rw.send(watch.Event{
					Type:   watch.Error,
					Object: &status,
				})
				// We are unable to restart the watch and have to stop the loop or this might cause lastResourceVersion inconsistency by skipping a potential RV with valid data!
				return true
			}
		}
	}
}

func (rw *RetryWatcher) receive() {
	defer close(rw.doneChan)

	for {
		select {
		case <-rw.stopChan:
			logrus.Debug("Stopping RetryWatcher.")
			return
		default:
			done := rw.doReceive()
			if done {
				return
			}
		}
	}
}

// ResultChan is from https://github.com/kubernetes/kubernetes/pull/50102.
func (rw *RetryWatcher) ResultChan() <-chan watch.Event {
	return rw.resultChan
}

// Stop is from https://github.com/kubernetes/kubernetes/pull/50102.
func (rw *RetryWatcher) Stop() {
	close(rw.stopChan)
}

// Done is from https://github.com/kubernetes/kubernetes/pull/50102.
func (rw *RetryWatcher) Done() <-chan struct{} {
	return rw.doneChan
}
