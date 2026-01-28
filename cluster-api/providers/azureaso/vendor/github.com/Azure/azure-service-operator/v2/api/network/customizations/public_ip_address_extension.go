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

var _ extensions.ErrorClassifier = &PublicIPAddressExtension{}

func (extension *PublicIPAddressExtension) ClassifyError(
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
	if isRetryablePublicIPAddressError(cloudError) {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}

func isRetryablePublicIPAddressError(err *genericarmclient.CloudError) bool {
	if err == nil {
		return false
	}

	// If a referenced resource is not yet provisioned, it may be coming soon
	if err.Code() == "ReferencedResourceNotProvisioned" {
		return true
	}

	return false
}
