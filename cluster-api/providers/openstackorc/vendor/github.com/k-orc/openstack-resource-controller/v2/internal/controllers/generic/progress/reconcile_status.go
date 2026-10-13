/*
Copyright 2024 The ORC Authors.

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

package progress

import (
	"errors"
	"fmt"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/go-logr/logr"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
)

// ReconcileStatus represents the status of the current reconcile.
//
// nil is a valid ReconcileStatus. It is safe to call all methods on a nil ReconcileStatus.
//
// You MUST use the return value of any method which returns ReconcileStatus:
// the returned object may not be equal to the receiver, for example when the
// receiver is nil.
type ReconcileStatus = *reconcileStatus

type reconcileStatus struct {
	messages []string
	requeue  time.Duration

	err error
}

// NewReconcileStatus returns an empty ReconcileStatus
func NewReconcileStatus() ReconcileStatus {
	return nil
}

// WrapError returns a ReconcileStatus containing the given error
func WrapError(err error) ReconcileStatus {
	return NewReconcileStatus().WithError(err)
}

// GetProgressMessages returns all progress messages which have been added to this ReconcileStatus
func (r ReconcileStatus) GetProgressMessages() []string {
	if r == nil {
		return nil
	}

	return r.messages
}

// GetError returns an error representing all errors which have been added to
// this ReconcileStatus. If multiple errors have been added they will have been
// combined with errors.Join()
func (r ReconcileStatus) GetError() error {
	if r == nil {
		return nil
	}

	return r.err
}

// GetRequeue returns the time after which the current object should be
// reconciled again. A value of 0 indicates that no requeue is requested.
func (r ReconcileStatus) GetRequeue() time.Duration {
	if r == nil {
		return 0
	}

	return r.requeue
}

// NeedsReschedule returns a boolean value indicating whether the
// ReconcileStatus will set the Progressing condition to true, and therefore
// that we intend to be scheduled again. It additionally returns any error
// associated with the ReconcileStatus.
//
// NeedsReschedule is used to shortcut reconciliation if any precondition has
// not been met.
func (r ReconcileStatus) NeedsReschedule() (bool, error) {
	if r == nil {
		return false, nil
	}

	return len(r.messages) > 0 || r.err != nil, r.err
}

// Return returns the the (ctrl.Result, error) expected by controller-runtime
// for a ReconcileStatus.
//
// If a ReconcileStatus contains a TerminalError, Return will log the error
// directly instead of returning it to controller-runtime, as this would cause
// an undesirable reschedule.
func (r ReconcileStatus) Return(log logr.Logger) (ctrl.Result, error) {
	if r == nil {
		return ctrl.Result{}, nil
	}

	var terminalError *orcerrors.TerminalError
	if errors.As(r.err, &terminalError) {
		log.V(logging.Info).Info("not scheduling further reconciles for terminal error", "err", r.err.Error())
		return ctrl.Result{}, nil
	}

	if r.err != nil {
		return ctrl.Result{}, r.err
	}

	return ctrl.Result{RequeueAfter: r.requeue}, nil
}

// WithProgressMessage returns a ReconcileStatus with the given progress
// messages in addition to any already present.
func (r ReconcileStatus) WithProgressMessage(msgs ...string) ReconcileStatus {
	if len(msgs) == 0 {
		return r
	}

	if r == nil {
		r = &reconcileStatus{}
	}

	r.messages = append(r.messages, msgs...)
	return r
}

// WithRequeue returns a ReconcileStatus with a request to requeue after the
// given time. If the ReconcileStatus already requests a requeue, the returned
// object will have the lesser of the existing and requested requeues.
func (r ReconcileStatus) WithRequeue(requeue time.Duration) ReconcileStatus {
	if requeue == 0 {
		return r
	}

	if r == nil {
		r = &reconcileStatus{}
	}

	if r.requeue == 0 || requeue < r.requeue {
		r.requeue = requeue
	}
	return r
}

// WithError returns a ReconcileStatus containing the given error joined to any
// existing errors.
func (r ReconcileStatus) WithError(err error) ReconcileStatus {
	if err == nil {
		return r
	}

	if r == nil {
		r = &reconcileStatus{}
	}

	r.err = errors.Join(r.err, err)
	return r
}

// WithReconcileStatus returns a ReconcileStatus combining all properties of the given ReconcileStatus.
func (r ReconcileStatus) WithReconcileStatus(o ReconcileStatus) ReconcileStatus {
	if r == nil {
		return o
	}

	return r.WithProgressMessage(o.GetProgressMessages()...).
		WithRequeue(o.GetRequeue()).
		WithError(o.GetError())
}

// WaitingOnEvent represents the type of event we are waiting on
type WaitingOnEvent int

const (
	// WaitingOnCreation indicates waiting for an object to be created
	WaitingOnCreation WaitingOnEvent = iota

	// WaitingOnReady indicates that an object exists but is not yet in the necessary state
	WaitingOnReady

	// WaitingOnDeletion indicates waiting for an object to be deleted
	WaitingOnDeletion
)

// WaitingOnObject adds a progress message indicating that we are waiting on a
// kubernetes object of type kind with name. We expect the controller to have an
// appropriate watch and handler for this event, so WaitingOnObject does not add
// an explicit requeue.
func (r ReconcileStatus) WaitingOnObject(kind, name string, waitingOn WaitingOnEvent) ReconcileStatus {
	var outcome string
	switch waitingOn {
	case WaitingOnCreation:
		outcome = "created"
	case WaitingOnReady:
		outcome = "ready"
	case WaitingOnDeletion:
		outcome = "deleted"
	}
	return r.WithProgressMessage(fmt.Sprintf("Waiting for %s/%s to be %s", kind, name, outcome))
}

// WaitingOnObject is a convenience method which returns a new ReconcileStatus with WaitingOnObject.
func WaitingOnObject(kind, name string, waitingOn WaitingOnEvent) ReconcileStatus {
	return NewReconcileStatus().WaitingOnObject(kind, name, waitingOn)
}

// WaitingOnFinalizer adds a progress message indicating that we are waiting for a specific finalizer to be removed.
func (r ReconcileStatus) WaitingOnFinalizer(finalizer string) ReconcileStatus {
	return r.WithProgressMessage(fmt.Sprintf("Waiting for finalizer %s to be removed", finalizer))
}

// WaitingOnFinalizer is a convenience method which returns a new ReconcileStatus with WaitingOnFinalizer.
func WaitingOnFinalizer(finalizer string) ReconcileStatus {
	return NewReconcileStatus().WaitingOnFinalizer(finalizer)
}

// WaitingOnOpenStack indicates that we are waiting for an event on the current
// OpenStack resource. It adds an appropriate progress message. It also adds a
// requeue with the requested polling period, as we are not able to receive
// triggers for OpenStack events.
func (r ReconcileStatus) WaitingOnOpenStack(waitingOn WaitingOnEvent, pollingPeriod time.Duration) ReconcileStatus {
	var outcome string
	switch waitingOn {
	case WaitingOnCreation:
		outcome = "created externally"
	case WaitingOnReady:
		outcome = "ready"
	case WaitingOnDeletion:
		outcome = "deleted"
	}

	return r.WithProgressMessage(fmt.Sprintf("Waiting for OpenStack resource to be %s", outcome)).
		WithRequeue(pollingPeriod)
}

// WaitingOnOpenStack is a convenience method which returns a new ReconcileStatus with WaitingOnOpenStack.
func WaitingOnOpenStack(waitingOn WaitingOnEvent, pollingPeriod time.Duration) ReconcileStatus {
	return NewReconcileStatus().WaitingOnOpenStack(waitingOn, pollingPeriod)
}

// NeedsRefresh indicates that the resource status needs to be refreshed. It
// sets an appropriate progress message and ensures that the object will be
// reconciled again immediately.
func (r ReconcileStatus) NeedsRefresh() ReconcileStatus {
	return r.WithProgressMessage("Resource status will be refreshed")
}

// NeedsRefresh is a convenience method which returns a new ReconcileStatus with NeedsRefresh.
func NeedsRefresh() ReconcileStatus {
	return NewReconcileStatus().NeedsRefresh()
}
