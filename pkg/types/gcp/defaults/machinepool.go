package defaults

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

// Apply sets values from the default machine platform to the machinepool.
func Apply(defaultMachinePlatform, machinePool *gcp.MachinePool) {
	tempMP := &gcp.MachinePool{}
	tempMP.Set(defaultMachinePlatform)
	tempMP.Set(machinePool)
	machinePool.Set(tempMP)
}

// SetMachinePoolDefaults sets the defaults for the platform.
func SetMachinePoolDefaults(platform *types.Platform, pool *gcp.MachinePool) {
	if pool == nil {
		return
	}

	if ek := pool.EncryptionKey; ek != nil {
		if kms := ek.KMSKey; kms != nil {
			if kms.ProjectID == "" {
				kms.ProjectID = platform.GCP.ProjectID
			}
			if kms.Location == "" {
				kms.Location = platform.GCP.Region
			}
		}
	}

	// Set the default Disk Type for the Instance type when the instance type is provided
	// by the user and the Disk type is not.
	if pool.InstanceType != "" && pool.OSDisk.DiskType == "" {
		family := gcp.GetGCPInstanceFamily(pool.InstanceType)
		if _, ok := gcp.InstanceTypeToDiskTypeMap[family]; ok {
			pool.OSDisk.DiskType = gcp.DefaultDiskTypeForInstanceAndProjectID(pool.InstanceType, platform.GCP.ProjectID, platform.GCP.Region)
		}
	}
}
