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
	"net"
	"reflect"

	pkgerrors "github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
	"sigs.k8s.io/cluster-api-provider-vsphere/internal/webhooks/conversion"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-vspheremachine,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=vspheremachines,versions=v1beta2,name=validation.vspheremachine.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-vspheremachine,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=vspheremachines,versions=v1beta2,name=default.vspheremachine.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1

// VSphereMachine implements a validation and defaulting webhook for VSphereMachine.
type VSphereMachine struct{}

var _ admission.Validator[*infrav1.VSphereMachine] = &VSphereMachine{}
var _ admission.Defaulter[*infrav1.VSphereMachine] = &VSphereMachine{}

func (webhook *VSphereMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.VSphereMachine{}).
		WithValidator(webhook).
		WithDefaulter(webhook).
		WithConverter(conversion.VSphereMachine).
		Complete()
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (webhook *VSphereMachine) Default(_ context.Context, objValue *infrav1.VSphereMachine) error {
	if objValue.Spec.Datacenter == "" {
		objValue.Spec.Datacenter = "*"
	}
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachine) ValidateCreate(_ context.Context, obj *infrav1.VSphereMachine) (admission.Warnings, error) {
	var allErrs field.ErrorList

	spec := obj.Spec

	for i, device := range spec.Network.Devices {
		for j, ip := range device.IPAddrs {
			if _, _, err := net.ParseCIDR(ip); err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "network", fmt.Sprintf("devices[%d]", i), fmt.Sprintf("ipAddrs[%d]", j)), ip, "ip addresses should be in the CIDR format"))
			}
		}
	}

	if spec.GuestSoftPowerOffTimeoutSeconds != 0 {
		if spec.PowerOffMode != infrav1.VirtualMachinePowerOpModeTrySoft {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "guestSoftPowerOffTimeout"), spec.GuestSoftPowerOffTimeoutSeconds, "should not be set in templates unless the powerOffMode is trySoft"))
		}
	}
	pciErrs := validatePCIDevices(spec.PciDevices)
	allErrs = append(allErrs, pciErrs...)

	return nil, AggregateObjErrors(obj.GroupVersionKind().GroupKind(), obj.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachine) ValidateUpdate(_ context.Context, oldTyped, newTyped *infrav1.VSphereMachine) (admission.Warnings, error) {
	var allErrs field.ErrorList

	if newTyped.Spec.GuestSoftPowerOffTimeoutSeconds != 0 {
		if newTyped.Spec.PowerOffMode != infrav1.VirtualMachinePowerOpModeTrySoft {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "guestSoftPowerOffTimeout"), newTyped.Spec.GuestSoftPowerOffTimeoutSeconds, "should not be set in templates unless the powerOffMode is trySoft"))
		}
	}

	newVSphereMachine, err := runtime.DefaultUnstructuredConverter.ToUnstructured(newTyped)
	if err != nil {
		return nil, apierrors.NewInternalError(pkgerrors.Wrap(err, "failed to convert new VSphereMachine to unstructured object"))
	}

	oldVSphereMachine, err := runtime.DefaultUnstructuredConverter.ToUnstructured(oldTyped)
	if err != nil {
		return nil, apierrors.NewInternalError(pkgerrors.Wrap(err, "failed to convert old VSphereMachine to unstructured object"))
	}

	newVSphereMachineSpec := newVSphereMachine["spec"].(map[string]interface{})
	oldVSphereMachineSpec := oldVSphereMachine["spec"].(map[string]interface{})

	allowChangeKeys := []string{"providerID", "powerOffMode", "guestSoftPowerOffTimeoutSeconds"}
	for _, key := range allowChangeKeys {
		delete(oldVSphereMachineSpec, key)
		delete(newVSphereMachineSpec, key)
	}

	newVSphereMachineNetwork := newVSphereMachineSpec["network"].(map[string]interface{})
	oldVSphereMachineNetwork := oldVSphereMachineSpec["network"].(map[string]interface{})

	// allow changes to the devices.
	delete(oldVSphereMachineNetwork, "devices")
	delete(newVSphereMachineNetwork, "devices")

	// validate that IPAddrs in updaterequest are valid.
	spec := newTyped.Spec
	for i, device := range spec.Network.Devices {
		for j, ip := range device.IPAddrs {
			if _, _, err := net.ParseCIDR(ip); err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "network", fmt.Sprintf("devices[%d]", i), fmt.Sprintf("ipAddrs[%d]", j)), ip, "ip addresses should be in the CIDR format"))
			}
		}
	}

	if !reflect.DeepEqual(oldVSphereMachineSpec, newVSphereMachineSpec) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"), "cannot be modified"))
	}

	return nil, AggregateObjErrors(newTyped.GroupVersionKind().GroupKind(), newTyped.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachine) ValidateDelete(_ context.Context, _ *infrav1.VSphereMachine) (admission.Warnings, error) {
	return nil, nil
}

func validatePCIDevices(devices []infrav1.PCIDeviceSpec) field.ErrorList {
	var allErrs field.ErrorList

	for i, device := range devices {
		if device.VGPUProfile != "" && device.DeviceID == nil && device.VendorID == nil {
			// Valid case for vGPU.
			continue
		}
		if device.VGPUProfile == "" && device.DeviceID != nil && device.VendorID != nil {
			// Valid case for PCI Passthrough.
			continue
		}
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "template", "spec", "pciDevices").Index(i), device, "should have either deviceId + vendorId or vGPUProfile set"))
	}
	return allErrs
}
