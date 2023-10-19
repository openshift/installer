/*
Copyright 2019 The Kubernetes Authors.

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

package azure

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/go-autorest/autorest"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// ResourceNotFound parses an error to check if its status code is Not Found (404).
func ResourceNotFound(err error) bool {
	return hasStatusCode(err, http.StatusNotFound)
}

// hasStatusCode returns true if an error is a DetailedError or ResponseError with a matching status code.
func hasStatusCode(err error, statusCode int) bool {
	derr := autorest.DetailedError{} // azure-sdk-for-go v1
	if errors.As(err, &derr) {
		return derr.StatusCode == statusCode
	}
	var rerr *azcore.ResponseError // azure-sdk-for-go v2
	return errors.As(err, &rerr) && rerr.StatusCode == statusCode
}

// VMDeletedError is returned when a virtual machine is deleted outside of capz.
type VMDeletedError struct {
	ProviderID string
}

// Error returns the error string.
func (vde VMDeletedError) Error() string {
	return fmt.Sprintf("VM with provider id %q has been deleted", vde.ProviderID)
}

// ReconcileError represents an error that is not automatically recoverable
// errorType indicates what type of action is required to recover. It can take two values:
// 1. `Transient` - Can be recovered through manual intervention, will be requeued after.
// 2. `Terminal` - Cannot be recovered, will not be requeued.
type ReconcileError struct {
	error
	errorType    ReconcileErrorType
	requestAfter time.Duration
}

// ReconcileErrorType represents the type of a ReconcileError.
type ReconcileErrorType string

const (
	// TransientErrorType can be recovered, will be requeued after a configured time interval.
	TransientErrorType ReconcileErrorType = "Transient"
	// TerminalErrorType cannot be recovered, will not be requeued.
	TerminalErrorType ReconcileErrorType = "Terminal"
)

// Error returns the error message for a ReconcileError.
func (t ReconcileError) Error() string {
	var errStr string
	if t.error != nil {
		errStr = t.error.Error()
	}
	switch t.errorType {
	case TransientErrorType:
		return fmt.Sprintf("%s. Object will be requeued after %s", errStr, t.requestAfter.String())
	case TerminalErrorType:
		return fmt.Sprintf("reconcile error that cannot be recovered occurred: %s. Object will not be requeued", errStr)
	default:
		return fmt.Sprintf("reconcile error occurred with unknown recovery type. The actual error is: %s", errStr)
	}
}

// IsTransient returns if the ReconcileError is recoverable.
func (t ReconcileError) IsTransient() bool {
	return t.errorType == TransientErrorType
}

// IsTerminal returns if the ReconcileError is recoverable.
func (t ReconcileError) IsTerminal() bool {
	return t.errorType == TerminalErrorType
}

// Is returns true if the target is a ReconcileError.
func (t ReconcileError) Is(target error) bool {
	return errors.As(target, &ReconcileError{})
}

// RequeueAfter returns requestAfter value.
func (t ReconcileError) RequeueAfter() time.Duration {
	return t.requestAfter
}

// WithTransientError wraps the error in a ReconcileError with errorType as `Transient`.
func WithTransientError(err error, requeueAfter time.Duration) ReconcileError {
	return ReconcileError{error: err, errorType: TransientErrorType, requestAfter: requeueAfter}
}

// WithTerminalError wraps the error in a ReconcileError with errorType as `Terminal`.
func WithTerminalError(err error) ReconcileError {
	return ReconcileError{error: err, errorType: TerminalErrorType}
}

// OperationNotDoneError is used to represent a long-running operation that is not yet complete.
type OperationNotDoneError struct {
	Future *infrav1.Future
}

// NewOperationNotDoneError returns a new OperationNotDoneError wrapping a Future.
func NewOperationNotDoneError(future *infrav1.Future) OperationNotDoneError {
	return OperationNotDoneError{
		Future: future,
	}
}

// Error returns the error represented as a string.
func (onde OperationNotDoneError) Error() string {
	return fmt.Sprintf("operation type %s on Azure resource %s/%s is not done", onde.Future.Type, onde.Future.ResourceGroup, onde.Future.Name)
}

// Is returns true if the target is an OperationNotDoneError.
func (onde OperationNotDoneError) Is(target error) bool {
	return IsOperationNotDoneError(target)
}

// IsOperationNotDoneError returns true if the target is an OperationNotDoneError.
func IsOperationNotDoneError(target error) bool {
	reconcileErr := &ReconcileError{}
	if errors.As(target, reconcileErr) {
		return IsOperationNotDoneError(reconcileErr.error)
	}
	return errors.As(target, &OperationNotDoneError{})
}

// IsContextDeadlineExceededOrCanceledError checks if it's a context deadline
// exceeded or canceled error.
func IsContextDeadlineExceededOrCanceledError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}
