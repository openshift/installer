package validation

import (
	"fmt"
	"net"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/kubevirt"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *kubevirt.Platform, fldPath *field.Path, c *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.Namespace == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("namespace"), "namespace is required"))
	}

	if p.NetworkName == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("networkName"), "networkName is required"))
	}

	if err := validate.IP(p.APIVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, err.Error()))
	} else if err := validateIPInMachineNetworkEntryList(c.MachineNetwork, p.APIVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, err.Error()))
	}

	if err := validate.IP(p.IngressVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, err.Error()))
	} else if err := validateIPInMachineNetworkEntryList(c.MachineNetwork, p.IngressVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, err.Error()))
	}

	return allErrs
}

func validateIPInMachineNetworkEntryList(machineNetworkEntryList []types.MachineNetworkEntry, ip string) error {
	ipAddr := net.ParseIP(ip)
	for _, machineNetworkEntry := range machineNetworkEntryList {
		if machineNetworkEntry.CIDR.Contains(ipAddr) {
			return nil
		}
	}
	return fmt.Errorf("IP must be in machine network range, machineNetworkEntryList: %s", machineNetworkEntryList)
}
