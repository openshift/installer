/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package extensions

import (
	"github.com/go-logr/logr"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
)

// ErrorClassifier can be implemented to customize how the reconciler reacts to specific errors
type ErrorClassifier interface {
	// ClassifyError evaluates the provided error, returning including whether it is fatal or can be retried.
	// cloudError is the error returned from ARM.
	// apiVersion is the ARM API version used for the request.
	// log is a logger than can be used for telemetry.
	// next is the next implementation to call.
	ClassifyError(
		cloudError *genericarmclient.CloudError,
		apiVersion string,
		log logr.Logger,
		next ErrorClassifierFunc) (core.CloudErrorDetails, error)
}

// ErrorClassifierFunc is the signature of a function that can be used to create a DefaultErrorClassifier
type ErrorClassifierFunc func(cloudError *genericarmclient.CloudError) (core.CloudErrorDetails, error)

func CreateErrorClassifier(
	host genruntime.ResourceExtension,
	classifier ErrorClassifierFunc,
	apiVersion string,
	log logr.Logger,
) ErrorClassifierFunc {
	impl, ok := host.(ErrorClassifier)
	if !ok {
		return classifier
	}

	return func(cloudError *genericarmclient.CloudError) (core.CloudErrorDetails, error) {
		log.V(Status).Info(
			"Classifying CloudError",
			"Message", cloudError.Message(),
			"Code", cloudError.Code(),
			"Target", cloudError.Target())

		result, err := impl.ClassifyError(cloudError, apiVersion, log, classifier)
		if err != nil {
			log.V(Status).Info(
				"CloudError classification failed",
				"Error", err.Error())

			return core.CloudErrorDetails{}, err
		}

		log.V(Status).Info(
			"CloudError classified",
			"Classification", result.Classification,
			"Code", result.Code,
			"Message", result.Message)

		return result, nil
	}
}
