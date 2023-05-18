// Copyright 2023 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dcl

import (
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
)

// NotFoundError is returned when a resource does not exist.
// Some APIs will also return it if a resource may exist but
// the current user does not have permission to view it.
// It wraps an error, usually a *googleapi.Error.
// It maps to HTTP 404.
type NotFoundError struct {
	Cause error
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.Cause)
}

// HasCode returns true if the given error is an HTTP response with the given code.
func HasCode(err error, code int) bool {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == code {
			return true
		}
	}
	return false
}

// IsNotFound returns true if the given error is a NotFoundError or is an HTTP 404.
func IsNotFound(err error) bool {
	if _, ok := err.(NotFoundError); ok {
		return true
	}
	return HasCode(err, 404)
}

// IsNotFoundOrCode returns true if the given error is a NotFoundError, an HTTP 404,
// or an HTTP response with the given code.
func IsNotFoundOrCode(err error, code int) bool {
	return IsNotFound(err) || HasCode(err, code)
}

// EnumInvalidError is returned when an enum is set (by a client) to a string
// value that is not valid for that enum.
// It maps to HTTP 400, although it is usually generated client-side before
// the enum is sent to the server.
type EnumInvalidError struct {
	Enum  string
	Value string
	Valid []string
}

func (e EnumInvalidError) Error() string {
	return fmt.Sprintf("%s not a valid %s (%v)", e.Value, e.Enum, e.Valid)
}

// NotDeletedError is returned when the resource should be deleted but has not
// been.  It is returned if the operation to delete the resource has apparently
// been successful, but Get() still fetches the resource successfully.
type NotDeletedError struct {
	ExistingResource interface{}
}

func (e NotDeletedError) Error() string {
	return fmt.Sprintf("resource not successfully deleted: %#v.", e.ExistingResource)
}

// IsRetryableGoogleError returns true if the error is retryable according to the given retryability.
func IsRetryableGoogleError(gerr *googleapi.Error, retryability Retryability, start time.Time) bool {
	return retryability.Retryable && retryability.regex.MatchString(gerr.Message) && time.Since(start) < retryability.Timeout
}

// IsRetryableHTTPError returns true if the error is retryable - in GCP that's a 500, 502, 503, or 429.
func IsRetryableHTTPError(err error, retryability map[int]Retryability, start time.Time) bool {
	if gerr, ok := err.(*googleapi.Error); ok {
		rtblt, ok := retryability[gerr.Code]
		return ok && IsRetryableGoogleError(gerr, rtblt, start)
	}
	return false
}

// IsNonRetryableHTTPError returns true if we know that the error is not retryable - in GCP that's a 400, 403, 404, or 409.
func IsNonRetryableHTTPError(err error, retryability map[int]Retryability, start time.Time) bool {
	if gerr, ok := err.(*googleapi.Error); ok {
		rtblt, ok := retryability[gerr.Code]
		return ok && !IsRetryableGoogleError(gerr, rtblt, start)
	}
	return false
}

// IsConflictError returns true if the error has conflict error code 409.
func IsConflictError(err error) bool {
	if gerr, ok := err.(*googleapi.Error); ok {
		return gerr.Code == 409
	}
	return false
}

// ApplyInfeasibleError is returned when lifecycle directives prevent an Apply from proceeding.
// This error means that no imperative requests were issued.
type ApplyInfeasibleError struct {
	Message string
}

func (e ApplyInfeasibleError) Error() string {
	return e.Message
}

// DiffAfterApplyError is returned when there are differences between the desired state and the
// intended state after Apply completes.  This usually indicates an error in the SDK, probably
// related to a failure to canonicalize properly.
type DiffAfterApplyError struct {
	Diffs []string
}

func (e DiffAfterApplyError) Error() string {
	return fmt.Sprintf("diffs exist after apply: %v", e.Diffs)
}

// OperationNotDone is returned when an API operation hasn't completed.
// It may wrap an error if the error means that the operation can be retried.
type OperationNotDone struct {
	Err error
}

func (e OperationNotDone) Error() string {
	return "operation not done."
}

// AttemptToIndexNilArray is returned when GetMapEntry is called with a path that includes an array
// index and that array is unset in the map.
type AttemptToIndexNilArray struct {
	FieldName string
}

func (e AttemptToIndexNilArray) Error() string {
	return fmt.Sprintf("field %s was nil, could not index array", e.FieldName)
}
