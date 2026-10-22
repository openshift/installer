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

	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta2"
	"sigs.k8s.io/cluster-api-provider-vsphere/feature"
	"sigs.k8s.io/cluster-api-provider-vsphere/internal/webhooks"
	"sigs.k8s.io/cluster-api-provider-vsphere/internal/webhooks/vmware/conversion"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/manager"
	pkgnetwork "sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/network"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-vmware-infrastructure-cluster-x-k8s-io-v1beta2-vspheremachine,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=vmware.infrastructure.cluster.x-k8s.io,resources=vspheremachines,versions=v1beta2,name=validation.vspheremachine.vmware.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1

// VSphereMachine implements a validation and defaulting webhook for VSphereMachine.
type VSphereMachine struct {
	// NetworkProvider is the network provider used by Supervisor based clusters
	NetworkProvider string
}

var _ admission.Validator[*vmwarev1.VSphereMachine] = &VSphereMachine{}

func (webhook *VSphereMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &vmwarev1.VSphereMachine{}).
		WithValidator(webhook).
		WithConverter(conversion.VSphereMachine).
		Complete()
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachine) ValidateCreate(_ context.Context, objTyped *vmwarev1.VSphereMachine) (admission.Warnings, error) {
	allErrs := validateNetwork(webhook.NetworkProvider, objTyped.Spec.Network, field.NewPath("spec", "network"))
	allErrs = append(allErrs, validatePolicies(objTyped.Spec.Policies, field.NewPath("spec", "policies"))...)

	return nil, webhooks.AggregateObjErrors(objTyped.GroupVersionKind().GroupKind(), objTyped.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachine) ValidateUpdate(_ context.Context, oldTyped, newTyped *vmwarev1.VSphereMachine) (admission.Warnings, error) {
	var allErrs field.ErrorList

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
	allErrs = append(allErrs, validatePolicies(newSpec.Policies, field.NewPath("spec", "policies"))...)

	return nil, webhooks.AggregateObjErrors(newTyped.GroupVersionKind().GroupKind(), newTyped.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereMachine) ValidateDelete(_ context.Context, _ *vmwarev1.VSphereMachine) (admission.Warnings, error) {
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
					primaryNetGVK := primary.NetworkRef.GroupVersionKind()
					if primaryNetGVK != pkgnetwork.NetworkGVKNSXTVPCSubnetSet {
						allErrs = append(allErrs, field.Invalid(
							fldPath.Child("interfaces", "primary", "network"),
							primaryNetGVK,
							fmt.Sprintf("only supports %s", pkgnetwork.NetworkGVKNSXTVPCSubnetSet)))
					}
				}
				for i, secondaryInterface := range network.Interfaces.Secondary {
					secondaryNetGVK := secondaryInterface.NetworkRef.GroupVersionKind()
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
					secondaryNetGVK := secondaryInterface.NetworkRef.GroupVersionKind()
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

	if len(network.VLANs) > 0 {
		allErrs = append(allErrs, validateVLANs(networkProvider, network, fldPath)...)
	}

	return allErrs
}

func validateVLANs(networkProvider string, network vmwarev1.VSphereMachineNetworkSpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if !feature.Gates.Enabled(feature.VLANSubinterface) {
		allErrs = append(allErrs, field.Forbidden(
			fldPath.Child("vlans"),
			"vlans can only be set when feature gate VLANSubinterface is enabled"))
		return allErrs
	}
	// vlan sub-interfaces feature only supports NSX-VPC Provider
	if networkProvider != manager.NSXVPCNetworkProvider {
		allErrs = append(allErrs, field.Forbidden(
			fldPath.Child("vlans"),
			fmt.Sprintf("vlans can only be set when network provider is %s", manager.NSXVPCNetworkProvider)))
		return allErrs
	}
	// vlan sub-interfaces only can link to a secondary interface
	if !network.Interfaces.IsDefined() || len(network.Interfaces.Secondary) == 0 {
		allErrs = append(allErrs, field.Required(
			fldPath.Child("vlans"),
			"vlans can only be specified if there are corresponding secondary interfaces"))
		return allErrs
	}
	// secondaryNames records all the secondary interface names, for checking vlan Link refers to an existing secondary interface name
	secondaryNames := map[string]struct{}{}
	// vlanNames records all the VLAN interface names, for checking duplicated VLAN interface name
	vlanNames := map[string]struct{}{}
	// vlanIDsPerLink maps a Link to its configured VLAN ID and name. Format: vlan Link(Secondary Interface Name) -> vlan ID -> vlan Name
	// for tracking assigned VLAN IDs per Link to detect duplicates
	vlanIDsPerLink := map[string]map[int32]string{}
	for _, s := range network.Interfaces.Secondary {
		secondaryNames[s.Name] = struct{}{}
	}
	for i, vlan := range network.VLANs {
		if _, ok := vlanNames[vlan.Name]; ok {
			allErrs = append(allErrs, field.Invalid(
				fldPath.Child("vlans").Index(i).Child("name"),
				vlan.Name,
				"VLAN name must be unique"))
		} else if vlan.Name == pkgnetwork.PrimaryInterfaceName {
			allErrs = append(allErrs, field.Invalid(
				fldPath.Child("vlans").Index(i).Child("name"),
				vlan.Name,
				"VLAN name is already in use by the primary interface"))
		} else if _, ok := secondaryNames[vlan.Name]; ok {
			allErrs = append(allErrs, field.Invalid(
				fldPath.Child("vlans").Index(i).Child("name"),
				vlan.Name,
				"VLAN name is already in use by a secondary interface"))
		} else {
			vlanNames[vlan.Name] = struct{}{}
		}
		if _, ok := secondaryNames[vlan.Link]; !ok {
			allErrs = append(allErrs, field.Invalid(
				fldPath.Child("vlans").Index(i).Child("link"),
				vlan.Link,
				"link must reference an existing secondary interface name"))
		}

		if vlan.ID == nil {
			allErrs = append(allErrs, field.Required(
				fldPath.Child("vlans").Index(i).Child("id"),
				"VLAN ID cannot be unset"))
			continue
		}

		if vlanIDsPerLink[vlan.Link] == nil {
			vlanIDsPerLink[vlan.Link] = map[int32]string{}
		}
		if existingVlanName, exists := vlanIDsPerLink[vlan.Link][*vlan.ID]; exists {
			allErrs = append(allErrs, field.Invalid(
				fldPath.Child("vlans").Index(i).Child("id"),
				vlan.ID,
				fmt.Sprintf("VLAN ID %d is already used by VLAN %q on the same link %q",
					*vlan.ID, existingVlanName, vlan.Link),
			))
		} else {
			vlanIDsPerLink[vlan.Link][*vlan.ID] = vlan.Name
		}
	}
	return allErrs
}

// validatePolicies validates the policies field is only set when the feature gate InfrastructurePolicies is enabled.
func validatePolicies(policies []vmwarev1.PolicyRef, fldPath *field.Path) field.ErrorList {
	if len(policies) == 0 {
		return nil
	}
	if feature.Gates.Enabled(feature.InfrastructurePolicies) {
		return nil
	}
	return field.ErrorList{
		field.Forbidden(fldPath, "policies can only be set when feature gate InfrastructurePolicies is enabled"),
	}
}
