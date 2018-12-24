package validation

import (
	"errors"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *openstack.Platform, fldPath *field.Path, fetcher ValidValuesFetcher) field.ErrorList {
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
		validImages, err := fetcher.GetImageNames(p.Cloud)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("baseImage"), errors.New("could not retrieve valid images")))
		} else if !isValidValue(p.BaseImage, validImages) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("baseImage"), p.BaseImage, validImages))
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
	}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	if err := validate.SubnetCIDR(&p.NetworkCIDRBlock.IPNet); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("NetworkCIDRBlock"), p.NetworkCIDRBlock, err.Error()))
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
