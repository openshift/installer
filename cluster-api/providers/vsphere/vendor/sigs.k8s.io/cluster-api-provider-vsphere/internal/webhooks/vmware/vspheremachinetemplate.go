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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/vmoperator"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-vmware-infrastructure-cluster-x-k8s-io-v1beta1-vspheremachinetemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=vmware.infrastructure.cluster.x-k8s.io,resources=vspheremachinetemplates,versions=v1beta1,name=validation.vspheremachinetemplate.vmware.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

// VSphereMachineTemplate implements a validation webhook for VSphereMachineTemplate.
type VSphereMachineTemplate struct{}

var _ webhook.CustomValidator = &VSphereMachineTemplate{}

func (webhook *VSphereMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&vmwarev1.VSphereMachineTemplate{}).
		WithValidator(webhook).
		Complete()
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachineTemplate) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	vSphereMachineTemplate, ok := obj.(*vmwarev1.VSphereMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereMachineTemplate but got a %T", obj))
	}
	return webhook.validate(ctx, nil, vSphereMachineTemplate)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachineTemplate) ValidateUpdate(ctx context.Context, _, newRaw runtime.Object) (admission.Warnings, error) {
	vSphereMachineTemplate, ok := newRaw.(*vmwarev1.VSphereMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereMachineTemplate but got a %T", newRaw))
	}
	return webhook.validate(ctx, nil, vSphereMachineTemplate)
}

func (webhook *VSphereMachineTemplate) validate(_ context.Context, _, newVSphereMachineTemplate *vmwarev1.VSphereMachineTemplate) (admission.Warnings, error) {
	var allErrs field.ErrorList

	// Validate namingStrategy
	namingStrategy := newVSphereMachineTemplate.Spec.Template.Spec.NamingStrategy
	if namingStrategy != nil &&
		namingStrategy.Template != nil {
		name, err := vmoperator.GenerateVirtualMachineName("machine", namingStrategy)
		templateFldPath := field.NewPath("spec", "template", "spec", "namingStrategy", "template")
		if err != nil {
			allErrs = append(allErrs,
				field.Invalid(
					templateFldPath,
					*namingStrategy.Template,
					fmt.Sprintf("invalid VirtualMachine name template: %v", err),
				),
			)
		} else {
			// Note: This validates that the resulting name is a valid Kubernetes object name.
			for _, err := range validation.IsDNS1123Subdomain(name) {
				allErrs = append(allErrs,
					field.Invalid(
						templateFldPath,
						*namingStrategy.Template,
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
func (webhook *VSphereMachineTemplate) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
