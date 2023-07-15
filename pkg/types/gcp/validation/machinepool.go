package validation

import (
	"fmt"
	"regexp"
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

	if p.OSDisk.DiskType != "" {
		diskTypes := sets.NewString("pd-standard", "pd-ssd")
		if !diskTypes.Has(p.OSDisk.DiskType) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("diskType"), p.OSDisk.DiskType, diskTypes.List()))
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

// ValidateServiceAccount checks that the service account is only supplied for control plane nodes and during
// a shared vpn installation.
func ValidateServiceAccount(platform *gcp.Platform, p *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.Platform.GCP.ServiceAccount != "" {
		if p.Name != "master" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("serviceAccount"), p.Platform.GCP.ServiceAccount, fmt.Sprintf("service accounts only valid for master nodes, provided for %s nodes", p.Name)))
		}
		if platform.NetworkProjectID == "" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("serviceAccount"), p.Platform.GCP.ServiceAccount, "service accounts only valid for xpn installs"))
		}
	}
	return allErrs
}

// ValidateMasterDiskType checks that the specified disk type is valid for control plane.
func ValidateMasterDiskType(p *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.Name == "master" && p.Platform.GCP.OSDisk.DiskType == "pd-standard" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("diskType"), p.Platform.GCP.OSDisk.DiskType, fmt.Sprintf("%s not compatible with control planes.", p.Platform.GCP.OSDisk.DiskType)))
	}

	return allErrs
}

// ValidateDefaultDiskType checks that the specified disk type is valid for default GCP Machine Platform.
func ValidateDefaultDiskType(p *gcp.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p != nil && p.OSDisk.DiskType != "" {
		diskTypes := sets.NewString("pd-ssd")

		if !diskTypes.Has(p.OSDisk.DiskType) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("diskType"), p.OSDisk.DiskType, diskTypes.List()))
		}
	}

	return allErrs
}
