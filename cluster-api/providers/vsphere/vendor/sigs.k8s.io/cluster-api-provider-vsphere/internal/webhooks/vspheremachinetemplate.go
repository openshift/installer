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

package webhooks

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/cluster-api/util/topology"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
)

const machineTemplateImmutableMsg = "VSphereMachineTemplate spec.template.spec field is immutable. Please create a new resource instead."

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-vspheremachinetemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=vspheremachinetemplates,versions=v1beta1,name=validation.vspheremachinetemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

// VSphereMachineTemplate implements a validation webhook for VSphereMachineTemplate.
type VSphereMachineTemplate struct{}

var _ webhook.CustomValidator = &VSphereMachineTemplate{}

func (webhook *VSphereMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&infrav1.VSphereMachineTemplate{}).
		WithValidator(webhook).
		Complete()
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachineTemplate) ValidateCreate(ctx context.Context, raw runtime.Object) (admission.Warnings, error) {
	obj, ok := raw.(*infrav1.VSphereMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereMachineTemplate but got a %T", raw))
	}

	var allErrs field.ErrorList
	spec := obj.Spec.Template.Spec

	if spec.Network.PreferredAPIServerCIDR != "" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "PreferredAPIServerCIDR"), spec.Network.PreferredAPIServerCIDR, "cannot be set, as it will be removed and is no longer used"))
	}
	if spec.ProviderID != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "providerID"), "cannot be set in templates"))
	}
	for _, device := range spec.Network.Devices {
		if len(device.IPAddrs) != 0 {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "network", "devices", "ipAddrs"), "cannot be set in templates"))
		}
	}
	if spec.HardwareVersion != "" {
		r := regexp.MustCompile("^vmx-[1-9][0-9]?$")
		if !r.MatchString(spec.HardwareVersion) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "template", "spec", "hardwareVersion"), spec.HardwareVersion, "should be a valid VM hardware version, example vmx-17"))
		}
	}
	if spec.GuestSoftPowerOffTimeout != nil {
		if spec.PowerOffMode != infrav1.VirtualMachinePowerOpModeTrySoft {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "template", "spec", "guestSoftPowerOffTimeout"), spec.GuestSoftPowerOffTimeout, "should not be set in templates unless the powerOffMode is trySoft"))
		}
		if spec.GuestSoftPowerOffTimeout.Duration <= 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "template", "spec", "guestSoftPowerOffTimeout"), spec.GuestSoftPowerOffTimeout, "should be greater than 0"))
		}
	}
	pciErrs := validatePCIDevices(spec.PciDevices)
	allErrs = append(allErrs, pciErrs...)

	templateErrs := validateVSphereVMNamingTemplate(ctx, obj)
	if len(templateErrs) > 0 {
		allErrs = append(allErrs, templateErrs...)
	}
	return nil, AggregateObjErrors(obj.GroupVersionKind().GroupKind(), obj.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachineTemplate) ValidateUpdate(ctx context.Context, oldRaw runtime.Object, newRaw runtime.Object) (admission.Warnings, error) {
	oldTyped, ok := oldRaw.(*infrav1.VSphereMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereMachineTemplate but got a %T", oldRaw))
	}
	newTyped, ok := newRaw.(*infrav1.VSphereMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereMachineTemplate but got a %T", newRaw))
	}

	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a admission.Request inside context: %v", err))
	}

	var allErrs field.ErrorList
	if !topology.IsDryRunRequest(req, newTyped) &&
		!reflect.DeepEqual(newTyped.Spec.Template.Spec, oldTyped.Spec.Template.Spec) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "template", "spec"), newTyped, machineTemplateImmutableMsg))
	}

	templateErrs := validateVSphereVMNamingTemplate(ctx, newTyped)
	if len(templateErrs) > 0 {
		allErrs = append(allErrs, templateErrs...)
	}
	return nil, AggregateObjErrors(newTyped.GroupVersionKind().GroupKind(), newTyped.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachineTemplate) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func validateVSphereVMNamingTemplate(_ context.Context, vsphereMachineTemplate *infrav1.VSphereMachineTemplate) field.ErrorList {
	var allErrs field.ErrorList
	namingStrategy := vsphereMachineTemplate.Spec.Template.Spec.NamingStrategy
	if namingStrategy != nil && namingStrategy.Template != nil {
		name, err := services.GenerateVSphereVMName("machine", namingStrategy)
		templateFldPath := field.NewPath("spec", "template", "spec", "namingStrategy", "template")
		if err != nil {
			allErrs = append(allErrs,
				field.Invalid(
					templateFldPath,
					*namingStrategy.Template,
					fmt.Sprintf("invalid VSphereVM name template: %v", err),
				),
			)
		} else {
			// Note: This validates that the resulting name is a valid Kubernetes object name.
			for _, err := range validation.IsDNS1123Subdomain(name) {
				allErrs = append(allErrs,
					field.Invalid(
						templateFldPath,
						*namingStrategy.Template,
						fmt.Sprintf("invalid VSphereVM name template, generated name is not a valid Kubernetes object name: %v", err),
					),
				)
			}
		}
	}
	return allErrs
}
