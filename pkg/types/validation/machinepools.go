package validation

import (
	"fmt"

	"github.com/cri-o/cpuset"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	awsvalidation "github.com/openshift/installer/pkg/types/aws/validation"
	"github.com/openshift/installer/pkg/types/azure"
	azurevalidation "github.com/openshift/installer/pkg/types/azure/validation"
	"github.com/openshift/installer/pkg/types/baremetal"
	baremetalvalidation "github.com/openshift/installer/pkg/types/baremetal/validation"
	"github.com/openshift/installer/pkg/types/gcp"
	gcpvalidation "github.com/openshift/installer/pkg/types/gcp/validation"
	"github.com/openshift/installer/pkg/types/kubevirt"
	kubevirtvalidation "github.com/openshift/installer/pkg/types/kubevirt/validation"
	"github.com/openshift/installer/pkg/types/libvirt"
	libvirtvalidation "github.com/openshift/installer/pkg/types/libvirt/validation"
	"github.com/openshift/installer/pkg/types/ovirt"
	ovirtvalidation "github.com/openshift/installer/pkg/types/ovirt/validation"
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
	allErrs = append(allErrs, validateMachinePoolPlatform(platform, &p.Platform, p, fldPath.Child("platform"))...)
	allErrs = append(allErrs, validateWorkloads(p.Workloads, fldPath.Child("workloads"))...)
	return allErrs
}

func validateMachinePoolPlatform(platform *types.Platform, p *types.MachinePoolPlatform, pool *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	platformName := platform.Name()
	validate := func(n string, value interface{}, validation func(*field.Path) field.ErrorList) {
		f := fldPath.Child(n)
		if platformName == n {
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
		validate(azure.Name, p.Azure, func(f *field.Path) field.ErrorList { return validateAzureMachinePool(p, pool, f) })
	}
	if p.GCP != nil {
		validate(gcp.Name, p.GCP, func(f *field.Path) field.ErrorList { return validateGCPMachinePool(platform, p, pool, f) })
	}
	if p.Libvirt != nil {
		validate(libvirt.Name, p.Libvirt, func(f *field.Path) field.ErrorList { return libvirtvalidation.ValidateMachinePool(p.Libvirt, f) })
	}
	if p.BareMetal != nil {
		validate(baremetal.Name, p.BareMetal, func(f *field.Path) field.ErrorList { return baremetalvalidation.ValidateMachinePool(p.BareMetal, f) })
	}
	if p.VSphere != nil {
		validate(vsphere.Name, p.VSphere, func(f *field.Path) field.ErrorList { return vspherevalidation.ValidateMachinePool(p.VSphere, f) })
	}
	if p.Ovirt != nil {
		validate(ovirt.Name, p.Ovirt, func(f *field.Path) field.ErrorList { return ovirtvalidation.ValidateMachinePool(p.Ovirt, f) })
	}
	if p.Kubevirt != nil {
		validate(kubevirt.Name, p.Kubevirt, func(f *field.Path) field.ErrorList { return kubevirtvalidation.ValidateMachinePool(p.Kubevirt, f) })
	}
	return allErrs
}

func validateGCPMachinePool(platform *types.Platform, p *types.MachinePoolPlatform, pool *types.MachinePool, f *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, gcpvalidation.ValidateMachinePool(platform.GCP, p.GCP, f)...)
	allErrs = append(allErrs, gcpvalidation.ValidateMasterDiskType(pool, f)...)

	return allErrs
}

func validateAzureMachinePool(p *types.MachinePoolPlatform, pool *types.MachinePool, f *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, azurevalidation.ValidateMachinePool(p.Azure, f)...)
	allErrs = append(allErrs, azurevalidation.ValidateMasterDiskType(pool, f)...)

	return allErrs
}

// validateWorkloads ensures that the given slice of workloads is either empty or:
//  - The name of the workload must be "management"
//  - Made of unique names (no duplicates)
//  - With a valid-looking CPU set description (note: Cannot actually verify against real hardware)
func validateWorkloads(workloads []types.Workload, f *field.Path) field.ErrorList {
	workloadNames := map[types.WorkloadName]bool{}
	allErrs := field.ErrorList{}
	for i, w := range workloads {
		fi := f.Index(i)
		if w.Name != types.ManagementWorkload {
			allErrs = append(allErrs, field.NotSupported(fi.Child("name"), w.Name, []string{string(types.ManagementWorkload)}))
		}
		if workloadNames[w.Name] {
			allErrs = append(allErrs, field.Duplicate(fi.Child("name"), w.Name))
		}
		workloadNames[w.Name] = true
		switch cpus, err := cpuset.Parse(w.CPUIDs); {
		case err != nil:
			allErrs = append(allErrs, field.Invalid(fi.Child("cpuIDs"), w.CPUIDs, "could not parse the cpuset"))
		case cpus.IsEmpty():
			allErrs = append(allErrs, field.Required(fi.Child("cpuIDs"), "must specify the CPU IDs"))
		}
	}
	return allErrs
}
