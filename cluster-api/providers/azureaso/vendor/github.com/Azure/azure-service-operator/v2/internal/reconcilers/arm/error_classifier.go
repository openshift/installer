/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package arm

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
	case "AnotherOperationInProgress",
		"AuthorizationFailed",
		"AllocationFailed",
		"Conflict",
		"FailedIdentityOperation",
		"InvalidResourceReference",
		"InvalidSubscriptionRegistrationState",
		"LinkedAuthorizationFailed",
		"MissingRegistrationForLocation",
		"MissingSubscriptionRegistration",
		"NoRegisteredProviderFound",
		"NotFound",
		"ParentResourceNotFound",
		"PrincipalNotFound",
		"ResourceGroupNotFound",
		"ResourceNotFound",
		"ResourceQuotaExceeded",
		"SubscriptionNotRegistered":
		result.Classification = core.ErrorRetryable
	case "ScopeLocked": // This error is raised when a resource or resource(s) are locked by a lock
		// Retry, but at a very slow rate to avoid spamming the Azure API with a request that is unlikely to succeed.
		result.Classification = core.ErrorRetryable
		result.Retry = retry.VerySlow
	case "BadRequestFormat",
		"BadRequest",
		"PublicIpForGatewayIsRequired", // TODO: There's not a great way to look at an arbitrary error returned by this API and determine if it's a 4xx or 5xx level... ugh
		"InvalidParameter",
		"InvalidParameterValue",
		"InvalidResourceGroupLocation",
		"InvalidResourceType",
		"InvalidRequestContent",
		"InvalidTemplate",
		"InvalidValuesForRequestParameters",
		"InvalidGatewaySkuProvidedForGatewayVpnType",
		"InvalidGatewaySize",
		"LocationRequired",
		"MethodNotAllowed",
		"MissingRequiredParameter",
		"PasswordTooLong",
		"PrivateIPAddressInReservedRange",
		"PrivateIPAddressNotInSubnet",
		"PropertyChangeNotAllowed",
		"RegionIsOfferRestricted",
		"RequestDisallowedByPolicy", // TODO: Technically could probably retry through this?
		"ReservedResourceName",
		"SkuNotAvailable",
		"SubscriptionNotFound":
		result.Classification = core.ErrorFatal
	default:
		// If we don't know what the error is use the HTTP status code to determine if we can retry
		result.Classification = classifyHTTPError(err)
	}

	return result
}

func classifyHTTPError(err *genericarmclient.CloudError) core.ErrorClassification {
	var httpError *azcore.ResponseError
	if !eris.As(err.Unwrap(), &httpError) {
		return core.ErrorRetryable
	}

	if httpError.StatusCode == 400 {
		return core.ErrorFatal
	}

	return core.ErrorRetryable
}
