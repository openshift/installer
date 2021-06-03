package validation

import (
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var (
	// Regions is a map of IBM Cloud regions where VPCs are supported.
	// The key of the map is the short name of the region. The value
	// of the map is the long name of the region.
	Regions = map[string]string{
		// https://cloud.ibm.com/docs/vpc?topic=vpc-creating-a-vpc-in-a-different-region
		"us-south": "US South (Dallas)",
		"us-east":  "US East (Washington DC)",
		"eu-gb":    "United Kindom (London)",
		"eu-de":    "EU Germany (Frankfurt)",
		"jp-tok":   "Japan (Tokyo)",
		"jp-osa":   "Japan (Osaka)",
		"au-syd":   "Australia (Sydney)",
		"ca-tor":   "Canada (Toronto)",
	}

	regionShortNames = func() []string {
		keys := make([]string, len(Regions))
		i := 0
		for r := range Regions {
			keys[i] = r
			i++
		}
		return keys
	}()
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *ibmcloud.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region must be specified"))
	} else if _, ok := Regions[p.Region]; !ok {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("region"), p.Region, regionShortNames))
	}

	allErrs = append(allErrs, validateVPCConfig(p, fldPath)...)

	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	return allErrs
}

func validateVPCConfig(p *ibmcloud.Platform, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.VPC != "" || len(p.Subnets) > 0 || p.VPCResourceGroupName != "" {
		if p.VPC == "" {
			allErrs = append(allErrs, field.Required(path.Child("vpc"), "vpc is required when specifying subnets or vpcResourceGroupName"))
		}
		if len(p.Subnets) == 0 {
			allErrs = append(allErrs, field.Required(path.Child("subnets"), "subnets is required when specifying vpc or vpcResourceGroupName"))
		}
		if p.VPCResourceGroupName == "" {
			allErrs = append(allErrs, field.Required(path.Child("vpcResourceGroupName"), "vpcResourceGroupName is required when specifying vpc or subnets"))
		}
	}
	return allErrs
}
