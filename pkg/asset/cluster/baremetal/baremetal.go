// Package baremetal extracts bare metal metadata from install
// configurations.
package baremetal

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

// Metadata converts an install configuration to bare metal metadata.
func Metadata(config *types.InstallConfig) *baremetal.Metadata {
	return &baremetal.Metadata{
		LibvirtURI:              config.Platform.BareMetal.LibvirtURI,
		BootstrapProvisioningIP: config.Platform.BareMetal.BootstrapProvisioningIP,
		ClusterProvisioningIP:   config.Platform.BareMetal.ClusterProvisioningIP,
	}
}
