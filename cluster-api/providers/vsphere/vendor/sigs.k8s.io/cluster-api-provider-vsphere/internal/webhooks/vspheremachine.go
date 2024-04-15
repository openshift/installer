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

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-vspheremachine,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=vspheremachines,versions=v1beta1,name=validation.vspheremachine.infrastructure.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-vspheremachine,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=vspheremachines,versions=v1beta1,name=default.vspheremachine.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

// VSphereMachineWebhook implements a validation and defaulting webhook for VSphereMachine.
type VSphereMachineWebhook struct{}

var _ webhook.CustomValidator = &VSphereMachineWebhook{}
var _ webhook.CustomDefaulter = &VSphereMachineWebhook{}

func (webhook *VSphereMachineWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&infrav1.VSphereMachine{}).
		WithValidator(webhook).
		WithDefaulter(webhook).
		Complete()
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (webhook *VSphereMachineWebhook) Default(_ context.Context, obj runtime.Object) error {
	objValue, ok := obj.(*infrav1.VSphereMachine)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereMachine but got a %T", obj))
	}
	if objValue.Spec.Datacenter == "" {
		objValue.Spec.Datacenter = "*"
	}
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachineWebhook) ValidateCreate(_ context.Context, raw runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList

	obj, ok := raw.(*infrav1.VSphereMachine)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereMachine but got a %T", raw))
	}
	spec := obj.Spec

	if spec.Network.PreferredAPIServerCIDR != "" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "PreferredAPIServerCIDR"), spec.Network.PreferredAPIServerCIDR, "cannot be set, as it will be removed and is no longer used"))
	}

	for i, device := range spec.Network.Devices {
		for j, ip := range device.IPAddrs {
			if _, _, err := net.ParseCIDR(ip); err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "network", fmt.Sprintf("devices[%d]", i), fmt.Sprintf("ipAddrs[%d]", j)), ip, "ip addresses should be in the CIDR format"))
			}
		}
	}

	if spec.GuestSoftPowerOffTimeout != nil {
		if spec.PowerOffMode != infrav1.VirtualMachinePowerOpModeTrySoft {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "guestSoftPowerOffTimeout"), spec.GuestSoftPowerOffTimeout, "should not be set in templates unless the powerOffMode is trySoft"))
		}
		if spec.GuestSoftPowerOffTimeout.Duration <= 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "guestSoftPowerOffTimeout"), spec.GuestSoftPowerOffTimeout, "should be greater than 0"))
		}
	}

	return nil, aggregateObjErrors(obj.GroupVersionKind().GroupKind(), obj.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachineWebhook) ValidateUpdate(_ context.Context, oldRaw runtime.Object, newRaw runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList

	newTyped, ok := newRaw.(*infrav1.VSphereMachine)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereMachine but got a %T", newRaw))
	}
	if newTyped.Spec.GuestSoftPowerOffTimeout != nil {
		if newTyped.Spec.PowerOffMode != infrav1.VirtualMachinePowerOpModeTrySoft {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "guestSoftPowerOffTimeout"), newTyped.Spec.GuestSoftPowerOffTimeout, "should not be set in templates unless the powerOffMode is trySoft"))
		}
		if newTyped.Spec.GuestSoftPowerOffTimeout.Duration <= 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "guestSoftPowerOffTimeout"), newTyped.Spec.GuestSoftPowerOffTimeout, "should be greater than 0"))
		}
	}

	newVSphereMachine, err := runtime.DefaultUnstructuredConverter.ToUnstructured(newRaw)
	if err != nil {
		return nil, apierrors.NewInternalError(errors.Wrap(err, "failed to convert new VSphereMachine to unstructured object"))
	}

	oldVSphereMachine, err := runtime.DefaultUnstructuredConverter.ToUnstructured(oldRaw)
	if err != nil {
		return nil, apierrors.NewInternalError(errors.Wrap(err, "failed to convert old VSphereMachine to unstructured object"))
	}

	newVSphereMachineSpec := newVSphereMachine["spec"].(map[string]interface{})
	oldVSphereMachineSpec := oldVSphereMachine["spec"].(map[string]interface{})

	allowChangeKeys := []string{"providerID", "powerOffMode", "guestSoftPowerOffTimeout"}
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

	return nil, aggregateObjErrors(newTyped.GroupVersionKind().GroupKind(), newTyped.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachineWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
