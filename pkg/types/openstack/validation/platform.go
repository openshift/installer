package validation

import (
	"errors"
        "net"

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
		netExts, err := fetcher.GetNetworkExtensionsAliases(p.Cloud)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("trunkSupport"), errors.New("could not retrieve networking extension aliases")))
		} else {
			if isValidValue("trunk", netExts) {
				p.TrunkSupport = "1"
			} else {
				p.TrunkSupport = "0"
			}
		}
		serviceCatalog, err := fetcher.GetServiceCatalog(p.Cloud)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("octaviaSupport"), errors.New("could not retrieve service catalog")))
		} else {
			if isValidValue("octavia", serviceCatalog) {
				p.OctaviaSupport = "1"
			} else {
				p.OctaviaSupport = "0"
			}
		}
	}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}

	if len(c.ObjectMeta.Name) > 14 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata", "name"), c.ObjectMeta.Name, "metadata name is too long, please restrict it to 14 characters"))
	}

	for _, ip := range p.ExternalDNS {
		if err := validate.IP(ip); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ExternalDNS"), p.ExternalDNS, err.Error()))
		}
	}

	// Check if snatCRSubnet is a valid subnet or IP
	if p.AciNetExt.ClusterSNATSubnet != "" {
		_, _, err := net.ParseCIDR(p.AciNetExt.ClusterSNATSubnet)
		if err != nil {
			ip := net.ParseIP(p.AciNetExt.ClusterSNATSubnet)
			if ip == nil {
				// error
				allErrs = append(allErrs, field.Invalid(fldPath.Child("ClusterSNATSubnet"), p.AciNetExt.ClusterSNATSubnet, err.Error()))
			}	
		}
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
