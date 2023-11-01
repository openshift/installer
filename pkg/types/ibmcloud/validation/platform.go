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
		"eu-es":    "Spain (Madrid)",
		"jp-tok":   "Japan (Tokyo)",
		"jp-osa":   "Japan (Osaka)",
		"au-syd":   "Australia (Sydney)",
		"ca-tor":   "Canada (Toronto)",
		"br-sao":   "Brazil (Sao Paulo)",
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

	if p.VPCName != "" {
		if p.ControlPlaneSubnets == nil {
			allErrs = append(allErrs, field.Required(fldPath.Child("controlPlaneSubnets"), "must provided at least one control plane subnet when a VPC is specified"))
		}
		if p.ComputeSubnets == nil {
			allErrs = append(allErrs, field.Required(fldPath.Child("computeSubnets"), "must provide at least one compute subnet when a VPC is specified"))
		}
	} else if p.ControlPlaneSubnets != nil || p.ComputeSubnets != nil {
		allErrs = append(allErrs, field.Required(fldPath.Child("vpcName"), "must provide a VPC name when supplying subnets"))
	}

	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p, p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	return allErrs
}
