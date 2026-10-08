// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package webhook

import (
	"context"

	"github.com/rotisserie/eris"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	v1 "github.com/Azure/azure-service-operator/v2/api/entra/v1"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

type SecurityGroup_Webhook struct{}

// +kubebuilder:webhook:path=/mutate-entra-azure-com-v1-securitygroup,mutating=true,sideEffects=None,matchPolicy=Exact,failurePolicy=fail,groups=entra.azure.com,resources=securitygroups,verbs=create;update,versions=v1,name=default.v1.securitygroups.entra.azure.com,admissionReviewVersions=v1

var _ webhook.CustomDefaulter = &SecurityGroup_Webhook{}

func (webhook *SecurityGroup_Webhook) Default(
	ctx context.Context,
	obj runtime.Object,
) error {
	resource, err := webhook.asSecurityGroup(obj)
	if err != nil {
		return eris.Wrapf(err, "setting defaults for resource")
	}

	webhook.defaultImpl(ctx, resource)

	return nil
}

func (*SecurityGroup_Webhook) asSecurityGroup(
	obj runtime.Object,
) (*v1.SecurityGroup, error) {
	typedObj, ok := obj.(*v1.SecurityGroup)
	if !ok {
		return nil, eris.Errorf("cannot modify resource that is not of type *v1.SecurityGroup. Type is %T", obj)
	}

	return typedObj, nil
}

// defaultImpl applies defaults to the SecurityGroup resource
func (webhook *SecurityGroup_Webhook) defaultImpl(
	_ context.Context,
	securityGroup *v1.SecurityGroup,
) {
	// Ensure we always have an OperatorSpec
	if securityGroup.Spec.OperatorSpec == nil {
		securityGroup.Spec.OperatorSpec = &v1.SecurityGroupOperatorSpec{}
	}

	// If CreationMode not specified, default to AdoptOrCreate
	if securityGroup.Spec.OperatorSpec.CreationMode == nil {
		securityGroup.Spec.OperatorSpec.CreationMode = to.Ptr(v1.AdoptOrCreate)
	}
}

// +kubebuilder:webhook:path=/validate-entra-azure-com-v1-securitygroup,mutating=false,sideEffects=None,matchPolicy=Exact,failurePolicy=fail,groups=entra.azure.com,resources=securitygroups,verbs=create;update,versions=v1,name=validate.v1.securitygroups.entra.azure.com,admissionReviewVersions=v1

var _ webhook.CustomValidator = &SecurityGroup_Webhook{}

// ValidateCreate implements admission.CustomValidator.
func (webhook *SecurityGroup_Webhook) ValidateCreate(
	ctx context.Context,
	obj runtime.Object,
) (warnings admission.Warnings, err error) {
	resource, err := webhook.asSecurityGroup(obj)
	if err != nil {
		return nil, eris.Wrapf(err, "validating creation of resource")
	}

	validations := webhook.createValidations()
	var temp any = webhook
	if runtimeValidator, ok := temp.(genruntime.Validator[*v1.SecurityGroup]); ok {
		validations = append(validations, runtimeValidator.CreateValidations()...)
	}

	return genruntime.ValidateCreate(ctx, resource, validations)
}

// ValidateDelete implements admission.CustomValidator.
func (webhook *SecurityGroup_Webhook) ValidateDelete(
	ctx context.Context,
	obj runtime.Object,
) (warnings admission.Warnings, err error) {
	resource, err := webhook.asSecurityGroup(obj)
	if err != nil {
		return nil, eris.Wrapf(err, "validating deletion of resource")
	}

	validations := webhook.deleteValidations()
	var temp any = webhook
	if runtimeValidator, ok := temp.(genruntime.Validator[*v1.SecurityGroup]); ok {
		validations = append(validations, runtimeValidator.DeleteValidations()...)
	}

	return genruntime.ValidateDelete(ctx, resource, validations)
}

// ValidateUpdate implements admission.CustomValidator.
func (webhook *SecurityGroup_Webhook) ValidateUpdate(
	ctx context.Context,
	oldObj runtime.Object,
	newObj runtime.Object,
) (warnings admission.Warnings, err error) {
	newResource, err := webhook.asSecurityGroup(newObj)
	if err != nil {
		return nil, eris.Wrapf(err, "validating update of resource")
	}

	oldResource, err := webhook.asSecurityGroup(oldObj)
	if err != nil {
		return nil, eris.Wrapf(err, "validating update of resource")
	}

	validations := webhook.updateValidations()
	var temp any = webhook
	if runtimeValidator, ok := temp.(genruntime.Validator[*v1.SecurityGroup]); ok {
		validations = append(validations, runtimeValidator.UpdateValidations()...)
	}

	return genruntime.ValidateUpdate(ctx, oldResource, newResource, validations)
}

// createValidations validates the creation of the resource
func (webhook *SecurityGroup_Webhook) createValidations() []func(ctx context.Context, obj *v1.SecurityGroup) (admission.Warnings, error) {
	return nil
}

// deleteValidations validates the deletion of the resource
func (webhook *SecurityGroup_Webhook) deleteValidations() []func(ctx context.Context, obj *v1.SecurityGroup) (admission.Warnings, error) {
	return nil
}

// updateValidations validates the update of the resource
func (webhook *SecurityGroup_Webhook) updateValidations() []func(ctx context.Context, oldObj *v1.SecurityGroup, newObj *v1.SecurityGroup) (admission.Warnings, error) {
	return nil
}
