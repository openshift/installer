/*
Copyright 2020 The Kubernetes Authors.

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

func IsRetryable(err error) bool {
	var errUnexpectedResponseCode gophercloud.ErrUnexpectedResponseCode
	if errors.As(err, &errUnexpectedResponseCode) {
		statusCode := errUnexpectedResponseCode.GetStatusCode()
		return statusCode >= 500 && statusCode != http.StatusNotImplemented
	}
	return false
}

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

func IsInvalidError(err error) bool {
	return gophercloud.ResponseCodeIs(err, http.StatusBadRequest)
}

func IsConflict(err error) bool {
	return gophercloud.ResponseCodeIs(err, http.StatusConflict)
}

func IsNotImplementedError(err error) bool {
	return gophercloud.ResponseCodeIs(err, http.StatusNotImplemented)
}

// The following types and constants are imported from CAPI and will be removed at some point once we
// implement the conditions that will be required in CAPI v1beta2
// See https://github.com/kubernetes-sigs/cluster-api/issues/10784

// DeprecatedCAPIMachineStatusError defines errors states for Machine objects.
type DeprecatedCAPIMachineStatusError string

const (
	// DeprecatedCAPIUpdateMachineError indicates an error while trying to update a Node that this
	// Machine represents. This may indicate a transient problem that will be
	// fixed automatically with time, such as a service outage,
	//
	// Example: error updating load balancers.
	DeprecatedCAPIUpdateMachineError DeprecatedCAPIMachineStatusError = "UpdateError"
)

// DeprecatedCAPIClusterStatusError defines errors states for Cluster objects.
type DeprecatedCAPIClusterStatusError string

const (
	// DeprecatedCAPOUpdateClusterError indicates that an error was encountered
	// when trying to update the cluster.
	DeprecatedCAPOUpdateClusterError DeprecatedCAPIClusterStatusError = "UpdateError"
)
