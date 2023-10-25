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

var _ extensions.ErrorClassifier = &NamespacesTopicsSubscriptionExtension{}

// ClassifyError evaluates the provided error, returning whether it is fatal or can be retried.
// A MessagingGatewayBadRequest error (400) is normally fatal, but NamespacesTopicsSubscription resource may return MessagingGatewayBadRequest whilst a dependency is being created,
// so we override for that case.
// cloudError is the error returned from ARM.
// apiVersion is the ARM API version used for the request.
// log is a logger than can be used for telemetry.
// next is the next implementation to call.
func (e *NamespacesTopicsSubscriptionExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	// Override is to treat MessagingGatewayBadRequest as retryable for NamespacesTopicsSubscription
	if isRetryableMessagingGatewayBadRequest(cloudError) {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}

// isRetryableMessagingGatewayBadRequest checks the passed error to see if it is a retryable MessagingGatewayBadRequest, returning true if it is.
func isRetryableMessagingGatewayBadRequest(err *genericarmclient.CloudError) bool {
	if err == nil {
		return false
	}

	return err.Code() == "MessagingGatewayBadRequest"
}
