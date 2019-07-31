package validation

import (
	"errors"
	"fmt"

	"github.com/apparentlymart/go-cidr/cidr"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *openstack.Platform, n *types.Networking, fldPath *field.Path, fetcher ValidValuesFetcher) field.ErrorList {
	allErrs := field.ErrorList{}
	validClouds, err := fetcher.GetCloudNames()
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath.Child("cloud"), errors.New("could not retrieve valid clouds")))
	} else if !isValidValue(p.Cloud, validClouds) {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("cloud"), p.Cloud, validClouds))
	} else {
		validRegions, err := fetcher.GetRegionNames(p.Cloud)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("region"), errors.New("could not retrieve valid regions")))
		} else if !isValidValue(p.Region, validRegions) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("region"), p.Region, validRegions))
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

	// Validate VIP Values
	if err := validate.IP(p.APIVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, err.Error()))
	} else if n != nil {
		expectedAPIVIP, _ := cidr.Host(&n.MachineCIDR.IPNet, 5)
		if p.APIVIP != expectedAPIVIP.String() {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, fmt.Errorf("the API VIP is expected to be %s", expectedAPIVIP).Error()))
		}
	}

	if err := validate.IP(p.DNSVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("dnsVIP"), p.DNSVIP, err.Error()))
	} else if n != nil {
		expectedDNSVIP, _ := cidr.Host(&n.MachineCIDR.IPNet, 6)
		if p.DNSVIP != expectedDNSVIP.String() {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("dnsVIP"), p.DNSVIP, fmt.Errorf("the DNS VIP is expected to be %s", expectedDNSVIP).Error()))
		}
	}

	if err := validate.IP(p.IngressVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, err.Error()))
	} else if n != nil {
		expectedIngressVIP, _ := cidr.Host(&n.MachineCIDR.IPNet, 7)
		if p.IngressVIP != expectedIngressVIP.String() {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, fmt.Errorf("the Ingress VIP is expected to be %s", expectedIngressVIP).Error()))
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
