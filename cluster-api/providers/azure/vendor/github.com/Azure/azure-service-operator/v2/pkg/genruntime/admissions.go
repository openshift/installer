/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Validator is similar to controller-runtime/pkg/webhook/admission Validator. Implementing this interface
// allows you to hook into the code generated validations and add custom handcrafted validations.
type Validator interface {
	// CreateValidations returns validation functions that should be run on create.
	CreateValidations() []func() (admission.Warnings, error)
	// UpdateValidations returns validation functions that should be run on update.
	UpdateValidations() []func(old runtime.Object) (admission.Warnings, error)
	// DeleteValidations returns validation functions that should be run on delete.
	DeleteValidations() []func() (admission.Warnings, error)
}

// Defaulter is similar to controller-runtime/pkg/webhook/admission Defaulter. Implementing this interface
// allows you to hook into the code generated defaults and add custom handcrafted defaults.
type Defaulter interface {
	// CustomDefault performs custom defaults that are run in addition to the code generated defaults.
	CustomDefault()
}

// ValidateWriteOnceProperties function validates the update on WriteOnce properties.
func ValidateWriteOnceProperties(oldObj ARMMetaObject, newObj ARMMetaObject) (admission.Warnings, error) {
	var errs []error

	if !IsResourceCreatedSuccessfully(newObj) {
		return nil, nil
	}

	if oldObj.AzureName() != newObj.AzureName() {
		errs = append(errs, errors.Errorf("updating 'AzureName' is not allowed for '%s : %s", oldObj.GetObjectKind().GroupVersionKind(), oldObj.GetName()))
	}

	// Ensure that owner has not been changed
	oldOwner := oldObj.Owner()
	newOwner := newObj.Owner()

	bothHaveOwner := oldOwner != nil && newOwner != nil
	ownerAdded := oldOwner == nil && newOwner != nil
	ownerRemoved := oldOwner != nil && newOwner == nil

	if (bothHaveOwner && oldOwner.Name != newOwner.Name) || ownerAdded {
		errs = append(errs, errors.Errorf("updating 'Owner.Name' is not allowed for '%s : %s", oldObj.GetObjectKind().GroupVersionKind(), oldObj.GetName()))
	} else if ownerRemoved {
		errs = append(errs, errors.Errorf("removing 'Owner' is not allowed for '%s : %s", oldObj.GetObjectKind().GroupVersionKind(), oldObj.GetName()))
	}

	return nil, kerrors.NewAggregate(errs)
}

func ValidateCreate(validations []func() (admission.Warnings, error)) (admission.Warnings, error) {
	var errs []error
	var warnings admission.Warnings
	for _, validation := range validations {
		warning, err := validation()
		if warning != nil {
			warnings = append(warnings, warning...)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return warnings, kerrors.NewAggregate(errs)
}

func ValidateDelete(validations []func() (admission.Warnings, error)) (admission.Warnings, error) {
	var errs []error
	var warnings admission.Warnings
	for _, validation := range validations {
		warning, err := validation()
		if warning != nil {
			warnings = append(warnings, warning...)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return warnings, kerrors.NewAggregate(errs)
}

func ValidateUpdate(old runtime.Object, validations []func(old runtime.Object) (admission.Warnings, error)) (admission.Warnings, error) {
	var errs []error
	var warnings admission.Warnings
	for _, validation := range validations {
		warning, err := validation(old)
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
