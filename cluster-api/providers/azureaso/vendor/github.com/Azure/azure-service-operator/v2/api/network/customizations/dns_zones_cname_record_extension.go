/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package customizations

import (
	"regexp"

	"github.com/go-logr/logr"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.ErrorClassifier = &DnsZonesCNAMERecordExtension{}

var referencedResourceNotFound = regexp.MustCompile("The referenced resource.*was not found")

func (extension *DnsZonesCNAMERecordExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc,
) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	if isRetryableDNSZoneRecordError(cloudError) {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}

// isRetryableDNSZoneRecordError determines if the error is a retryable DNS Zone Record error.
// This should be generic across all DNS Zone Record types.
func isRetryableDNSZoneRecordError(err *genericarmclient.CloudError) bool {
	if err == nil {
		return false
	}

	// If a referenced resource is not yet provisioned, it may be coming soon
	if err.Code() == "BadRequest" && referencedResourceNotFound.MatchString(err.Message()) {
		return true
	}

	return false
}
