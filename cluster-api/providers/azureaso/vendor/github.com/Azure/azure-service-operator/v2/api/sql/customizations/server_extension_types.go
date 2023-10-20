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

var _ extensions.ErrorClassifier = &ServerExtension{}

// ClassifyError evaluates the provided error, returning whether it is fatal or can be retried.
// cloudError is the error returned from ARM.
// apiVersion is the ARM API version used for the request.
// log is a logger than can be used for telemetry.
// next is the next implementation to call.
func (e *ServerExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	if cloudError.Code() == "PasswordNotComplex" {
		details.Classification = core.ErrorFatal
	}

	return details, nil
}
