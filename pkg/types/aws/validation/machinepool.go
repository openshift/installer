package validation

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

var (
	validArchitectures = map[types.Architecture]bool{
		types.ArchitectureAMD64: true,
		types.ArchitectureARM64: true,
	}

	// ValidArchitectureValues lists the supported arches for AWS
	ValidArchitectureValues = func() []string {
		v := make([]string, 0, len(validArchitectures))
		for m := range validArchitectures {
			v = append(v, string(m))
		}
		return v
	}()
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(platform *aws.Platform, p *aws.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for i, zone := range p.Zones {
		if !strings.HasPrefix(zone, platform.Region) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("zones").Index(i), zone, fmt.Sprintf("Zone not in configured region (%s)", platform.Region)))
		}
	}

	if p.IOPS < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("iops"), p.IOPS, "Storage IOPS must be positive"))
	}
	if p.Size < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("size"), p.Size, "Storage size must be positive"))
	}
	return allErrs
}

// ValidateAMIID check the AMI ID is set for a machine pool.
func ValidateAMIID(platform *aws.Platform, p *aws.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	pool := &aws.MachinePool{AMIID: platform.AMIID}
	pool.Set(platform.DefaultMachinePlatform)
	pool.Set(p)

	// regions is a list of regions for which the user should set AMI ID as copying the AMI to these regions
	// is known to not be supported.
	regions := sets.NewString("us-gov-west-1", "us-gov-east-1", "us-iso-east-1", "cn-north-1", "cn-northwest-1")
	if pool.AMIID == "" && regions.Has(platform.Region) {
		allErrs = append(allErrs, field.Required(fldPath, fmt.Sprintf("AMI ID must be provided for regions %s", strings.Join(regions.List(), ", "))))
	}
	return allErrs
}

// ValidateMachinePoolArchitecture checks that a valid architecture is set for a machine pool.
func ValidateMachinePoolArchitecture(pool *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if !validArchitectures[pool.Architecture] {
		allErrs = append(allErrs, field.NotSupported(fldPath, pool.Architecture, ValidArchitectureValues))
	}
	return allErrs
}
