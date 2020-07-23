package validation

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"

	guuid "github.com/google/uuid"
	"github.com/openshift/installer/pkg/types/openstack"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)

type flavorRequirements struct {
	RAM, VCPUs, Disk int
}

var (
	ctrlPlaneFlavorMinimums = flavorRequirements{
		RAM:   16,
		VCPUs: 4,
		Disk:  25,
	}
	computeFlavorMinimums = flavorRequirements{
		RAM:   8,
		VCPUs: 2,
		Disk:  25,
	}
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *openstack.MachinePool, ci *CloudInfo, controlPlane bool, fldPath *field.Path) field.ErrorList {
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

	if p.FlavorName != "" {
		if controlPlane {
			allErrs = append(allErrs, validateMpoolFlavor(ci.Flavors[p.FlavorName], p.FlavorName, ctrlPlaneFlavorMinimums, fldPath)...)
		} else {
			allErrs = append(allErrs, validateMpoolFlavor(ci.Flavors[p.FlavorName], p.FlavorName, computeFlavorMinimums, fldPath)...)
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

func validateMpoolFlavor(flavor *flavors.Flavor, name string, req flavorRequirements, fldPath *field.Path) field.ErrorList {
	if flavor == nil {
		return field.ErrorList{field.NotFound(fldPath.Child("flavorName"), name)}
	}

	errs := []string{}
	if flavor.RAM < req.RAM {
		errs = append(errs, fmt.Sprintf("Must have minimum of %d GB RAM, had %d GB", req.RAM, flavor.RAM))
	}
	if flavor.VCPUs < req.VCPUs {
		errs = append(errs, fmt.Sprintf("Must have minimum of %d VCPUs, had %d", req.VCPUs, flavor.VCPUs))
	}
	if flavor.Disk < req.Disk {
		errs = append(errs, fmt.Sprintf("Must have minimum of %d GB Disk, had %d GB", req.Disk, flavor.Disk))
	}

	if len(errs) == 0 {
		return field.ErrorList{}
	}

	errString := "Flavor did not meet the following minimum requirements: "
	for i, err := range errs {
		errString = errString + err
		if i != len(errs)-1 {
			errString = errString + "; "
		}
	}

	return field.ErrorList{field.Invalid(fldPath.Child("flavorName"), flavor.Name, errString)}
}
