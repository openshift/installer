/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"context"
	"fmt"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api/util/topology"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// AzureMachineTemplateImmutableMsg ...
const (
	AzureMachineTemplateImmutableMsg                      = "AzureMachineTemplate spec.template.spec field is immutable. Please create new resource instead. ref doc: https://cluster-api.sigs.k8s.io/tasks/updating-machine-templates.html"
	AzureMachineTemplateRoleAssignmentNameMsg             = "AzureMachineTemplate spec.template.spec.roleAssignmentName field can't be set"
	AzureMachineTemplateSystemAssignedIdentityRoleNameMsg = "AzureMachineTemplate spec.template.spec.systemAssignedIdentityRole.name field can't be set"
)

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (r *AzureMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(azureMachineTemplateWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-azuremachinetemplate,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azuremachinetemplates,versions=v1beta1,name=default.azuremachinetemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azuremachinetemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azuremachinetemplates,versions=v1beta1,name=validation.azuremachinetemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type azureMachineTemplateWebhook struct{}

var _ webhook.CustomDefaulter = &azureMachineTemplateWebhook{}
var _ webhook.CustomValidator = &azureMachineTemplateWebhook{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*azureMachineTemplateWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r := obj.(*AzureMachineTemplate)
	spec := r.Spec.Template.Spec

	allErrs := ValidateAzureMachineSpec(spec)

	if spec.RoleAssignmentName != "" {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("AzureMachineTemplate", "spec", "template", "spec", "roleAssignmentName"), r, AzureMachineTemplateRoleAssignmentNameMsg),
		)
	}

	if spec.SystemAssignedIdentityRole != nil && spec.SystemAssignedIdentityRole.Name != "" {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("AzureMachineTemplate", "spec", "template", "spec", "systemAssignedIdentityRole", "name"), r, AzureMachineTemplateSystemAssignedIdentityRoleNameMsg),
		)
	}

	if (r.Spec.Template.Spec.NetworkInterfaces != nil) && len(r.Spec.Template.Spec.NetworkInterfaces) > 0 && r.Spec.Template.Spec.SubnetName != "" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("AzureMachineTemplate", "spec", "template", "spec", "networkInterfaces"), r.Spec.Template.Spec.NetworkInterfaces, "cannot set both NetworkInterfaces and machine SubnetName"))
	}

	if (r.Spec.Template.Spec.NetworkInterfaces != nil) && len(r.Spec.Template.Spec.NetworkInterfaces) > 0 && r.Spec.Template.Spec.AcceleratedNetworking != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("AzureMachineTemplate", "spec", "template", "spec", "acceleratedNetworking"), r.Spec.Template.Spec.NetworkInterfaces, "cannot set both NetworkInterfaces and machine AcceleratedNetworking"))
	}

	for i, networkInterface := range r.Spec.Template.Spec.NetworkInterfaces {
		if networkInterface.PrivateIPConfigs < 1 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("AzureMachineTemplate", "spec", "template", "spec", "networkInterfaces", "privateIPConfigs"), r.Spec.Template.Spec.NetworkInterfaces[i].PrivateIPConfigs, "networkInterface privateIPConfigs must be set to a minimum value of 1"))
		}
	}

	if ptr.Deref(r.Spec.Template.Spec.DisableExtensionOperations, false) && len(r.Spec.Template.Spec.VMExtensions) > 0 {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("AzureMachineTemplate", "spec", "template", "spec", "vmExtensions"), "VMExtensions must be empty when DisableExtensionOperations is true"))
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(GroupVersion.WithKind(AzureMachineTemplateKind).GroupKind(), r.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *azureMachineTemplateWebhook) ValidateUpdate(ctx context.Context, oldRaw runtime.Object, newRaw runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList
	old := oldRaw.(*AzureMachineTemplate)
	t := newRaw.(*AzureMachineTemplate)

	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a admission.Request inside context: %v", err))
	}

	if !topology.ShouldSkipImmutabilityChecks(req, t) &&
		!reflect.DeepEqual(t.Spec.Template.Spec, old.Spec.Template.Spec) {
		// The equality failure could be because of default mismatch between v1alpha3 and v1beta1. This happens because
		// the new object `r` will have run through the default webhooks but the old object `old` would not have so.
		// This means if the old object was in v1alpha3, it would not get the new defaults set in v1beta1 resulting
		// in object inequality. To workaround this, we set the v1beta1 defaults here so that the old object also gets
		// the new defaults.

		// We need to set ssh key explicitly, otherwise Default() will create a new one.
		if old.Spec.Template.Spec.SSHPublicKey == "" {
			old.Spec.Template.Spec.SSHPublicKey = t.Spec.Template.Spec.SSHPublicKey
		}

		if err := r.Default(ctx, old); err != nil {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("AzureMachineTemplate"), t, fmt.Sprintf("Unable to apply defaults: %v", err)),
			)
		}

		// if it's still not equal, return error.
		if !reflect.DeepEqual(t.Spec.Template.Spec, old.Spec.Template.Spec) {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("AzureMachineTemplate", "spec", "template", "spec"), t, AzureMachineTemplateImmutableMsg),
			)
		}
	}

	if len(allErrs) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(GroupVersion.WithKind(AzureMachineTemplateKind).GroupKind(), t.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (*azureMachineTemplateWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// Default implements webhookutil.defaulter so a webhook will be registered for the type.
func (*azureMachineTemplateWebhook) Default(_ context.Context, obj runtime.Object) error {
	t := obj.(*AzureMachineTemplate)
	if err := t.Spec.Template.Spec.SetDefaultSSHPublicKey(); err != nil {
		ctrl.Log.WithName("SetDefault").Error(err, "SetDefaultSSHPublicKey failed")
	}
	t.Spec.Template.Spec.SetDataDisksDefaults()
	t.Spec.Template.Spec.SetNetworkInterfacesDefaults()
	return nil
}
