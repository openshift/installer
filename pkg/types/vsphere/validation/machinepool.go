package validation

import (
	"github.com/openshift/installer/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/validate"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(platform *vsphere.Platform, machinePool *types.MachinePool, fldPath *field.Path) field.ErrorList {
	vspherePool := machinePool.Platform.VSphere
	allErrs := field.ErrorList{}
	if vspherePool.DiskSizeGB < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("diskSizeGB"), vspherePool.DiskSizeGB, "storage disk size must be positive"))
	}
	if vspherePool.MemoryMiB < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("memoryMB"), vspherePool.MemoryMiB, "memory size must be positive"))
	}
	if vspherePool.NumCPUs < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("cpus"), vspherePool.NumCPUs, "number of CPUs must be positive"))
	}
	if vspherePool.NumCoresPerSocket < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("coresPerSocket"), vspherePool.NumCoresPerSocket, "cores per socket must be positive"))
	}

	defaultCoresPerSocket := int32(4)
	defaultNumCPUs := int32(4)
	if vspherePool.NumCPUs > 0 {
		if vspherePool.NumCoresPerSocket > vspherePool.NumCPUs {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("coresPerSocket"), vspherePool.NumCoresPerSocket, "cores per socket must be less than number of CPUs"))
		} else if vspherePool.NumCoresPerSocket > 0 && vspherePool.NumCPUs%vspherePool.NumCoresPerSocket != 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("cpus"), vspherePool.NumCPUs, "numCPUs specified should be a multiple of cores per socket"))
		} else if vspherePool.NumCPUs%defaultCoresPerSocket != 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("cpus"), vspherePool.NumCPUs, "numCPUs specified should be a multiple of cores per socket which is by default 4"))
		}
	} else if vspherePool.NumCoresPerSocket > defaultNumCPUs {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("coresPerSocket"), vspherePool.NumCoresPerSocket, "cores per socket must be less than number of CPUs which is by default 4"))
	}

	if len(vspherePool.Zones) > 0 {
		if len(platform.FailureDomains) == 0 {
			return append(allErrs, field.Required(fldPath.Child("zones"), "failureDomains must be defined if zones are defined"))
		}
		for _, zone := range vspherePool.Zones {
			err := validate.ClusterName1035(zone)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("zones"), vspherePool.Zones, err.Error()))
			}
			zoneDefined := false
			for _, failureDomain := range platform.FailureDomains {
				if failureDomain.Name == zone {
					zoneDefined = true
				}
			}
			if zoneDefined == false {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("zones"), zone, "zone not defined in failureDomains"))
			}
		}
	} else if len(platform.FailureDomains) > 0 {
		for _, failureDomain := range platform.FailureDomains {
			vspherePool.Zones = append(vspherePool.Zones, failureDomain.Name)
		}
	}
	return allErrs
}
