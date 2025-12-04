package validation

import (
	"fmt"
	"regexp"
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

	// validArchitectureValues lists the supported arches for AWS
	validArchitectureValues = func() []string {
		v := make([]string, 0, len(validArchitectures))
		for m := range validArchitectures {
			v = append(v, string(m))
		}
		return v
	}()

	validMetadataAuthValues = sets.NewString("Required", "Optional")

	validConfidentialComputePolicy = []aws.ConfidentialComputePolicy{
		aws.ConfidentialComputePolicyDisabled,
		aws.ConfidentialComputePolicySEVSNP,
	}

	// awsDedicatedHostNamePattern is a regex expression that defines the dedicated host id format.
	awsDedicatedHostNamePattern = regexp.MustCompile(`^h-[0-9a-f]{17}$`)
)

// AWS has a limit of 16 security groups. See:
// https://docs.aws.amazon.com/vpc/latest/userguide/amazon-vpc-limits.html
// We set a user limit of 10 and reserve 6 for use by OpenShift.
const maxUserSecurityGroupsCount = 10

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(platform *aws.Platform, p *aws.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for i, zone := range p.Zones {
		if !strings.HasPrefix(zone, platform.Region) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("zones").Index(i), zone, fmt.Sprintf("Zone not in configured region (%s)", platform.Region)))
		}
	}

	if p.EC2RootVolume.Type != "" {
		allErrs = append(allErrs, validateVolumeSize(p, fldPath)...)
		allErrs = append(allErrs, validateIOPS(p, fldPath)...)
	}

	if p.EC2Metadata.Authentication != "" && !validMetadataAuthValues.Has(p.EC2Metadata.Authentication) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("authentication"), p.EC2Metadata.Authentication, "must be either Required or Optional"))
	}

	allErrs = append(allErrs, validateSecurityGroups(platform, p, fldPath)...)
	allErrs = append(allErrs, ValidateCPUOptions(p, fldPath)...)

	if p.HostPlacement != nil {
		allErrs = append(allErrs, validateHostPlacement(p, fldPath.Child("hostPlacement"))...)
	}

	return allErrs
}

func validateHostPlacement(p *aws.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.HostPlacement.Affinity == nil {
		allErrs = append(allErrs, field.Required(fldPath.Child("affinity"), "affinity is required when hostPlacement is configured"))
		return allErrs // Can't validate further without affinity
	}

	switch *p.HostPlacement.Affinity {
	case aws.HostAffinityAnyAvailable:
		if len(p.HostPlacement.DedicatedHost) > 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("dedicatedHost"), "dedicatedHost is required when 'affinity' is set to DedicatedHost, and forbidden otherwise"))
		}
	case aws.HostAffinityDedicatedHost:
		if len(p.HostPlacement.DedicatedHost) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("dedicatedHost"), "dedicatedHost is required when 'affinity' is set to DedicatedHost, and forbidden otherwise"))
		} else {
			for index, host := range p.HostPlacement.DedicatedHost {
				hostPath := fldPath.Child("dedicatedHost").Index(index)
				if len(host.ID) == 0 {
					allErrs = append(allErrs, field.Required(hostPath.Child("id"), "a hostID must be specified when configuring 'dedicatedHost'"))
				} else if !awsDedicatedHostNamePattern.MatchString(host.ID) {
					allErrs = append(allErrs, field.Invalid(hostPath.Child("id"), host.ID, "id must start with 'h-' followed by 17 lowercase hexadecimal characters (0-9 and a-f)"))
				}
			}
		}
	default:
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("affinity"), p.HostPlacement.Affinity, []aws.HostAffinity{aws.HostAffinityAnyAvailable, aws.HostAffinityDedicatedHost}))
	}

	return allErrs
}

func validateSecurityGroups(platform *aws.Platform, p *aws.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(p.AdditionalSecurityGroupIDs) > 0 && len(platform.VPC.Subnets) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("platform.subnets"), "subnets must be provided when additional security groups are present"))
	}

	// The installer also creates a security group: `${var.cluster_id}-master-sg/${var.cluster_id}-worker-sg`
	if count := len(p.AdditionalSecurityGroupIDs); count > maxUserSecurityGroupsCount {
		allErrs = append(allErrs, field.TooMany(fldPath, count, maxUserSecurityGroupsCount))
	}

	return allErrs
}

func validateVolumeSize(p *aws.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	volumeSize := p.EC2RootVolume.Size

	if volumeSize <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("size"), volumeSize, "volume size value must be a positive number"))
	}

	return allErrs
}

func validateIOPS(p *aws.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	volumeType := strings.ToLower(p.EC2RootVolume.Type)
	iops := p.EC2RootVolume.IOPS

	switch volumeType {
	case "io1", "io2":
		if iops <= 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("iops"), iops, "iops must be a positive number"))
		}
	case "gp3":
		if iops < 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("iops"), iops, "iops must be a positive number"))
		}
	case "gp2", "st1", "sc1", "standard":
		if iops != 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("iops"), iops, fmt.Sprintf("iops not supported for type %s", volumeType)))
		}
	default:
		allErrs = append(allErrs, field.Invalid(fldPath.Child("type"), volumeType, fmt.Sprintf("failed to find volume type %s", volumeType)))
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
	regions := sets.NewString("us-iso-east-1", "us-isob-east-1", "us-iso-west-1", "cn-north-1", "cn-northwest-1")
	if pool.AMIID == "" && regions.Has(platform.Region) {
		allErrs = append(allErrs, field.Required(fldPath, fmt.Sprintf("AMI ID must be provided for regions %s", strings.Join(regions.List(), ", "))))
	}
	return allErrs
}

// ValidateMachinePoolArchitecture checks that a valid architecture is set for a machine pool.
func ValidateMachinePoolArchitecture(pool *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if !validArchitectures[pool.Architecture] {
		allErrs = append(allErrs, field.NotSupported(fldPath, pool.Architecture, validArchitectureValues))
	}
	return allErrs
}

// ValidateCPUOptions checks that valid CPU options are set for a machine pool.
func ValidateCPUOptions(p *aws.MachinePool, fldPath *field.Path) field.ErrorList {
	if p.CPUOptions == nil {
		return nil
	}

	allErrs := field.ErrorList{}

	if *p.CPUOptions == (aws.CPUOptions{}) {
		allErrs = append(
			allErrs,
			field.Invalid(
				fldPath.Child("cpuOptions"),
				"{}",
				"At least one field must be set if cpuOptions is provided",
			),
		)
	}

	if p.CPUOptions.ConfidentialCompute != nil {
		switch *p.CPUOptions.ConfidentialCompute {
		case aws.ConfidentialComputePolicyDisabled, aws.ConfidentialComputePolicySEVSNP:
			// Valid values
		default:
			allErrs = append(
				allErrs,
				field.NotSupported(
					fldPath.Child("confidentialCompute"),
					p.CPUOptions.ConfidentialCompute,
					validConfidentialComputePolicy,
				),
			)
		}
	}

	return allErrs
}
