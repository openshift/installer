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

package elb

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
)

var _ error = &ELBError{}

// ELBError is an error exposed to users of this library.
type ELBError struct {
	msg string

	Code int
}

// Error implements the Error interface.
func (e *ELBError) Error() string {
	return e.msg
}

// NewNotFound returns an error which indicates that the resource of the kind and the name was not found.
func NewNotFound(msg string) error {
	return &ELBError{
		msg:  msg,
		Code: http.StatusNotFound,
	}
}

// NewConflict returns an error which indicates that the request cannot be processed due to a conflict.
func NewConflict(msg string) error {
	return &ELBError{
		msg:  msg,
		Code: http.StatusConflict,
	}
}

// NewInstanceNotRunning returns an error which indicates that the request cannot be processed due to the instance not
// being in a running state.
func NewInstanceNotRunning(msg string) error {
	return &ELBError{
		msg:  msg,
		Code: http.StatusTooEarly,
	}
}

// IsNotFound returns true if the error was created by NewNotFound.
func IsNotFound(err error) bool {
	if ReasonForError(err) == http.StatusNotFound {
		return true
	}
	if code, ok := awserrors.Code(errors.Cause(err)); ok {
		if code == elb.ErrCodeAccessPointNotFoundException {
			return true
		}
	}
	return false
}

// IsAccessDenied returns true if the error is AccessDenied.
func IsAccessDenied(err error) bool {
	if code, ok := awserrors.Code(errors.Cause(err)); ok {
		if code == "AccessDenied" {
			return true
		}
	}
	return false
}

// IsConflict returns true if the error was created by NewConflict.
func IsConflict(err error) bool {
	return ReasonForError(err) == http.StatusConflict
}

// IsSDKError returns true if the error is of type awserr.Error.
func IsSDKError(err error) (ok bool) {
	_, ok = errors.Cause(err).(awserr.Error)
	return
}

// IsInstanceNotRunning returns true if the error was created by NewInstanceNotRunning.
func IsInstanceNotRunning(err error) (ok bool) {
	return ReasonForError(err) == http.StatusTooEarly
}

// ReasonForError returns the HTTP status for a particular error.
func ReasonForError(err error) int {
	if t, ok := errors.Cause(err).(*ELBError); ok {
		return t.Code
	}
	return -1
}
