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

var _ extensions.ErrorClassifier = &SecurityPolicyExtension{}

// ClassifyError evaluates the provided error, returning whether it is fatal or can be retried.
// A BadRequest error (400) is normally fatal, but for AFD Security Policies, it is returned if a policy is attempted to be
// added to an endpoint that isn't in a completed provisoningState, which we must retry on.
// cloudError is the error returned from ARM.
// apiVersion is the ARM API version used for the request.
// log is a logger than can be used for telemetry.
// next is the next implementation to call.
func (e *SecurityPolicyExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc,
) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	// Override is to treat Conflict as retryable for Redis, if the message contains "try again later"
	if isSecurityPolicyRetryableBadRequest(cloudError) {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}

// isRetryableBadRequest checks the passed error to see if it is a retryable conflict, returning true if it is.
func isSecurityPolicyRetryableBadRequest(err *genericarmclient.CloudError) bool {
	if err == nil {
		return false
	}

	// string matching for error detection is fragile but unfortunately there's no better way to determine this given the
	// shape the API returns
	return err.Code() == "BadRequest" &&
		strings.Contains(err.Message(), "is not associated with the AFDX profile") &&
		strings.Contains(err.Message(), "is in a invalid provisioning state")
}
