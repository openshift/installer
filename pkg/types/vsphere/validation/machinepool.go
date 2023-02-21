package validation

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/validate"
)

const (
	defaultCoresPerSocket = int32(4)
	defaultNumCPUs        = int32(4)
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
	numCPUs := vspherePool.NumCPUs
	if numCPUs < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("cpus"), numCPUs, "number of CPUs must be positive"))
	}
	numCoresPerSocket := vspherePool.NumCoresPerSocket
	if numCoresPerSocket < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("coresPerSocket"), numCoresPerSocket, "cores per socket must be positive"))
	}

	// Either the number set by the user or a default value
	if numCPUs <= 0 {
		numCPUs = defaultNumCPUs
	}
	if numCoresPerSocket <= 0 {
		numCoresPerSocket = defaultCoresPerSocket
	}

	if numCoresPerSocket > numCPUs {
		errorMsg := fmt.Sprintf("cores per socket must be less than the number of CPUs (which is by default %d)", defaultNumCPUs)
		allErrs = append(allErrs, field.Invalid(fldPath.Child("coresPerSocket"), numCoresPerSocket, errorMsg))
	} else if numCPUs%numCoresPerSocket != 0 {
		errMsg := fmt.Sprintf("numCPUs specified should be a multiple of cores per socket (which is by default %d)", defaultCoresPerSocket)
		allErrs = append(allErrs, field.Invalid(fldPath.Child("cpus"), numCPUs, errMsg))
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
			if !zoneDefined {
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
