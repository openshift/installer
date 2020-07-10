// Package packet extracts packet metadata from install configurations.
package packet

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/packet"
)

// Metadata converts an install configuration to ovirt metadata.
func Metadata(config *types.InstallConfig) *packet.Metadata {
	m := packet.Metadata{}
	return &m
}
