// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package webhook

import (
	"context"
	"fmt"

	"github.com/rotisserie/eris"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	v1 "github.com/Azure/azure-service-operator/v2/api/sql/v1"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

type User_Webhook struct{}

// +kubebuilder:webhook:path=/mutate-sql-azure-com-v1-user,mutating=true,sideEffects=None,matchPolicy=Exact,failurePolicy=fail,groups=sql.azure.com,resources=users,verbs=create;update,versions=v1,name=default.v1.users.sql.azure.com,admissionReviewVersions=v1

var _ webhook.CustomDefaulter = &User_Webhook{}

func (webhook *User_Webhook) Default(ctx context.Context, obj runtime.Object) error {
	resource, ok := obj.(*v1.User)
	if !ok {
		return fmt.Errorf("expected github.com/Azure/azure-service-operator/v2/api/dbforpostgresql/v1/User, but got %T", obj)
	}
	err := webhook.defaultImpl(ctx, resource)
	if err != nil {
		return err
	}
	var temp any = webhook
	if runtimeDefaulter, ok := temp.(genruntime.Defaulter); ok {
		err = runtimeDefaulter.CustomDefault(ctx, resource)
		if err != nil {
			return err
		}
	}
	return nil
}

// defaultAzureName defaults the Azure name of the resource to the Kubernetes name
func (webhook *User_Webhook) defaultAzureName(_ context.Context, user *v1.User) error {
	if user.Spec.AzureName == "" {
		user.Spec.AzureName = user.Name
	}
	return nil
}

// defaultImpl applies the code generated defaults to the FlexibleServer resource
func (webhook *User_Webhook) defaultImpl(ctx context.Context, user *v1.User) error {
	err := webhook.defaultAzureName(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// +kubebuilder:webhook:path=/validate-sql-azure-com-v1-user,mutating=false,sideEffects=None,matchPolicy=Exact,failurePolicy=fail,groups=sql.azure.com,resources=users,verbs=create;update,versions=v1,name=validate.v1.users.sql.azure.com,admissionReviewVersions=v1

var _ webhook.CustomValidator = &User_Webhook{}

// ValidateCreate validates the creation of the resource
func (webhook *User_Webhook) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	resource, ok := obj.(*v1.User)
	if !ok {
		return nil, fmt.Errorf("expected github.com/Azure/azure-service-operator/v2/api/dbforpostgresql/v1/User, but got %T", obj)
	}
	validations := webhook.createValidations()
	var temp any = webhook
	if runtimeValidator, ok := temp.(genruntime.Validator[*v1.User]); ok {
		validations = append(validations, runtimeValidator.CreateValidations()...)
	}
	return genruntime.ValidateCreate(ctx, resource, validations)
}

// ValidateDelete validates the deletion of the resource
func (webhook *User_Webhook) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	resource, ok := obj.(*v1.User)
	if !ok {
		return nil, fmt.Errorf("expected github.com/Azure/azure-service-operator/v2/api/dbforpostgresql/v1/User, but got %T", obj)
	}
	validations := webhook.deleteValidations()
	var temp any = webhook
	if runtimeValidator, ok := temp.(genruntime.Validator[*v1.User]); ok {
		validations = append(validations, runtimeValidator.DeleteValidations()...)
	}
	return genruntime.ValidateDelete(ctx, resource, validations)
}

// ValidateUpdate validates an update of the resource
func (webhook *User_Webhook) ValidateUpdate(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (admission.Warnings, error) {
	newResource, ok := newObj.(*v1.User)
	if !ok {
		return nil, fmt.Errorf("expected github.com/Azure/azure-service-operator/v2/api/dbforpostgresql/v1/User, but got %T", newObj)
	}
	oldResource, ok := oldObj.(*v1.User)
	if !ok {
		return nil, fmt.Errorf("expected github.com/Azure/azure-service-operator/v2/api/dbforpostgresql/v1/User, but got %T", oldObj)
	}
	validations := webhook.updateValidations()
	var temp any = webhook
	if runtimeValidator, ok := temp.(genruntime.Validator[*v1.User]); ok {
		validations = append(validations, runtimeValidator.UpdateValidations()...)
	}
	return genruntime.ValidateUpdate(
		ctx,
		oldResource,
		newResource,
		validations)
}

// createValidations validates the creation of the resource
func (webhook *User_Webhook) createValidations() []func(ctx context.Context, obj *v1.User) (admission.Warnings, error) {
	return nil
}

// deleteValidations validates the deletion of the resource
func (webhook *User_Webhook) deleteValidations() []func(ctx context.Context, obj *v1.User) (admission.Warnings, error) {
	return nil
}

// updateValidations validates the update of the resource
func (webhook *User_Webhook) updateValidations() []func(ctx context.Context, oldObj *v1.User, newObj *v1.User) (admission.Warnings, error) {
	return []func(ctx context.Context, oldObj *v1.User, newObj *v1.User) (admission.Warnings, error){
		webhook.validateWriteOncePropertiesNotChanged,
	}
}

// validateWriteOncePropertiesNotChanged function validates the update on WriteOnce properties.
// TODO: Note this should be kept in sync with admissions.ValidateWriteOnceProperties
func (webhook *User_Webhook) validateWriteOncePropertiesNotChanged(_ context.Context, oldObj *v1.User, newObj *v1.User) (admission.Warnings, error) {
	var errs []error

	// If we don't have a finalizer yet, it's OK to change things
	hasFinalizer := controllerutil.ContainsFinalizer(oldObj, genruntime.ReconcilerFinalizer)
	if !hasFinalizer {
		return nil, nil
	}

	if oldObj.Spec.AzureName != newObj.Spec.AzureName {
		err := eris.Errorf(
			"updating 'AzureName' is not allowed for '%s : %s",
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

	if (bothHaveOwner && oldOwner.Name != newOwner.Name) || ownerAdded {
		err := eris.Errorf(
			"updating 'Owner.Name' is not allowed for '%s : %s",
			oldObj.GetObjectKind().GroupVersionKind(),
			oldObj.GetName())

		errs = append(errs, err)
	} else if ownerRemoved {
		err := eris.Errorf(
			"removing 'Owner' is not allowed for '%s : %s",
			oldObj.GetObjectKind().GroupVersionKind(),
			oldObj.GetName())

		errs = append(errs, err)
	}

	return nil, kerrors.NewAggregate(errs)
}
