package defaults

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

func findDefaultValue(value, valueTwo string) string {
	if value != "" {
		return value
	}
	return valueTwo
}

// SetMachinePoolDefaults sets the defaults for the platform.
func SetMachinePoolDefaults(platform *types.Platform, pool *gcp.MachinePool) {
	defaultMachinePool := platform.GCP.DefaultMachinePlatform
	if defaultMachinePool == nil {
		return
	}

	// Disk Type and Disk Size do not have default empty values, so their values are currently being
	// taken as what is provided. The encryption key can default to no value, so it must be set.
	if pool.EncryptionKey != nil {
		if defaultMachinePool.EncryptionKey != nil {
			if pool.EncryptionKey.KMSKeyServiceAccount == "" {
				pool.EncryptionKey.KMSKeyServiceAccount = defaultMachinePool.EncryptionKey.KMSKeyServiceAccount
			}

			if defaultMachinePool.EncryptionKey.KMSKey != nil {
				if pool.EncryptionKey.KMSKey != nil {
					if pool.EncryptionKey.KMSKey.ProjectID == "" {
						pool.EncryptionKey.KMSKey.ProjectID = findDefaultValue(defaultMachinePool.EncryptionKey.KMSKey.ProjectID, platform.GCP.ProjectID)
					}
					if pool.EncryptionKey.KMSKey.KeyRing == "" {
						pool.EncryptionKey.KMSKey.KeyRing = defaultMachinePool.EncryptionKey.KMSKey.KeyRing
					}
					if pool.EncryptionKey.KMSKey.Location == "" {
						pool.EncryptionKey.KMSKey.Location = findDefaultValue(defaultMachinePool.EncryptionKey.KMSKey.Location, platform.GCP.Region)
					}
					if pool.EncryptionKey.KMSKey.Name == "" {
						pool.EncryptionKey.KMSKey.Name = defaultMachinePool.EncryptionKey.KMSKey.Name
					}
				} else {
					pool.EncryptionKey.KMSKey = defaultMachinePool.EncryptionKey.KMSKey

					if pool.EncryptionKey.KMSKey.ProjectID == "" {
						pool.EncryptionKey.KMSKey.ProjectID = platform.GCP.ProjectID
					}
					if pool.EncryptionKey.KMSKey.Location == "" {
						pool.EncryptionKey.KMSKey.Location = platform.GCP.Region
					}
				}
			}
		}
	} else {
		// In the event that defaultMachinePool.EncryptionKey is nil this will not do anything.
		pool.EncryptionKey = defaultMachinePool.EncryptionKey
	}

	// The name and project in the OSImage are required values. The only time
	// these will be overridden is when the default machine pool is set but the
	// machine pool is not.
	if pool.OSImage == nil && defaultMachinePool.OSImage != nil {
		pool.OSImage = defaultMachinePool.OSImage
	}

	if pool.InstanceType == "" {
		pool.InstanceType = defaultMachinePool.InstanceType
	}

	// Zones are not combined to a single set that would include the machine pool
	// and the default machine pool.
	if len(pool.Zones) == 0 {
		pool.Zones = defaultMachinePool.Zones
	}

	// Tags are not combined to a single set from the default platform.
	// If any tags exist those tags will be used for the machine pool.
	if len(pool.Tags) == 0 {
		pool.Tags = defaultMachinePool.Tags
	}

	// Only override when no value is set.
	if pool.SecureBoot == "" {
		pool.SecureBoot = defaultMachinePool.SecureBoot
	}

	// Only override when no value is set
	if pool.OnHostMaintenance == "" {
		pool.OnHostMaintenance = defaultMachinePool.OnHostMaintenance
	}

	// Only override when no value is set
	if pool.ConfidentialCompute == "" {
		pool.ConfidentialCompute = defaultMachinePool.ConfidentialCompute
	}

	// Only override when no value is set
	if pool.ServiceAccount == "" {
		pool.ServiceAccount = defaultMachinePool.ServiceAccount
	}
}
