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
	"net/http"

	"github.com/gophercloud/gophercloud"
)

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
	// Some discussion here: https://github.com/gophercloud/gophercloud/issues/2279
	var errDefault404 gophercloud.ErrDefault404
	var pErrDefault404 *gophercloud.ErrDefault404
	var errNotFound gophercloud.ErrResourceNotFound
	var pErrNotFound *gophercloud.ErrResourceNotFound
	if errors.As(err, &errDefault404) || errors.As(err, &pErrDefault404) || errors.As(err, &errNotFound) || errors.As(err, &pErrNotFound) {
		return true
	}

	var errUnexpectedResponseCode gophercloud.ErrUnexpectedResponseCode
	if errors.As(err, &errUnexpectedResponseCode) {
		if errUnexpectedResponseCode.Actual == http.StatusNotFound {
			return true
		}
	}

	return false
}

func IsInvalidError(err error) bool {
	var errDefault400 gophercloud.ErrDefault400
	if errors.As(err, &errDefault400) {
		return true
	}

	var errUnexpectedResponseCode gophercloud.ErrUnexpectedResponseCode
	if errors.As(err, &errUnexpectedResponseCode) {
		if errUnexpectedResponseCode.Actual == http.StatusBadRequest {
			return true
		}
	}

	return false
}

func IsConflict(err error) bool {
	var errDefault409 gophercloud.ErrDefault409
	if errors.As(err, &errDefault409) {
		return true
	}

	var errUnexpectedResponseCode gophercloud.ErrUnexpectedResponseCode
	if errors.As(err, &errUnexpectedResponseCode) {
		if errUnexpectedResponseCode.Actual == http.StatusConflict {
			return true
		}
	}

	return false
}

func IsNotImplementedError(err error) bool {
	var errUnexpectedResponseCode gophercloud.ErrUnexpectedResponseCode
	if errors.As(err, &errUnexpectedResponseCode) {
		if errUnexpectedResponseCode.Actual == http.StatusNotImplemented {
			return true
		}
	}

	return false
}
