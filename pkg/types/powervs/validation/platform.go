package validation

import (
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/powervs"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *powervs.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// validate Zone
	if p.Zone == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("zone"), "zone must be specified"))
		// Region checking is nonsense if Zone is invalid
		return allErrs
	} else if ok := powervs.ValidateZone(p.Zone); !ok {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("zone"), p.Zone, powervs.ZoneNames()))
		// Region checking is nonsense if Zone is invalid
		return allErrs
	}

	// validate Region
	if p.Region == "" {
		p.Region = powervs.RegionFromZone(p.Zone)
	}
	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region not findable from specified zone"))
	} else if _, ok := powervs.Regions[p.Region]; !ok {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("region"), p.Region, powervs.RegionShortNames()))
	}

	// validate ServiceInstanceID
	if p.ServiceInstanceID != "" {
		_, err := uuid.Parse(p.ServiceInstanceID)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ServiceInstanceID"), p.ServiceInstanceID, "ServiceInstanceID must be a valid UUID"))
		}
	}

	// validate DefaultMachinePlatform
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	return allErrs
}
