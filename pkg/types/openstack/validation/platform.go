package validation

import (
	"errors"
	"fmt"
	"net"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *openstack.Platform, n *types.Networking, fldPath *field.Path, fetcher ValidValuesFetcher, c *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	if allErrs = validateCloud(p, fldPath, fetcher); len(allErrs) == 0 {

		// validate BYO machinesSubnet usage
		allErrs = append(allErrs, validateMachinesSubnet(p, n, fldPath, fetcher)...)

		// validate the externalNetwork
		allErrs = append(allErrs, validateExternalNetwork(p, fldPath, fetcher)...)

		// validate compute flavor
		allErrs = append(allErrs, validateComputeFlavor(p, fldPath, fetcher)...)
	}

	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}

	if len(c.ObjectMeta.Name) > 14 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata", "name"), c.ObjectMeta.Name, "metadata name is too long, please restrict it to 14 characters"))
	}

	for _, ip := range p.ExternalDNS {
		if err := validate.IP(ip); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("externalDNS"), p.ExternalDNS, err.Error()))
		}
	}

	err := validateVIP(p.APIVIP, n)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, err.Error()))
	}

	err = validateVIP(p.IngressVIP, n)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, err.Error()))
	}

	return allErrs
}

func isValidValue(s string, validValues []string) bool {
	for _, v := range validValues {
		if s == v {
			return true
		}
	}
	return false
}

// validateVIP is a convenience function for validating VIP port and usage
func validateVIP(vip string, n *types.Networking) error {
	if vip != "" {
		if err := validate.IP(vip); err != nil {
			return err
		}

		if !n.MachineNetwork[0].CIDR.Contains(net.ParseIP(vip)) {
			return errors.New("IP is not in the machineNetwork")
		}
	}
	return nil
}

// validateCloud validates the cloud selected by the user as well as the clouds.yaml
func validateCloud(p *openstack.Platform, fldPath *field.Path, fetcher ValidValuesFetcher) (allErrs field.ErrorList) {
	validClouds, err := fetcher.GetCloudNames()
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath.Child("cloud"), errors.New("could not retrieve valid clouds")))
	} else if !isValidValue(p.Cloud, validClouds) {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("cloud"), p.Cloud, validClouds))
	}
	return allErrs
}

// validateMachinesSubnet validates the machines subnet and enforces proper byo subnet usage and returns a list of all validation errors
func validateMachinesSubnet(p *openstack.Platform, n *types.Networking, fldPath *field.Path, fetcher ValidValuesFetcher) (allErrs field.ErrorList) {
	if p.MachinesSubnet != "" {
		if len(p.ExternalDNS) > 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("externalDNS"), p.ExternalDNS, "externalDNS is set, externalDNS is not supported when machinesSubnet is set"))
		}
		if !validUUIDv4(p.MachinesSubnet) {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("machinesSubnet"), errors.New("invalid subnet ID")))
		} else {
			cidr, err := fetcher.GetSubnetCIDR(p.Cloud, p.MachinesSubnet)
			if err != nil {
				allErrs = append(allErrs, field.InternalError(fldPath.Child("machinesSubnet"), fmt.Errorf("invalid subnet %v", err)))
			}
			if n.MachineNetwork[0].CIDR.String() != cidr {
				allErrs = append(allErrs, field.InternalError(fldPath.Child("machinesSubnet"), fmt.Errorf("the first CIDR in machineNetwork, %s, doesn't match the CIDR of the machineSubnet, %s", n.MachineNetwork[0].CIDR.String(), cidr)))
			}
		}
	}

	if len(p.ExternalDNS) > 0 && p.MachinesSubnet != "" {
		allErrs = append(allErrs, field.InternalError(fldPath.Child("machinesSubnet"), fmt.Errorf("externalDNS can't be set when using a custom machinesSubnet")))
	}
	return allErrs
}

// validateExternalNetwork validates the user's input for the externalNetwork and returns a list of all validation errors
func validateExternalNetwork(p *openstack.Platform, fldPath *field.Path, fetcher ValidValuesFetcher) (allErrs field.ErrorList) {
	validNetworks, err := fetcher.GetNetworkNames(p.Cloud)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath.Child("externalNetwork"), errors.New("could not retrieve valid networks")))
	} else if !isValidValue(p.ExternalNetwork, validNetworks) {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("externalNetwork"), p.ExternalNetwork, validNetworks))
	}
	return allErrs
}

// validateComputeFlavor validates the compute flavor and returns a list of all validatoin errors
func validateComputeFlavor(p *openstack.Platform, fldPath *field.Path, fetcher ValidValuesFetcher) (allErrs field.ErrorList) {
	validFlavors, err := fetcher.GetFlavorNames(p.Cloud)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath.Child("computeFlavor"), errors.New("could not retrieve valid flavors")))
	} else if !isValidValue(p.FlavorName, validFlavors) {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("computeFlavor"), p.FlavorName, validFlavors))
	}
	return allErrs
}
