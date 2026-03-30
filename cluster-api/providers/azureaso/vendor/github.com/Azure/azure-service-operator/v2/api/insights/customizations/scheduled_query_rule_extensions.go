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

var _ extensions.ErrorClassifier = &ScheduledQueryRuleExtension{}

// ClassifyError evaluates the provided error, returning whether it is fatal or can be retried.
// A badrequest (400) is normally fatal, but ScheduledQueryRule resources may return 400 whilst
// a dependency is being created, so we override for that case.
// cloudError is the error returned from ARM.
// apiVersion is the ARM API version used for the request.
// log is a logger than can be used for telemetry.
// next is the next implementation to call.
func (e *ScheduledQueryRuleExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc,
) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	// Override is to treat BadRequests as retryable for ScheduledQueryRules
	if isRetryableError(cloudError) {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}

// isRetryableError checks the passed error to see if it is a retryable bad request, returning true if it is.
func isRetryableError(err *genericarmclient.CloudError) bool {
	if err == nil {
		// No error, so no need for a retry
		return false
	}

	// When any of the scopes to which a ScheduledQueryRule is suposed to apply are missing,
	// we get a BadRequest error with a messaage like:
	// "Scope 'target-law' does not exist"
	// These should be retried to allow the rule to be created once the target scope resource exists.
	if err.Code() == "BadRequest" && strings.Contains(err.Message(), "does not exist") {
		return true
	}

	// When a managed identity is used, if the role assignment was not created yet
	// we get a BadRequest error with a message like:
	// "The provided credentials have insufficient access to perform the requested operation"
	// These should be retried to allow the role assignment to be created once the identity is ready.
	if err.Code() == "BadRequest" && strings.Contains(err.Message(), "insufficient access") {
		return true
	}

	return false
}
