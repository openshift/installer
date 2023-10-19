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

var _ extensions.ErrorClassifier = &MongodbDatabaseExtension{}

// ClassifyError evaluates the provided error, returning whether it is fatal or can be retried.
// A BadRequest (400) is normally fatal, but Mongodb Databases may return 400 if database creation is attempted while
// the parent account is still being created, so we make BadRequest retryable for this case.
// cloudError is the error returned from ARM.
// next is the next implementation to call.
func (extension *MongodbDatabaseExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	_ string,
	_ logr.Logger,
	next extensions.ErrorClassifierFunc) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	// Override is to treat BadRequest as retryable for Mongodb Databases
	if isRetryableBadRequest(cloudError) {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}
