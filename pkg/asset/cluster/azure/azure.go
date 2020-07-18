// Package azure extracts AZURE metadata from install configurations.
package azure

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

// Metadata converts an install configuration to Azure metadata.
func Metadata(config *types.InstallConfig) *azure.Metadata {
	return &azure.Metadata{
		CloudName:         config.Platform.Azure.CloudName,
		Region:            config.Platform.Azure.Region,
		ResourceGroupName: config.Azure.ResourceGroupName,
	}
}
