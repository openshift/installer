package nutanix

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/nutanix"
)

// Metadata converts an install configuration to Nutanix metadata.
func Metadata(config *types.InstallConfig) *nutanix.Metadata {
	return &nutanix.Metadata{
		PrismCentral: config.Nutanix.PrismCentral,
		Username:     config.Nutanix.Username,
		Password:     config.Nutanix.Password,
		Port:         config.Nutanix.Port,
	}
}
