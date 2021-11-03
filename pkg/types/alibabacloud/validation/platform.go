package validation

import (
	"github.com/openshift/installer/pkg/types/alibabacloud"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *alibabacloud.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region must be specified"))
	}

	if p.ResourceGroupID == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("resourceGroupID"), "resource group ID must be specified"))
	}

	return allErrs
}
