/*
Copyright 2024 The Kubernetes Authors.

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

// Package vmware is the package for webhooks of vmware resources.
package vmware

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/cluster-api/util/topology"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta2"
	"sigs.k8s.io/cluster-api-provider-vsphere/internal/webhooks/vmware/conversion"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/vmoperator"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-vmware-infrastructure-cluster-x-k8s-io-v1beta2-vspheremachinetemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=vmware.infrastructure.cluster.x-k8s.io,resources=vspheremachinetemplates,versions=v1beta2,name=validation.vspheremachinetemplate.vmware.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-vmware-infrastructure-cluster-x-k8s-io-v1beta2-vspheremachinetemplate,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=vmware.infrastructure.cluster.x-k8s.io,resources=vspheremachinetemplates,versions=v1beta2,name=default.vspheremachinetemplate.vmware.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1

// VSphereMachineTemplate implements a validation webhook for VSphereMachineTemplate.
type VSphereMachineTemplate struct {
	// NetworkProvider is the network provider used by Supervisor based clusters
	NetworkProvider string
}

var _ admission.Defaulter[*vmwarev1.VSphereMachineTemplate] = &VSphereMachineTemplate{}
var _ admission.Validator[*vmwarev1.VSphereMachineTemplate] = &VSphereMachineTemplate{}

func (webhook *VSphereMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &vmwarev1.VSphereMachineTemplate{}).
		WithDefaulter(webhook).
		WithValidator(webhook).
		WithConverter(conversion.VSphereMachineTemplate).
		Complete()
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (webhook *VSphereMachineTemplate) Default(ctx context.Context, c *vmwarev1.VSphereMachineTemplate) error {
	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return apierrors.NewBadRequest(fmt.Sprintf("expected an admission.Request inside context: %v", err))
	}

	if topology.IsDryRunRequest(req, c) {
		// In case of dry-run requests from the topology controller, apply defaults from older versions of CAPV
		// so we do not trigger rollouts when dealing with objects created before dropping those defaults.
		applyPreviousVSphereMachineTemplateDefaults(c)
	}

	return nil
}

func applyPreviousVSphereMachineTemplateDefaults(c *vmwarev1.VSphereMachineTemplate) {
	if c.Spec.Template.Spec.PowerOffMode == "" {
		c.Spec.Template.Spec.PowerOffMode = vmwarev1.VirtualMachinePowerOpModeHard
	}
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachineTemplate) ValidateCreate(ctx context.Context, vSphereMachineTemplate *vmwarev1.VSphereMachineTemplate) (admission.Warnings, error) {
	return webhook.validate(ctx, nil, vSphereMachineTemplate)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachineTemplate) ValidateUpdate(ctx context.Context, oldObj, newObj *vmwarev1.VSphereMachineTemplate) (admission.Warnings, error) {
	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a admission.Request inside context: %v", err))
	}
	if !topology.IsDryRunRequest(req, newObj) {
		// Apply defaults from older versions of CAPV so the following checks do not report differences when
		// dealing with objects created before dropping those defaults.
		applyPreviousVSphereMachineTemplateDefaults(oldObj)
		applyPreviousVSphereMachineTemplateDefaults(newObj)

		equal, diff, err := util.Diff(oldObj.Spec.Template.Spec, newObj.Spec.Template.Spec)
		if err != nil {
			return nil, apierrors.NewBadRequest(fmt.Sprintf("failed to compare old and new VSphereMachineTemplate: %v", err))
		}
		if !equal {
			return nil, apierrors.NewInvalid(vmwarev1.GroupVersion.WithKind("VSphereMachineTemplate").GroupKind(), newObj.Name, field.ErrorList{
				field.Invalid(field.NewPath("spec", "template", "spec"), newObj, fmt.Sprintf("VSphereMachineTemplate spec.template.spec field is immutable. Please create a new resource instead. Diff: %s", diff)),
			})
		}
	}

	return webhook.validate(ctx, nil, newObj)
}

func (webhook *VSphereMachineTemplate) validate(_ context.Context, _, newVSphereMachineTemplate *vmwarev1.VSphereMachineTemplate) (admission.Warnings, error) {
	allErrs := validateNetwork(webhook.NetworkProvider, newVSphereMachineTemplate.Spec.Template.Spec.Network, field.NewPath("spec", "template", "spec", "network"))
	allErrs = append(allErrs, validatePolicies(newVSphereMachineTemplate.Spec.Template.Spec.Policies, field.NewPath("spec", "template", "spec", "policies"))...)

	// Validate namingStrategy
	namingStrategy := newVSphereMachineTemplate.Spec.Template.Spec.Naming
	if namingStrategy.Template != "" {
		name, err := vmoperator.GenerateVirtualMachineName("machine", namingStrategy)
		templateFldPath := field.NewPath("spec", "template", "spec", "namingStrategy", "template")
		if err != nil {
			allErrs = append(allErrs,
				field.Invalid(
					templateFldPath,
					namingStrategy.Template,
					fmt.Sprintf("invalid VirtualMachine name template: %v", err),
				),
			)
		} else {
			// Note: This validates that the resulting name is a valid Kubernetes object name.
			for _, err := range validation.IsDNS1123Subdomain(name) {
				allErrs = append(allErrs,
					field.Invalid(
						templateFldPath,
						namingStrategy.Template,
						fmt.Sprintf("invalid VirtualMachine name template, generated name is not a valid Kubernetes object name: %v", err),
					),
				)
			}
		}
	}

	if len(allErrs) > 0 {
		return nil, apierrors.NewInvalid(vmwarev1.GroupVersion.WithKind("VSphereMachineTemplate").GroupKind(), newVSphereMachineTemplate.Name, allErrs)
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachineTemplate) ValidateDelete(_ context.Context, _ *vmwarev1.VSphereMachineTemplate) (admission.Warnings, error) {
	return nil, nil
}
