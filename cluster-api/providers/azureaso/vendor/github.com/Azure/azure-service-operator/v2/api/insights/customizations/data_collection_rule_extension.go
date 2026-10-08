/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package customizations

import (
	"github.com/go-logr/logr"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.ErrorClassifier = &DataCollectionRuleExtension{}

// ClassifyError evaluates the provided error, returning whether it is fatal or can be retried.
// A badrequest (400) is normally fatal, but DataCollectionRule resources may return 400 whilst
// a dependency is being created, so we override for that case.
// cloudError is the error returned from ARM.
// apiVersion is the ARM API version used for the request.
// log is a logger than can be used for telemetry.
// next is the next implementation to call.
func (e *DataCollectionRuleExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc,
) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	// Override is to treat BadRequests as retryable for DataCollectionRules
	if isRetryableErrorForDataCollectionRule(cloudError) {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}

// isRetryableError checks the passed error to see if it is a retryable bad request, returning true if it is.
func isRetryableErrorForDataCollectionRule(err *genericarmclient.CloudError) bool {
	if err == nil {
		// No error, so no need for a retry
		return false
	}

	// If the workspace required by a DataCollectionRules is missing, we get a 400 response with code WorkspaceNotFound
	// We can retry as the Workspace may already be under construction
	if err.Code() == "WorkspaceNotFound" {
		return true
	}

	return false
}
