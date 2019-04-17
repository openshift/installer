// Package ovirt extracts ovirt metadata from install configurations.
package ovirt

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ovirt"
)

// Metadata converts an install configuration to ovirt metadata.
func Metadata(config *types.InstallConfig) *ovirt.Metadata {
	return &ovirt.Metadata{
		ClusterID:  config.Ovirt.ClusterID,
	}
}
