package validation

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *openstack.Platform, n *types.Networking, ci *CloudInfo) field.ErrorList {
	var allErrs field.ErrorList
	fldPath := field.NewPath("platform", "openstack")

	// validate BYO machinesSubnet usage
	allErrs = append(allErrs, validateMachinesSubnet(p, n, ci, fldPath)...)

	// validate the externalNetwork
	allErrs = append(allErrs, validateExternalNetwork(p, ci, fldPath)...)

	// validate platform flavor
	allErrs = append(allErrs, validatePlatformFlavor(p, ci, fldPath)...)

	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, ci, fldPath.Child("defaultMachinePlatform"))...)
	}

	return allErrs
}

// validateMachinesSubnet validates the machines subnet and enforces proper byo subnet usage and returns a list of all validation errors
func validateMachinesSubnet(p *openstack.Platform, n *types.Networking, ci *CloudInfo, fldPath *field.Path) (allErrs field.ErrorList) {
	if p.MachinesSubnet != "" {
		if len(p.ExternalDNS) > 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("externalDNS"), p.ExternalDNS, "externalDNS is set, externalDNS is not supported when machinesSubnet is set"))
		}
		if !validUUIDv4(p.MachinesSubnet) {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("machinesSubnet"), errors.New("invalid subnet ID")))
		} else {
			if n.MachineNetwork[0].CIDR.String() != ci.MachinesSubnet.CIDR {
				allErrs = append(allErrs, field.InternalError(fldPath.Child("machinesSubnet"), fmt.Errorf("the first CIDR in machineNetwork, %s, doesn't match the CIDR of the machineSubnet, %s", n.MachineNetwork[0].CIDR.String(), ci.MachinesSubnet.CIDR)))
			}
		}
	}

	if len(p.ExternalDNS) > 0 && p.MachinesSubnet != "" {
		allErrs = append(allErrs, field.InternalError(fldPath.Child("machinesSubnet"), fmt.Errorf("externalDNS can't be set when using a custom machinesSubnet")))
	}
	return allErrs
}

// validateExternalNetwork validates the user's input for the externalNetwork and returns a list of all validation errors
func validateExternalNetwork(p *openstack.Platform, ci *CloudInfo, fldPath *field.Path) (allErrs field.ErrorList) {
	if ci.ExternalNetwork == nil {
		allErrs = append(allErrs, field.NotFound(fldPath.Child("externalNetwork"), p.ExternalNetwork))
	}
	return allErrs
}

// validatePlatformFlavor validates the platform flavor and returns a list of all validatoin errors
func validatePlatformFlavor(p *openstack.Platform, ci *CloudInfo, fldPath *field.Path) (allErrs field.ErrorList) {
	if ci.PlatformFlavor == nil {
		allErrs = append(allErrs, field.NotFound(fldPath.Child("computeFlavor"), p.FlavorName))
	}
	return allErrs
}
