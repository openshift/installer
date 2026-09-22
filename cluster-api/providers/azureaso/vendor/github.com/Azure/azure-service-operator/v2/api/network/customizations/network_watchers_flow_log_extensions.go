// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customizations

import (
	"github.com/go-logr/logr"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.ErrorClassifier = &NetworkWatchersFlowLogExtension{}

// ClassifyError evaluates the provided error, returning whether it is fatal or can be retried.
// A fatal error will be recorded and no further retries will be attempted for that resource.
// A retryable error will be retried after a delay.
func (e *NetworkWatchersFlowLogExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc,
) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	if isRetryableFlowLogError(cloudError) {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}

func isRetryableFlowLogError(err *genericarmclient.CloudError) bool {
	if err == nil {
		return false
	}

	// If the referenced Network Security Group is not yet provisioned, it may be coming soon
	if err.Code() == "NetworkSecurityGroupNotFoundForFlowLog" {
		return true
	}

	// If the reference storage account is not yet provisioned, it may be coming soon
	if err.Code() == "StorageProvisioningStateNotSucceeded" {
		return true
	}

	return false
}
