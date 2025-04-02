package validation

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(platform *gcp.Platform, p *gcp.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for i, zone := range p.Zones {
		if !strings.HasPrefix(zone, platform.Region) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("zones").Index(i), zone, fmt.Sprintf("Zone not in configured region (%s)", platform.Region)))
		}
	}
	if p.OSDisk.DiskSizeGB != 0 {
		if p.OSDisk.DiskSizeGB < 16 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("diskSizeGB"), p.OSDisk.DiskSizeGB, "must be at least 16GB in size"))
		} else if p.OSDisk.DiskSizeGB > 65536 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("diskSizeGB"), p.OSDisk.DiskSizeGB, "exceeding maximum GCP disk size limit, must be below 65536"))
		}
	}

	if diskType := p.OSDisk.DiskType; diskType != "" && !gcp.ComputeSupportedDisks.Has(diskType) {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("diskType"), diskType, sets.List(gcp.ComputeSupportedDisks)))
	}

	if p.ConfidentialCompute != "" && p.ConfidentialCompute != string(gcp.DisabledFeature) {
		if p.OnHostMaintenance != string(gcp.OnHostMaintenanceTerminate) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("OnHostMaintenance"), p.OnHostMaintenance, fmt.Sprintf("OnHostMaintenace must be set to Terminate when ConfidentialCompute is %s", p.ConfidentialCompute)))
		}

		instanceType, _, _ := strings.Cut(p.InstanceType, "-")
		confidentialCompute := gcp.ConfidentialComputePolicy(p.ConfidentialCompute)
		if confidentialCompute == gcp.ConfidentialComputePolicy(gcp.EnabledFeature) {
			confidentialCompute = gcp.ConfidentialComputePolicySEV
		}
		supportedMachineTypes, ok := gcp.ConfidentialComputePolicyToSupportedInstanceType[confidentialCompute]
		if !ok {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("confidentialCompute"), p.ConfidentialCompute, fmt.Sprintf("Unknown confidential computing technology %s", p.ConfidentialCompute)))
		} else if !slices.Contains(supportedMachineTypes, instanceType) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("type"), p.InstanceType, fmt.Sprintf("Machine type do not support %s. Machine types supporting %s: %s", p.ConfidentialCompute, p.ConfidentialCompute, strings.Join(supportedMachineTypes, ", "))))
		}
	}

	for i, tag := range p.Tags {
		if tag == "" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("tags").Index(i), tag, fmt.Sprintf("tag can not be empty")))
		} else if !unicode.IsLetter(rune(tag[0])) || (!unicode.IsLetter(rune(tag[len(tag)-1])) && !unicode.IsNumber(rune(tag[len(tag)-1]))) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("tags").Index(i), tag, fmt.Sprintf("tag can only start with a letter and must end with a letter or a number")))
		} else if !regexp.MustCompile(`^[a-z0-9-]*$`).MatchString(tag) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("tags").Index(i), tag, fmt.Sprintf("tag can only contain lowercase letters, numbers, and dashes")))
		} else if len(tag) > 63 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("tags").Index(i), tag, fmt.Sprintf("maximum number of characters is 63")))
		}
	}
	return allErrs
}

// ValidateServiceAccount does not do any checks on the service account since it can be set for all nodes and
// in non-shared vpn installations.
func ValidateServiceAccount(platform *gcp.Platform, p *types.MachinePool, fldPath *field.Path) field.ErrorList {
	return field.ErrorList{}
}

// ValidateMasterDiskType checks that the specified disk type is valid for control plane.
func ValidateMasterDiskType(p *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.Name == "master" && p.Platform.GCP.OSDisk.DiskType != "" {
		if !gcp.ControlPlaneSupportedDisks.Has(p.Platform.GCP.OSDisk.DiskType) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("diskType"), p.Platform.GCP.OSDisk.DiskType, sets.List(gcp.ControlPlaneSupportedDisks)))
		}
	}
	return allErrs
}

// ValidateDefaultDiskType checks that the specified disk type is valid for default GCP Machine Platform.
func ValidateDefaultDiskType(p *gcp.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p != nil && p.OSDisk.DiskType != "" {
		if !gcp.ControlPlaneSupportedDisks.Has(p.OSDisk.DiskType) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("diskType"), p.OSDisk.DiskType, sets.List(gcp.ControlPlaneSupportedDisks)))
		}
	}

	return allErrs
}
