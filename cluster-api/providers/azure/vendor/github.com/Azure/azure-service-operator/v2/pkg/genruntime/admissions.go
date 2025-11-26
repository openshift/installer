/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"context"

	"github.com/rotisserie/eris"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Validator is similar to controller-runtime/pkg/webhook/admission Validator. Implementing this interface
// allows you to hook into the code generated validations and add custom handcrafted validations.
type Validator[T runtime.Object] interface {
	// CreateValidations returns validation functions that should be run on create.
	CreateValidations() []func(ctx context.Context, obj T) (admission.Warnings, error)
	// UpdateValidations returns validation functions that should be run on update.
	UpdateValidations() []func(ctx context.Context, oldObj T, newObj T) (admission.Warnings, error)
	// DeleteValidations returns validation functions that should be run on delete.
	DeleteValidations() []func(ctx context.Context, obj T) (admission.Warnings, error)
}

// Defaulter is similar to controller-runtime/pkg/webhook/admission Defaulter. Implementing this interface
// allows you to hook into the code generated defaults and add custom handcrafted defaults.
type Defaulter interface {
	// CustomDefault performs custom defaults that are run in addition to the code generated defaults.
	CustomDefault(ctx context.Context, obj runtime.Object) error
}

// ValidateWriteOnceProperties function validates the update on WriteOnce properties.
func ValidateWriteOnceProperties(oldObj ARMMetaObject, newObj ARMMetaObject) (admission.Warnings, error) {
	var errs []error

	if !IsResourceCreatedSuccessfully(newObj) {
		return nil, nil
	}

	// Prohibit changing the AzureName,
	// but allow it to be set if it's empty.
	//
	// https://github.com/Azure/azure-service-operator/issues/4306
	oldName := oldObj.AzureName()
	if oldName != "" && oldName != newObj.AzureName() {
		err := eris.Errorf(
			"updating 'spec.azureName' is not allowed for '%s : %s",
			oldObj.GetObjectKind().GroupVersionKind(),
			oldObj.GetName())
		errs = append(errs, err)
	}

	// Ensure that owner has not been changed
	oldOwner := oldObj.Owner()
	newOwner := newObj.Owner()

	bothHaveOwner := oldOwner != nil && newOwner != nil
	ownerAdded := oldOwner == nil && newOwner != nil
	ownerRemoved := oldOwner != nil && newOwner == nil

	ownerNameChanged := bothHaveOwner && oldOwner.Name != newOwner.Name
	ownerARMIDChanged := bothHaveOwner && oldOwner.ARMID != newOwner.ARMID

	if ownerAdded {
		// This error may not be possible to trigger in practice, as it requires an Azure resource that supports existing without an owner
		// or with an owner. There aren't any resources that meet those criteria that we know of, so this check is primarily us being
		// defensive.
		errs = append(errs, eris.Errorf("adding an owner to an already created resource is not allowed for '%s : %s", oldObj.GetObjectKind().GroupVersionKind(), oldObj.GetName()))
	} else if ownerNameChanged {
		errs = append(errs, eris.Errorf("updating 'spec.owner.name' is not allowed for '%s : %s", oldObj.GetObjectKind().GroupVersionKind(), oldObj.GetName()))
	} else if ownerARMIDChanged {
		errs = append(errs, eris.Errorf("updating 'spec.owner.armId' is not allowed for '%s : %s", oldObj.GetObjectKind().GroupVersionKind(), oldObj.GetName()))
	} else if ownerRemoved {
		errs = append(errs, eris.Errorf("removing 'spec.owner' is not allowed for '%s : %s", oldObj.GetObjectKind().GroupVersionKind(), oldObj.GetName()))
	}

	return nil, kerrors.NewAggregate(errs)
}

func ValidateCreate[T runtime.Object](ctx context.Context, resource T, validations []func(ctx context.Context, resource T) (admission.Warnings, error)) (admission.Warnings, error) {
	var errs []error
	var warnings admission.Warnings
	for _, validation := range validations {
		warning, err := validation(ctx, resource)
		if warning != nil {
			warnings = append(warnings, warning...)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return warnings, kerrors.NewAggregate(errs)
}

func ValidateDelete[T runtime.Object](ctx context.Context, resource T, validations []func(ctx context.Context, resource T) (admission.Warnings, error)) (admission.Warnings, error) {
	var errs []error
	var warnings admission.Warnings
	for _, validation := range validations {
		warning, err := validation(ctx, resource)
		if warning != nil {
			warnings = append(warnings, warning...)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return warnings, kerrors.NewAggregate(errs)
}

func ValidateUpdate[T runtime.Object](ctx context.Context, old T, new T, validations []func(ctx context.Context, old T, new T) (admission.Warnings, error)) (admission.Warnings, error) {
	var errs []error
	var warnings admission.Warnings
	for _, validation := range validations {
		warning, err := validation(ctx, old, new)
		if warning != nil {
			warnings = append(warnings, warning...)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return warnings, kerrors.NewAggregate(errs)
}

func IsResourceCreatedSuccessfully(obj ARMMetaObject) bool {
	return GetResourceIDOrDefault(obj) != ""
}
