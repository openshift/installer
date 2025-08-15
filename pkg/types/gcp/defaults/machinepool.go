package defaults

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

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
}
