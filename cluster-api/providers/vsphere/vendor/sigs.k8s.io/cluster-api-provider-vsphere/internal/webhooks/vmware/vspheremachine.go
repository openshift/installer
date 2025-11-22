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
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/feature"
	"sigs.k8s.io/cluster-api-provider-vsphere/internal/webhooks"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/manager"
	pkgnetwork "sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/network"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-vmware-infrastructure-cluster-x-k8s-io-v1beta1-vspheremachine,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=vmware.infrastructure.cluster.x-k8s.io,resources=vspheremachines,versions=v1beta1,name=validation.vspheremachine.vmware.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-vmware-infrastructure-cluster-x-k8s-io-v1beta1-vspheremachine,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=vmware.infrastructure.cluster.x-k8s.io,resources=vspheremachines,versions=v1beta1,name=default.vspheremachine.vmware.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

// VSphereMachine implements a validation and defaulting webhook for VSphereMachine.
type VSphereMachine struct {
	// NetworkProvider is the network provider used by Supervisor based clusters
	NetworkProvider string
}

var _ webhook.CustomValidator = &VSphereMachine{}
var _ webhook.CustomDefaulter = &VSphereMachine{}

func (webhook *VSphereMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&vmwarev1.VSphereMachine{}).
		WithValidator(webhook).
		WithDefaulter(webhook, admission.DefaulterRemoveUnknownOrOmitableFields).
		Complete()
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (webhook *VSphereMachine) Default(_ context.Context, _ runtime.Object) error {
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachine) ValidateCreate(_ context.Context, object runtime.Object) (admission.Warnings, error) {
	objTyped, ok := object.(*vmwarev1.VSphereMachine)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereMachine but got a %T", object))
	}

	allErrs := validateNetwork(webhook.NetworkProvider, objTyped.Spec.Network, field.NewPath("spec", "network"))

	return nil, webhooks.AggregateObjErrors(objTyped.GroupVersionKind().GroupKind(), objTyped.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachine) ValidateUpdate(_ context.Context, oldRaw runtime.Object, newRaw runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList

	newTyped, ok := newRaw.(*vmwarev1.VSphereMachine)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereMachine but got a %T", newRaw))
	}

	oldTyped, ok := oldRaw.(*vmwarev1.VSphereMachine)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereMachine but got a %T", oldRaw))
	}

	newSpec, oldSpec := newTyped.Spec, oldTyped.Spec

	// In VM operator, following fields are immutable, so CAPV should not allow to update them.
	// - ImageName
	// - ClassName
	// - StorageClass
	// - MinHardwareVersion
	if newSpec.ImageName != oldSpec.ImageName {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "imageName"), "cannot be modified"))
	}

	if newSpec.ClassName != oldSpec.ClassName {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "className"), "cannot be modified"))
	}

	if newSpec.StorageClass != oldSpec.StorageClass {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "storageClass"), "cannot be modified"))
	}

	if newSpec.MinHardwareVersion != oldSpec.MinHardwareVersion {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "minHardwareVersion"), "cannot be modified"))
	}

	if !reflect.DeepEqual(newSpec.Network.Interfaces, oldSpec.Network.Interfaces) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "network", "interfaces"), "cannot be modified"))
	}

	allErrs = append(allErrs, validateNetwork(webhook.NetworkProvider, newSpec.Network, field.NewPath("spec", "network"))...)

	return nil, webhooks.AggregateObjErrors(newTyped.GroupVersionKind().GroupKind(), newTyped.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachine) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func validateNetwork(networkProvider string, network vmwarev1.VSphereMachineNetworkSpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if network.Interfaces.IsDefined() {
		if !feature.Gates.Enabled(feature.MultiNetworks) {
			allErrs = append(allErrs, field.Forbidden(
				fldPath.Child("interfaces"),
				"interfaces can only be set when feature gate MultiNetworks is enabled"))
		} else {
			// Validate network type is supported
			switch networkProvider {
			case manager.NSXVPCNetworkProvider:
				primary := network.Interfaces.Primary
				if primary.IsDefined() {
					primaryNetGVK := primary.Network.GroupVersionKind()
					if primaryNetGVK != pkgnetwork.NetworkGVKNSXTVPCSubnetSet {
						allErrs = append(allErrs, field.Invalid(
							fldPath.Child("interfaces", "primary", "network"),
							primaryNetGVK,
							fmt.Sprintf("only supports %s", pkgnetwork.NetworkGVKNSXTVPCSubnetSet)))
					}
				}
				for i, secondaryInterface := range network.Interfaces.Secondary {
					secondaryNetGVK := secondaryInterface.Network.GroupVersionKind()
					if secondaryNetGVK != pkgnetwork.NetworkGVKNSXTVPCSubnetSet && secondaryNetGVK != pkgnetwork.NetworkGVKNSXTVPCSubnet {
						allErrs = append(allErrs, field.Invalid(
							fldPath.Child("interfaces", "secondary").Index(i).Child("network"),
							secondaryNetGVK,
							fmt.Sprintf("only supports %s or %s", pkgnetwork.NetworkGVKNSXTVPCSubnetSet, pkgnetwork.NetworkGVKNSXTVPCSubnet)))
					}
				}
			case manager.VDSNetworkProvider:
				if network.Interfaces.Primary.IsDefined() {
					allErrs = append(allErrs, field.Forbidden(
						fldPath.Child("interfaces", "primary"),
						"primary interface can not be set when network provider is vsphere-network"))
				}
				for i, secondaryInterface := range network.Interfaces.Secondary {
					secondaryNetGVK := secondaryInterface.Network.GroupVersionKind()
					if secondaryNetGVK != pkgnetwork.NetworkGVKNetOperator {
						allErrs = append(allErrs, field.Invalid(
							fldPath.Child("interfaces", "secondary").Index(i).Child("network"),
							secondaryNetGVK,
							fmt.Sprintf("only supports %s", pkgnetwork.NetworkGVKNetOperator)))
					}
				}
			default:
				allErrs = append(allErrs, field.Forbidden(fldPath.Child("interfaces"), fmt.Sprintf("interfaces can not be set when network provider is %s", networkProvider)))
			}

			// Validate interface names are unique
			interfaceNames := map[string]struct{}{pkgnetwork.PrimaryInterfaceName: {}}
			for i, secondaryInterface := range network.Interfaces.Secondary {
				if _, ok := interfaceNames[secondaryInterface.Name]; ok {
					allErrs = append(allErrs, field.Invalid(
						fldPath.Child("interfaces", "secondary").Index(i).Child("name"),
						secondaryInterface.Name,
						"interface name is already in use"))
				} else {
					interfaceNames[secondaryInterface.Name] = struct{}{}
				}
			}
		}
	}
	return allErrs
}
