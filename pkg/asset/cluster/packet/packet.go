// Package packet extracts packet metadata from install configurations.
package packet

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/packet"
)

// Metadata converts an install configuration to Packet metadata.
func Metadata(config *types.InstallConfig) *packet.Metadata {
	m := packet.Metadata{
		FacilityCode: config.Platform.Packet.FacilityCode,
		ProjectID:    config.Platform.Packet.ProjectID,
	}
	return &m
}
