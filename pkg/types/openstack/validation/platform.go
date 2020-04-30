package validation

import (
	"errors"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *openstack.Platform, n *types.Networking, fldPath *field.Path, fetcher ValidValuesFetcher, c *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	validClouds, err := fetcher.GetCloudNames()
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath.Child("cloud"), errors.New("could not retrieve valid clouds")))
	} else if !isValidValue(p.Cloud, validClouds) {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("cloud"), p.Cloud, validClouds))
	} else {
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
		validNetworks, err := fetcher.GetNetworkNames(p.Cloud)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("externalNetwork"), errors.New("could not retrieve valid networks")))
		} else if !isValidValue(p.ExternalNetwork, validNetworks) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("externalNetwork"), p.ExternalNetwork, validNetworks))
		}
		validFlavors, err := fetcher.GetFlavorNames(p.Cloud)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("computeFlavor"), errors.New("could not retrieve valid flavors")))
		} else if !isValidValue(p.FlavorName, validFlavors) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("computeFlavor"), p.FlavorName, validFlavors))
		}
		p.TrunkSupport = "0"
		netExts, err := fetcher.GetNetworkExtensionsAliases(p.Cloud)
		if err != nil {
			logrus.Warning("Could not retrieve networking extension aliases. Assuming trunk ports are not supported.")
		} else {
			if isValidValue("trunk", netExts) {
				p.TrunkSupport = "1"
			}
		}
		p.OctaviaSupport = "0"
		serviceCatalog, err := fetcher.GetServiceCatalog(p.Cloud)
		if err != nil {
			logrus.Warning("Could not retrieve service catalog. Assuming there is no Octavia load balancer service available.")
		} else {
			if isValidValue("octavia", serviceCatalog) {
				p.OctaviaSupport = "1"
			}
		}
	}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}

	if len(c.ObjectMeta.Name) > 14 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata", "name"), c.ObjectMeta.Name, "metadata name is too long, please restrict it to 14 characters"))
	}

	if len(p.ExternalDNS) > 0 && p.MachinesSubnet != "" {
		allErrs = append(allErrs, field.InternalError(fldPath.Child("machinesSubnet"), fmt.Errorf("externalDNS can't be set when using a custom machinesSubnet")))
	}

	for _, ip := range p.ExternalDNS {
		if err := validate.IP(ip); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("externalDNS"), p.ExternalDNS, err.Error()))
		}
	}

	err = validateVIP(p.APIVIP, n)
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
