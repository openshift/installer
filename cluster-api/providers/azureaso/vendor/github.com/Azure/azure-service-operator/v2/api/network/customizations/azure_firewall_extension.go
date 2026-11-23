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

var firewallReferenceNotFound = regexp.MustCompile("(?i)failed to reference Firewall Policy")

func (extension *AzureFirewallExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc,
) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	if isRetryableFirewallError(cloudError) {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}

func isRetryableFirewallError(err *genericarmclient.CloudError) bool {
	if err == nil {
		return false
	}

	// If a referenced resource is not yet provisioned, it may be coming soon
	if err.Code() == "AzfwAddToFirewallPolicyFailed" && firewallReferenceNotFound.MatchString(err.Message()) {
		return true
	}

	return false
}
