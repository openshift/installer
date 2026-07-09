/*
Copyright 2020 The ORC Authors.

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

package errors

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
)

var (
	ErrFilterMatch     = fmt.Errorf("filter match error")
	ErrMultipleMatches = multipleMatchesError{}
	ErrNoMatches       = noMatchesError{}
)

type (
	multipleMatchesError struct{}
	noMatchesError       struct{}
)

func (e multipleMatchesError) Error() string {
	return "filter matched more than one resource"
}

func (e multipleMatchesError) Is(err error) bool {
	return err == ErrFilterMatch
}

func (e noMatchesError) Error() string {
	return "filter matched no resources"
}

func (e noMatchesError) Is(err error) bool {
	return err == ErrFilterMatch
}

// IsRetryable returns true if err may succeed when retried without changes to
// the spec. This includes HTTP error responses other than 409 Conflict, since
// some HTTP errors (like 400 Bad Request) can be transient when a dependency
// is not yet ready in OpenStack, and server errors (5xx) are typically
// transient.
//
// Non-HTTP errors from gophercloud (e.g. client-side validation such as
// banned value_spec keys), 409 Conflict, and 501 Not Implemented are not
// retryable. The exception is Neutron quota-exceeded errors, which are
// returned as 409 but are retryable because quota can free up without spec
// changes.
func IsRetryable(err error) bool {
	if IsConflict(err) {
		// Neutron returns 409 for quota-exceeded errors, but these are
		// retryable because quota can free up without spec changes.
		return isNeutronQuotaError(err)
	}

	if IsNotImplementedError(err) {
		return false
	}

	var errUnexpectedResponseCode gophercloud.ErrUnexpectedResponseCode
	return errors.As(err, &errUnexpectedResponseCode)
}

// isNeutronQuotaError returns true if err is an HTTP error response whose
// body indicates a Neutron quota-exceeded condition. Neutron returns quota
// errors as 409 Conflict with an "OverQuota" type in the response body.
func isNeutronQuotaError(err error) bool {
	var errUnexpectedResponseCode gophercloud.ErrUnexpectedResponseCode
	if !errors.As(err, &errUnexpectedResponseCode) {
		return false
	}
	return bytes.Contains(errUnexpectedResponseCode.Body, []byte("OverQuota"))
}

// IsNotFound returns true if err indicates the requested OpenStack resource
// was not found (HTTP 404 or gophercloud's ErrResourceNotFound).
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	// Gophercloud is not consistent in how it returns 404 errors. Sometimes
	// it returns a pointer to the error, sometimes it returns the error
	// directly.
	// Some discussion here: https://github.com/gophercloud/gophercloud/v2/issues/2279
	var errNotFound gophercloud.ErrResourceNotFound
	var pErrNotFound *gophercloud.ErrResourceNotFound
	if errors.As(err, &errNotFound) || errors.As(err, &pErrNotFound) {
		return true
	}

	return gophercloud.ResponseCodeIs(err, http.StatusNotFound)
}

// IsInvalidError returns true if err is an HTTP 400 Bad Request response.
func IsInvalidError(err error) bool {
	return gophercloud.ResponseCodeIs(err, http.StatusBadRequest)
}

// IsConflict returns true if err is an HTTP 409 Conflict response.
func IsConflict(err error) bool {
	return gophercloud.ResponseCodeIs(err, http.StatusConflict)
}

// IsNotImplementedError returns true if err is an HTTP 501 Not Implemented response.
func IsNotImplementedError(err error) bool {
	return gophercloud.ResponseCodeIs(err, http.StatusNotImplemented)
}
