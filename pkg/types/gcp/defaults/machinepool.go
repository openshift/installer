package defaults

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

// SetMachinePoolDefaults sets the defaults for the platform.
func SetMachinePoolDefaults(platform *types.Platform, pool *gcp.MachinePool) {
	if pool.EncryptionKey != nil {
		if pool.EncryptionKey.KMSKey.ProjectID == "" {
			pool.EncryptionKey.KMSKey.ProjectID = platform.GCP.ProjectID
		}
		if pool.EncryptionKey.KMSKey.Location == "" {
			pool.EncryptionKey.KMSKey.Location = platform.GCP.Region
		}
	}
}
