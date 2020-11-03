package validation

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	guuid "github.com/google/uuid"
	"github.com/openshift/installer/pkg/types/openstack"
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

	if controlPlane {
		allErrs = append(allErrs, validateFlavor(p.FlavorName, ci, ctrlPlaneFlavorMinimums, fldPath.Child("type"))...)
	} else {
		allErrs = append(allErrs, validateFlavor(p.FlavorName, ci, computeFlavorMinimums, fldPath.Child("type"))...)
	}

	allErrs = append(allErrs, validateZones(p.Zones, ci.Zones, fldPath.Child("zones"))...)
	allErrs = append(allErrs, validateUUIDV4s(p.AdditionalNetworkIDs, fldPath.Child("additionalNetworkIDs"))...)
	allErrs = append(allErrs, validateUUIDV4s(p.AdditionalSecurityGroupIDs, fldPath.Child("additionalSecurityGroupIDs"))...)

	return allErrs
}

func validateZones(input []string, available []string, fldPath *field.Path) field.ErrorList {
	// check if machinepool default
	if len(input) == 1 && input[0] == "" {
		return field.ErrorList{}
	}

	allErrs := field.ErrorList{}
	availableZones := sets.NewString(available...)
	for idx, zone := range input {
		if !availableZones.Has(zone) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("zone").Index(idx), zone, "Zone either does not exist in this cloud, or is not available"))
		}
	}

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

// validate flavor checks to make sure that a given flavor exists and meets the minimum requrement to run a cluster
// this function does not validate proper install config usage
func validateFlavor(flavorName string, ci *CloudInfo, req flavorRequirements, fldPath *field.Path) field.ErrorList {
	if flavorName == "" {
		return nil
	}

	flavor, _ := ci.Flavors[flavorName]
	if flavor.Flavor == nil {
		return field.ErrorList{field.NotFound(fldPath, flavorName)}
	}

	// OpenStack administrators don't always fill in accurate metadata for
	// baremetal flavors. Skipping validation.
	if flavor.Baremetal {
		return nil
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
		return nil
	}

	errString := "Flavor did not meet the following minimum requirements: "
	for i, err := range errs {
		errString = errString + err
		if i != len(errs)-1 {
			errString = errString + "; "
		}
	}

	return field.ErrorList{field.Invalid(fldPath, flavor.Name, errString)}
}
