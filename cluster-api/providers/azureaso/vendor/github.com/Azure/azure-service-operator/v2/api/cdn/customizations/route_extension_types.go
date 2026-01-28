/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package customizations

import (
	"strings"

	"github.com/go-logr/logr"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.ErrorClassifier = &RouteExtension{}

// ClassifyError evaluates the provided error, returning whether it is fatal or can be retried.
// A BadRequest error (400) is normally fatal, but for AFD Routes, it is returned if a route is attempted to be
// added to an originGroup that doesn't exist or doesn't have any origins added. This error is not actually fatal
// so we retry on it.
// cloudError is the error returned from ARM.
// apiVersion is the ARM API version used for the request.
// log is a logger than can be used for telemetry.
// next is the next implementation to call.
func (e *RouteExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc,
) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	if isRouteRetryableBadRequest(cloudError) {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}

// isRouteRetryableBadRequest checks the passed error to see if it is a retryable error, returning true if it is.
func isRouteRetryableBadRequest(err *genericarmclient.CloudError) bool {
	if err == nil {
		return false
	}

	// string matching for error detection is fragile but unfortunately there's no better way to determine this given the
	// shape the API returns
	return err.Code() == "BadRequest" &&
		strings.Contains(err.Message(), "originGroup is created successfully") &&
		strings.Contains(err.Message(), "at least one enabled origin is created")
}
