/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package extensions

import (
	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// SuccessfulCreationHandler can be implemented to customize the resource upon successful creation in Azure.
// This extension is invoked once after the initial ARM PUT operation succeeds, giving resources the opportunity
// to perform one-time initialization that depends on the Azure-assigned resource ID.
// Implement this extension when:
// - Resource ID needs custom computation or override after creation
// - Child resources need special ID references set on the parent
// - One-time initialization is required after the resource exists in Azure
type SuccessfulCreationHandler interface {
	// Success performs custom logic after the resource is successfully created in Azure for the first time.
	// obj is the resource that was just created, with populated status including the Azure resource ID.
	// Returns an error if the success handling fails, which will prevent the Ready condition from being set.
	Success(obj genruntime.ARMMetaObject) error
}

// SuccessFunc is the signature of a function that can be used to create a default SuccessfulCreationHandler
type SuccessFunc = func(obj genruntime.ARMMetaObject) error

// CreateSuccessfulCreationHandler creates a SuccessFunc if the resource implements SuccessfulCreationHandler.
// If the resource did not implement SuccessfulCreationHandler a default handler that does nothing is returned.
func CreateSuccessfulCreationHandler(
	host genruntime.ResourceExtension,
	log logr.Logger,
) SuccessFunc {
	impl, ok := host.(SuccessfulCreationHandler)
	if !ok {
		return func(obj genruntime.ARMMetaObject) error {
			return nil
		}
	}

	return func(obj genruntime.ARMMetaObject) error {
		log.V(Status).Info("Handling successful resource creation")
		err := impl.Success(obj)
		if err != nil {
			return eris.Wrapf(err, "custom resource success handler failed")
		}
		return nil
	}
}
