package validation

import (
	"fmt"
	"regexp"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/strings/slices"

	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/validate"
)

const (
	defaultCoresPerSocket = int32(4)
	defaultNumCPUs        = int32(4)
	// Maximum number of data disks allowed to be added to a VM.
	maxVSphereDataDisks = 29
	// Max length of a DataDisk name.
	maxVSphereDataDiskNameLength = 80
	// Max size of any data disk in vSphere is 62 TiB.  We are currently limiting to 16TiB (16384 GiB) as a starting point.
	maxVSphereDataDiskSize = 16384
)

var (
	// validProvisioningModes lists all valid data disk provisioning modes.
	validProvisioningModes = []machinev1beta1.ProvisioningMode{
		machinev1beta1.ProvisioningModeThin,
		machinev1beta1.ProvisioningModeThick,
		machinev1beta1.ProvisioningModeEagerlyZeroed,
	}
	// vSphereDataDiskNamePattern is used to validate the name of a data disk.
	vSphereDataDiskNamePattern = regexp.MustCompile(`^[a-zA-Z0-9]([-_a-zA-Z0-9]*[a-zA-Z0-9])?$`)
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

	// We need to validate the data disks to make sure configs are valid.
	dataDisks := vspherePool.DataDisks
	if len(dataDisks) > 0 {
		disksPath := fldPath.Child("disks")
		if len(dataDisks) > maxVSphereDataDisks {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("dataDisks"), len(dataDisks), fmt.Sprintf("data disk count must not exceed %d", maxVSphereDataDisks)))
		}

		// Check each data disk
		for i, disk := range dataDisks {
			diskPath := disksPath.Index(i)

			// Validate data disk name
			if len(disk.Name) == 0 {
				allErrs = append(allErrs, field.Required(diskPath.Child("name"), "data disk name must be set"))
			} else {
				if len(disk.Name) > maxVSphereDataDiskNameLength {
					allErrs = append(allErrs, field.Invalid(diskPath.Child("name"), len(disk.Name), fmt.Sprintf("data disk name must not exceed %d", maxVSphereDataDiskNameLength)))
				}
				if vSphereDataDiskNamePattern.FindStringSubmatch(disk.Name) == nil {
					allErrs = append(allErrs, field.Invalid(diskPath.Child("name"), disk.Name, "data disk name must consist only of alphanumeric characters, hyphens and underscores, and must start and end with an alphanumeric character."))
				}
			}

			// Check if disk size is set
			if disk.SizeGiB == 0 {
				allErrs = append(allErrs, field.Required(diskPath.Child("sizeGiB"), "data disk size must be set"))
			}

			// Check data disk does not exceed max disk size
			if disk.SizeGiB > maxVSphereDataDiskSize {
				allErrs = append(allErrs, field.Invalid(diskPath.Child("sizeGiB"), disk.SizeGiB, fmt.Sprintf("data disk size (GiB) must not exceed %d", maxVSphereDataDiskSize)))
			}

			// Validate provisioning modes
			if len(disk.ProvisioningMode) > 0 {
				validModesSet := sets.NewString()
				for _, m := range validProvisioningModes {
					validModesSet.Insert(string(m))
				}
				if !validModesSet.Has(string(disk.ProvisioningMode)) {
					allErrs = append(allErrs, field.NotSupported(diskPath, disk.ProvisioningMode, validModesSet.List()))
				}
			}
		}
	}

	if len(vspherePool.Zones) > 0 {
		var zoneRefs []string
		if len(platform.FailureDomains) == 0 {
			return append(allErrs, field.Required(fldPath.Child("zones"), "failureDomains must be defined if zones are defined"))
		}
		for index, zone := range vspherePool.Zones {
			err := validate.ClusterName1035(zone)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("zones"), vspherePool.Zones, err.Error()))
			}

			// This checks to make sure ref is not duplicated
			if slices.Contains(zoneRefs, zone) {
				allErrs = append(allErrs, field.Duplicate(fldPath.Child("zones").Index(index), zone))
			} else {
				zoneRefs = append(zoneRefs, zone)
			}

			// Verify zone references a valid failure domain
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
