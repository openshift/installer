package validation

import (
	"github.com/openshift/installer/pkg/types/azure"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *azure.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region should be set to one of the supported Azure regions"))
	}
	if p.BaseDomainResourceGroupName == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("baseDomainResourceGroupName"), "baseDomainResourceGroupName is the resource group name where the azure dns zone is deployed"))
	}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	return allErrs
}
