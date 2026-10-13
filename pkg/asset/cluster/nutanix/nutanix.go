package nutanix

import (
	"strconv"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/nutanix"
)

// Metadata converts an install configuration to Nutanix metadata.
func Metadata(config *types.InstallConfig) *nutanix.Metadata {
	return &nutanix.Metadata{
		PrismCentral: config.Nutanix.PrismCentral.Endpoint.Address,
		Username:     config.Nutanix.PrismCentral.Username,
		Password:     config.Nutanix.PrismCentral.Password,
		Port:         strconv.Itoa(int(config.Nutanix.PrismCentral.Endpoint.Port)),
	}
}
