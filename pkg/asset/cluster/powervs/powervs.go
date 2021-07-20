// Package powervs extracts Power VS metadata from install configurations.
package powervs

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervs"
)

// Metadata converts an install configuration to PowerVS metadata.
func Metadata(config *types.InstallConfig) *powervs.Metadata {
	return &powervs.Metadata{
		Region: config.Platform.PowerVS.Region,
		Zone:   config.Platform.PowerVS.Zone,
	}
}
