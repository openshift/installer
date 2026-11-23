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
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/cluster-api/util/topology"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
	"sigs.k8s.io/cluster-api-provider-vsphere/internal/webhooks/conversion"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
)

const machineTemplateImmutableMsg = "VSphereMachineTemplate spec.template.spec field is immutable. Please create a new resource instead."

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-vspheremachinetemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=vspheremachinetemplates,versions=v1beta2,name=validation.vspheremachinetemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-vspheremachinetemplate,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=vspheremachinetemplates,versions=v1beta2,name=default.vspheremachinetemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1

// VSphereMachineTemplate implements a validation webhook for VSphereMachineTemplate.
type VSphereMachineTemplate struct{}

var _ admission.Defaulter[*infrav1.VSphereMachineTemplate] = &VSphereMachineTemplate{}
var _ admission.Validator[*infrav1.VSphereMachineTemplate] = &VSphereMachineTemplate{}

func (webhook *VSphereMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.VSphereMachineTemplate{}).
		WithDefaulter(webhook).
		WithValidator(webhook).
		WithConverter(conversion.VSphereMachineTemplate).
		Complete()
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (webhook *VSphereMachineTemplate) Default(ctx context.Context, c *infrav1.VSphereMachineTemplate) error {
	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return apierrors.NewBadRequest(fmt.Sprintf("expected an admission.Request inside context: %v", err))
	}

	if topology.IsDryRunRequest(req, c) {
		// In case of dry-run requests from the topology controller, apply defaults from older versions of CAPV
		// so we do not trigger rollouts when dealing with objects created before dropping those defaults.
		if c.Spec.Template.Spec.PowerOffMode == "" {
			c.Spec.Template.Spec.PowerOffMode = infrav1.VirtualMachinePowerOpModeHard
		}
	}

	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachineTemplate) ValidateCreate(ctx context.Context, obj *infrav1.VSphereMachineTemplate) (admission.Warnings, error) {
	var allErrs field.ErrorList
	spec := obj.Spec.Template.Spec

	if spec.ProviderID != "" {
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
	if spec.GuestSoftPowerOffTimeoutSeconds != 0 {
		if spec.PowerOffMode != infrav1.VirtualMachinePowerOpModeTrySoft {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "template", "spec", "guestSoftPowerOffTimeout"), spec.GuestSoftPowerOffTimeoutSeconds, "should not be set in templates unless the powerOffMode is trySoft"))
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
func (webhook *VSphereMachineTemplate) ValidateUpdate(ctx context.Context, oldTyped, newTyped *infrav1.VSphereMachineTemplate) (admission.Warnings, error) {
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
func (webhook *VSphereMachineTemplate) ValidateDelete(_ context.Context, _ *infrav1.VSphereMachineTemplate) (admission.Warnings, error) {
	return nil, nil
}

func validateVSphereVMNamingTemplate(_ context.Context, vsphereMachineTemplate *infrav1.VSphereMachineTemplate) field.ErrorList {
	var allErrs field.ErrorList
	namingStrategy := vsphereMachineTemplate.Spec.Template.Spec.Naming
	if namingStrategy.Template != "" {
		name, err := services.GenerateVSphereVMName("machine", namingStrategy)
		templateFldPath := field.NewPath("spec", "template", "spec", "namingStrategy", "template")
		if err != nil {
			allErrs = append(allErrs,
				field.Invalid(
					templateFldPath,
					namingStrategy.Template,
					fmt.Sprintf("invalid VSphereVM name template: %v", err),
				),
			)
		} else {
			// Note: This validates that the resulting name is a valid Kubernetes object name.
			for _, err := range validation.IsDNS1123Subdomain(name) {
				allErrs = append(allErrs,
					field.Invalid(
						templateFldPath,
						namingStrategy.Template,
						fmt.Sprintf("invalid VSphereVM name template, generated name is not a valid Kubernetes object name: %v", err),
					),
				)
			}
		}
	}
	return allErrs
}
