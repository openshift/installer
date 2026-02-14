package validation

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	awsvalidation "github.com/openshift/installer/pkg/types/aws/validation"
	"github.com/openshift/installer/pkg/types/azure"
	azurevalidation "github.com/openshift/installer/pkg/types/azure/validation"
	"github.com/openshift/installer/pkg/types/baremetal"
	baremetalvalidation "github.com/openshift/installer/pkg/types/baremetal/validation"
	"github.com/openshift/installer/pkg/types/gcp"
	gcpvalidation "github.com/openshift/installer/pkg/types/gcp/validation"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	ibmcloudvalidation "github.com/openshift/installer/pkg/types/ibmcloud/validation"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackvalidation "github.com/openshift/installer/pkg/types/openstack/validation"
	"github.com/openshift/installer/pkg/types/ovirt"
	ovirtvalidation "github.com/openshift/installer/pkg/types/ovirt/validation"
	"github.com/openshift/installer/pkg/types/powervc"
	"github.com/openshift/installer/pkg/types/powervs"
	powervsvalidation "github.com/openshift/installer/pkg/types/powervs/validation"
	"github.com/openshift/installer/pkg/types/vsphere"
	vspherevalidation "github.com/openshift/installer/pkg/types/vsphere/validation"
)

var (
	validHyperthreadingModes = map[types.HyperthreadingMode]bool{
		types.HyperthreadingDisabled: true,
		types.HyperthreadingEnabled:  true,
	}

	validHyperthreadingModeValues = func() []string {
		v := make([]string, 0, len(validHyperthreadingModes))
		for m := range validHyperthreadingModes {
			v = append(v, string(m))
		}
		return v
	}()

	validArchitectures = map[types.Architecture]bool{
		types.ArchitectureAMD64:   true,
		types.ArchitectureS390X:   true,
		types.ArchitecturePPC64LE: true,
		types.ArchitectureARM64:   true,
	}

	validArchitectureValues = func() []string {
		v := make([]string, 0, len(validArchitectures))
		for m := range validArchitectures {
			v = append(v, string(m))
		}
		return v
	}()
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(platform *types.Platform, p *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.Replicas != nil {
		if *p.Replicas < 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("replicas"), p.Replicas, "number of replicas must not be negative"))
		}
	} else {
		allErrs = append(allErrs, field.Required(fldPath.Child("replicas"), "replicas is required"))
	}
	if !validHyperthreadingModes[p.Hyperthreading] {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("hyperthreading"), p.Hyperthreading, validHyperthreadingModeValues))
	}
	if !validArchitectures[p.Architecture] {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("architecture"), p.Architecture, validArchitectureValues))
	}
	if platform.AWS != nil {
		allErrs = append(allErrs, awsvalidation.ValidateMachinePoolArchitecture(p, fldPath.Child("architecture"))...)
	}

	allErrs = append(allErrs, validateDiskSetup(p, fldPath.Child("diskSetup"))...)

	allErrs = append(allErrs, validateMachinePoolPlatform(platform, &p.Platform, p, fldPath.Child("platform"))...)
	return allErrs
}

func validateDiskSetup(p *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	foundEtcd := false
	foundSwap := false
	for _, ds := range p.DiskSetup {
		// outputting the yaml to make recognizing the issue easier for the user
		dsBytes, err := yaml.Marshal(ds)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath, err))
		}
		dsYaml := string(dsBytes)
		switch ds.Type {
		case types.UserDefined:
			if ds.UserDefined == nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("userDefined"), dsYaml, "userDefined configuration must be created"))
				continue
			}
			if len(ds.UserDefined.PlatformDiskID) > 12 {
				userDefinedPath := fldPath.Child("userDefined")
				allErrs = append(allErrs, field.Invalid(userDefinedPath.Child("platformDiskId"), dsYaml, "cannot be longer than 12 characters"))
				continue
			}
		case types.Etcd:
			if foundEtcd {
				allErrs = append(allErrs, field.TooMany(fldPath.Child("etcd"), 2, 1))
				continue
			}
			if ds.Etcd == nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("etcd"), dsYaml, "etcd configuration must be created"))
				continue
			}
			// etcd should only be setup on control plane, not any other machine type.
			if p.Name != "master" {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("etcd"), dsYaml, "cannot specify etcd on worker machine pools"))
				continue
			}
			foundEtcd = true
		case types.Swap:
			if p.Name == "master" {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("swap"), dsYaml, "swap is unsupported on control plane nodes"))
				continue
			}

			if foundSwap {
				allErrs = append(allErrs, field.TooMany(fldPath.Child("swap"), 2, 1))
				continue
			}
			if ds.Swap == nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("swap"), dsYaml, "swap configuration must be created"))
				continue
			}
			foundSwap = true
		}
	}

	return allErrs
}

