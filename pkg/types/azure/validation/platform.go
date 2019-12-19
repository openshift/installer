package validation

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *azure.Platform, publish types.PublishingStrategy, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region should be set to one of the supported Azure regions"))
	}
	if publish == types.ExternalPublishingStrategy {
		if p.BaseDomainResourceGroupName == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("baseDomainResourceGroupName"), "baseDomainResourceGroupName is the resource group name where the azure dns zone is deployed"))
		}
	}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	if p.VirtualNetwork != "" {
		if p.ComputeSubnet == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("computeSubnet"), "must provide a compute subnet when a virtual network is specified"))
		}
		if p.ControlPlaneSubnet == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("controlPlaneSubnet"), "must provide a control plane subnet when a virtual network is specified"))
		}
		if p.NetworkResourceGroupName == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("networkResourceGroupName"), "must provide a network resource group when a virtual network is specified"))
		}
	}
	if (p.ComputeSubnet != "" || p.ControlPlaneSubnet != "") && (p.VirtualNetwork == "" || p.NetworkResourceGroupName == "") {
		if p.VirtualNetwork == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("virtualNetwork"), "must provide a virtual network when supplying subnets"))
		}
		if p.NetworkResourceGroupName == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("networkResourceGroupName"), "must provide a network resource group when supplying subnets"))
		}
	}
	if (p.ResourceGroupName != "" && p.UserAssignedIdentity == "") {
		allErrs = append(allErrs, field.Required(fldPath.Child("userAssignedIdentity"), "must provide user assigned identity when supplying resource group"))
	}
	if (p.ResourceGroupName == "" && p.UserAssignedIdentity != "") {
		allErrs = append(allErrs, field.Required(fldPath.Child("resourceGroupName"), "must provide resource group name when supplying user assigned identity"))
	}
	return allErrs
}
