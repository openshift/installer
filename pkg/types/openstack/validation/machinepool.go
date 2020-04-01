package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	guuid "github.com/google/uuid"
	"github.com/openshift/installer/pkg/types/openstack"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *openstack.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// Validate Root Volumes
	if p.RootVolume != nil {
		if p.RootVolume.Type == "" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("rootVolume").Child("type"), p.RootVolume.Type, "Volume type must be specified to use root volumes"))
		}
		if p.RootVolume.Size <= 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("rootVolume").Child("size"), p.RootVolume.Size, "Volume size must be greater than zero to use root volumes"))
		}
	}

	allErrs = append(allErrs, validateUUIDV4s(p.AdditionalNetworkIDs, fldPath.Child("additionalNetworkIDs"))...)
	allErrs = append(allErrs, validateUUIDV4s(p.AdditionalSecurityGroupIDs, fldPath.Child("additionalSecurityGroupIDs"))...)

	return allErrs
}

func validateUUIDV4s(input []string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for idx, uuid := range input {
		if !validUUIDv4(uuid) {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(idx), uuid, "valid UUID v4 must be specified"))
		}
	}

	return allErrs
}

// validUUIDv4 checks if string is in UUID v4 format
// For more information: https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)
func validUUIDv4(s string) bool {
	uuid, err := guuid.Parse(s)
	if err != nil {
		return false
	}

	// check that version of the uuid
	if uuid.Version().String() != "VERSION_4" {
		return false
	}

	return true
}
