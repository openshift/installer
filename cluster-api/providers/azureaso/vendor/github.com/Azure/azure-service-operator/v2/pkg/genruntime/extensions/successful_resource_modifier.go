/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package extensions

import (
	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// SuccessfulCreationHandler can be implemented to customize the resource upon successful creation
type SuccessfulCreationHandler interface {
	// Success modifies the resource based on a successful creation
	Success(obj genruntime.ARMMetaObject) error
}

// SuccessFunc is the signature of a function that can be used to create a default SuccessfulCreationHandler
type SuccessFunc = func(obj genruntime.ARMMetaObject) error

// CreateSuccessfulCreationHandler creates a SuccessFunc if the resource implements SuccessfulCreationHandler.
// If the resource did not implement SuccessfulCreationHandler a default handler that does nothing is returned.
func CreateSuccessfulCreationHandler(
	host genruntime.ResourceExtension,
	log logr.Logger) SuccessFunc {

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
			return errors.Wrapf(err, "custom resource success handler failed")
		}
		return nil
	}
}
