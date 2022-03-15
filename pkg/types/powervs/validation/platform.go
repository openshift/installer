package validation

import (
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/powervs"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *powervs.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	//Validate Region
	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region must be specified"))
	} else if _, ok := powervs.Regions[p.Region]; !ok {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("region"), p.Region, powervs.RegionShortNames()))
	}

	//validate ServiceInstanceID
	if p.ServiceInstanceID != "" {
		_, err := uuid.Parse(p.ServiceInstanceID)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ServiceInstanceID"), p.ServiceInstanceID, "ServiceInstanceID must be a valid UUID"))
		}
	}

	//validate DefaultMachinePlatform
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	return allErrs
}
