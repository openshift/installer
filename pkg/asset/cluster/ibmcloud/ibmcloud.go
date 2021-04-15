// Package ibmcloud extracts IBM Cloud metadata from install configurations.
package ibmcloud

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ibmcloud"
)

// Metadata converts an install configuration to IBM Cloud metadata.
func Metadata(infraID string, config *types.InstallConfig) *ibmcloud.Metadata {
	return &ibmcloud.Metadata{
		Region:            config.Platform.IBMCloud.Region,
		ResourceGroupName: config.Platform.IBMCloud.ClusterResourceGroupName(infraID),
	}
}
