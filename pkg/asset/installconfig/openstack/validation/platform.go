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

	// validate floating ips
	allErrs = append(allErrs, validateFloatingIPs(p, ci, fldPath)...)

	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, ci, true, fldPath.Child("defaultMachinePlatform"))...)
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
	// Return an error if external network was specified in the install config, but hasn't been found
	if p.ExternalNetwork != "" && ci.ExternalNetwork == nil {
		allErrs = append(allErrs, field.NotFound(fldPath.Child("externalNetwork"), p.ExternalNetwork))
	}
	return allErrs
}

// validatePlatformFlavor validates the platform flavor and returns a list of all validation errors
func validatePlatformFlavor(p *openstack.Platform, ci *CloudInfo, fldPath *field.Path) (allErrs field.ErrorList) {
	flavor, ok := ci.Flavors[p.FlavorName]
	if !ok {
		allErrs = append(allErrs, field.NotFound(fldPath.Child("computeFlavor"), p.FlavorName))
		return allErrs
	}

	// OpenStack administrators don't always fill in accurate metadata for
	// baremetal flavors. Skipping validation.
	if flavor.Baremetal {
		return allErrs
	}

	errs := []string{}
	req := ctrlPlaneFlavorMinimums
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

	return append(allErrs, field.Invalid(fldPath.Child("flavorName"), flavor.Name, errString))
}

func validateFloatingIPs(p *openstack.Platform, ci *CloudInfo, fldPath *field.Path) (allErrs field.ErrorList) {
	if p.LbFloatingIP != "" {
		if ci.APIFIP == nil {
			allErrs = append(allErrs, field.NotFound(fldPath.Child("lbFloatingIP"), p.LbFloatingIP))
		} else if ci.APIFIP.Status != "DOWN" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("lbFloatingIP"), p.LbFloatingIP, "Floating IP already in use"))
		} else if p.ExternalNetwork == "" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("lbFloatingIP"), p.LbFloatingIP, "Cannot set floating ips when external network not specified"))
		}
	}

	if p.IngressFloatingIP != "" {
		if ci.IngressFIP == nil {
			allErrs = append(allErrs, field.NotFound(fldPath.Child("ingressFloatingIP"), p.IngressFloatingIP))
		} else if ci.IngressFIP.Status != "DOWN" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressFloatingIP"), p.IngressFloatingIP, "Floating IP already in use"))
		} else if p.ExternalNetwork == "" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressFloatingIP"), p.IngressFloatingIP, "Cannot set floating ips when external network not specified"))
		}
	}
	return allErrs
}
