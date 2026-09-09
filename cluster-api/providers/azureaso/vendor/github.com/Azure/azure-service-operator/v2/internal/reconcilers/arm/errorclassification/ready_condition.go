/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package errorclassification

import (
	"github.com/rotisserie/eris"

	"github.com/Azure/azure-service-operator/v2/internal/errwrap"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/retry"
)

// MakeReadyConditionImpactingErrorFromError creates a ReadyConditionImpactingError from an error.
// It uses the provided classifier to classify the error, and returns a ReadyConditionImpactingError with the
// appropriate severity and reason based on the classification.
// If the error is already a ReadyConditionImpactingError, it returns it unchanged.
func MakeReadyConditionImpactingErrorFromError(azureErr error, classifier extensions.ErrorClassifierFunc) error {
	var readyConditionError *conditions.ReadyConditionImpactingError
	isReadyConditionImpactingError := eris.As(azureErr, &readyConditionError)
	if isReadyConditionImpactingError {
		// The error has already been classified. This currently only happens in test with the go-vcr injected
		// http client
		return azureErr
	}

	var cloudError *genericarmclient.CloudError
	isCloudErr := eris.As(azureErr, &cloudError)
	if !isCloudErr {
		// This shouldn't happen, as all errors from ARM should be in one of the shapes that CloudError supports. In case
		// we've somehow gotten one that isn't formatted correctly, create a sensible default error
		return conditions.NewReadyConditionImpactingError(
			azureErr,
			conditions.ConditionSeverityWarning,
			conditions.MakeReason(core.UnknownErrorCode, retry.Slow))
	}

	details, err := classifier(cloudError)
	if err != nil {
		return eris.Wrapf(
			err,
			"Unable to classify cloud error (%s)",
			cloudError.Error())
	}

	var severity conditions.ConditionSeverity
	switch details.Classification {
	case core.ErrorRetryable:
		severity = conditions.ConditionSeverityWarning
	case core.ErrorFatal:
		severity = conditions.ConditionSeverityError
	default:
		return eris.Errorf(
			"unknown error classification %q while making Ready condition",
			details.Classification)
	}

	// Stick errorDetails.Message into an error so that it will be displayed as the message on the condition
	err = errwrap.Hide(cloudError, details.Message)
	reason := conditions.MakeReason(details.Code, details.Retry)
	result := conditions.NewReadyConditionImpactingError(err, severity, reason)

	return result
}
