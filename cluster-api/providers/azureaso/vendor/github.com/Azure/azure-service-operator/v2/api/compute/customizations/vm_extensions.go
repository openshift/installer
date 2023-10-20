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

var _ extensions.ErrorClassifier = &VirtualMachineExtension{}

// ClassifyError evaluates the provided error, returning whether it is fatal or can be retried.
func (e *VirtualMachineExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	// It's weird to retry on OperationNotAllowed as it certainly sounds like a fatal error, but
	// it primarily happens for quota errors on VM/VMSS, which we do want to retry on as the quota may free
	// up at some point in the future.
	if details.Code == "OperationNotAllowed" {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}
