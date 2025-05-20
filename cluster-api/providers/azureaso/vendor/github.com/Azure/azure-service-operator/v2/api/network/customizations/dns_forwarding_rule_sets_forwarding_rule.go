// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package customizations

import (
	"strings"

	"github.com/go-logr/logr"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.ErrorClassifier = &DnsForwardingRuleSetsForwardingRuleExtension{}

func (extension *DnsForwardingRuleSetsForwardingRuleExtension) ClassifyError(
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
	if isRetryableConflict(cloudError) {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}

// isRetryableConflict checks the passed error to see if it is a retryable conflict, returning true if it is.
func isRetryableConflict(err *genericarmclient.CloudError) bool {
	if err == nil {
		return false
	}

	// We retry on this case as it's not actually fatal
	return strings.Contains(err.Message(), "Another operation is pending on the specified resource")
}