func validateMachinePoolPlatform(platform *types.Platform, p *types.MachinePoolPlatform, pool *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	platformName := platform.Name()
	validate := func(n string, value interface{}, validation func(*field.Path) field.ErrorList) {
		f := fldPath.Child(n)
		if platformName == n || (platformName == powervc.Name && n == openstack.Name) {
			allErrs = append(allErrs, validation(f)...)
		} else {
			allErrs = append(allErrs, field.Invalid(f, value, fmt.Sprintf("cannot specify %q for machine pool when cluster is using %q", n, platformName)))
		}
	}
	if platform.AWS != nil {
		allErrs = append(allErrs, awsvalidation.ValidateAMIID(platform.AWS, p.AWS, fldPath.Child("aws"))...)
	}
	if p.AWS != nil {
		validate(aws.Name, p.AWS, func(f *field.Path) field.ErrorList { return awsvalidation.ValidateMachinePool(platform.AWS, p.AWS, f) })
	}
	if p.Azure != nil {
		validate(azure.Name, p.Azure, func(f *field.Path) field.ErrorList {
			return azurevalidation.ValidateMachinePool(p.Azure, pool.Name, platform.Azure, pool, f)
		})
	}
	if p.GCP != nil {
		validate(gcp.Name, p.GCP, func(f *field.Path) field.ErrorList { return validateGCPMachinePool(platform, p, pool, f) })
	}
	if p.IBMCloud != nil {
		validate(ibmcloud.Name, p.IBMCloud, func(f *field.Path) field.ErrorList {
			return ibmcloudvalidation.ValidateMachinePool(platform.IBMCloud, p.IBMCloud, f)
		})
	}
	if p.BareMetal != nil {
		validate(baremetal.Name, p.BareMetal, func(f *field.Path) field.ErrorList { return baremetalvalidation.ValidateMachinePool(p.BareMetal, f) })
	}
	if p.VSphere != nil {
		validate(vsphere.Name, p.VSphere, func(f *field.Path) field.ErrorList {
			return vspherevalidation.ValidateMachinePool(platform.VSphere, pool, f)
		})
	}
	if p.Ovirt != nil {
		validate(ovirt.Name, p.Ovirt, func(f *field.Path) field.ErrorList { return ovirtvalidation.ValidateMachinePool(p.Ovirt, f) })
	}
	if p.PowerVS != nil {
		validate(powervs.Name, p.PowerVS, func(f *field.Path) field.ErrorList {
			return powervsvalidation.ValidateMachinePool(platform.PowerVS, p.PowerVS, f)
		})
	}
	if p.OpenStack != nil {
		validate(openstack.Name, p.OpenStack, func(f *field.Path) field.ErrorList {
			return openstackvalidation.ValidateMachinePool(platform.OpenStack, p.OpenStack, pool.Name, f)
		})
	}
	return allErrs
}

func validateGCPMachinePool(platform *types.Platform, p *types.MachinePoolPlatform, pool *types.MachinePool, f *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, gcpvalidation.ValidateMachinePool(platform.GCP, p.GCP, f)...)
	allErrs = append(allErrs, gcpvalidation.ValidateMasterDiskType(pool, f)...)
	allErrs = append(allErrs, gcpvalidation.ValidateServiceAccount(platform.GCP, pool, f)...)

	return allErrs
}
