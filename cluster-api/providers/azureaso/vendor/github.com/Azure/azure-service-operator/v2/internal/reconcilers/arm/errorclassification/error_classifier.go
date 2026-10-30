/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package errorclassification

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/rotisserie/eris"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/retry"
)

func ClassifyCloudError(err *genericarmclient.CloudError) (core.CloudErrorDetails, error) {
	if err == nil {
		// Default to retrying if we're asked to classify a nil error
		result := core.CloudErrorDetails{
			Classification: core.ErrorRetryable,
			Code:           core.UnknownErrorCode,
			Message:        core.UnknownErrorMessage,
		}
		return result, nil
	}

	return classifyCloudError(err), nil
}

func classifyCloudError(err *genericarmclient.CloudError) core.CloudErrorDetails {
	// See https://docs.microsoft.com/en-us/azure/azure-resource-manager/templates/common-deployment-errors
	// for a breakdown of common deployment error codes. Note that the error codes documented there are
	// the inner error codes we're parsing here.

	code := err.Code()
	if code == "" {
		// If there's no code, assume we can retry on it
		return core.CloudErrorDetails{
			Classification: core.ErrorRetryable,
			Code:           err.Code(),
			Message:        err.Message(),
		}
	}

	result := core.CloudErrorDetails{
		Code:    err.Code(),
		Message: err.Message(),
	}

	switch code {
	case "AnotherOperationInProgress", // Retry after the other operation has completed may work
		"AuthorizationFailed",                  // Retry after authorization updated (or propagated) may work
		"AllocationFailed",                     // Retry later may work as there may be more capacity then
		"Conflict",                             // Retry after conflicting operation completed may work
		"FailedIdentityOperation",              // Retry after identity updated (or propagated) may work
		"InvalidResourceReference",             // Often raised when the referenced resource is still being created, retry once creation has finished may work
		"InvalidSubscriptionRegistrationState", // Retry after registration updated (or propagated) may work
		"LinkedAuthorizationFailed",            // Retry after authorization updated (or propagated) may work
		"MissingRegistrationForLocation",       // Retry after registration updated (or propagated) may work
		"MissingSubscriptionRegistration",      // Retry after registration updated (or propagated) may work
		"NoRegisteredProviderFound",            // Retry after provider registered may work
		"NotFound",                             // Retry after resource created may work
		"ParentResourceNotFound",               // Retry after parent created may work
		"PrincipalNotFound",                    // Retry after principal propagation may work
		"ResourceGroupNotFound",                // Often raised when the referenced group is still being created, retry once creation has finished may work
		"ResourceNotFound",                     // Often raised when the referenced resource is still being created, retry once creation has finished may work
		"ResourceQuotaExceeded",                // Retry after quota increased (or other resources deleted) may work
		"RequestDisallowedByPolicy",            // Retry after policy updated/removed may work
		"SubscriptionNotRegistered":            // Retry after subscription registration updated (or propagated) may work
		result.Classification = core.ErrorRetryable

	// Codes here are probably fatal, but we've seen reports that they can come up transiently from time to time,
	// so we treat them as retryable at a slow rate.
	case "BadRequestFormat",
		"BadRequest",
		"InvalidResourceGroupLocation",
		"InvalidParameter",
		"InvalidParameterValue",
		"MethodNotAllowed",
		"ReservedResourceName",
		"RegionIsOfferRestricted",
		"ScopeLocked", // This error is raised when a resource or resource(s) are locked by a lock
		"SkuNotAvailable",
		"SubscriptionNotFound":
		// Retry, but at a very slow rate to avoid spamming the Azure API with a request that is unlikely to succeed.
		result.Classification = core.ErrorRetryable
		result.Retry = retry.VerySlow

	case "PublicIpForGatewayIsRequired", // TODO: There's not a great way to look at an arbitrary error returned by this API and determine if it's a 4xx or 5xx level... ugh
		"InvalidResourceType",
		"InvalidRequestContent",
		"InvalidTemplate",
		"InvalidValuesForRequestParameters",
		"InvalidGatewaySkuProvidedForGatewayVpnType",
		"InvalidGatewaySize",
		"K8sVersionNotSupported",
		"LocationRequired",
		"MissingRequiredParameter",
		"PasswordTooLong",
		"PrivateIPAddressInReservedRange",
		"PrivateIPAddressNotInSubnet",
		"PropertyChangeNotAllowed":
		result.Classification = core.ErrorFatal

	default:
		// If we don't know what the error is use the HTTP status code to determine if we can retry
		result.Classification, result.Retry = classifyHTTPError(err)
	}

	return result
}

func classifyHTTPError(err *genericarmclient.CloudError) (core.ErrorClassification, retry.Classification) { //nolint:unparam
	var httpError *azcore.ResponseError
	if !eris.As(err.Unwrap(), &httpError) {
		// If we can't determine the error type, assume we can retry
		return core.ErrorRetryable, retry.Slow
	}

	if httpError.StatusCode == 400 {
		// HTTP 400 errors are generally fatal, but we don't know that for sure and Azure has a lot of services which return
		// 400 even when the error is transient, so we treat them all as retryable (at very slow speed) for now.
		return core.ErrorRetryable, retry.VerySlow
	}

	return core.ErrorRetryable, retry.Slow
}
