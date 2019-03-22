// Package azure extracts AZURE metadata from install configurations.
package azure

import (
	"fmt"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

// Metadata converts an install configuration to Azure metadata.
func Metadata(clusterID, infraID string, config *types.InstallConfig) *azure.Metadata {
	return &azure.Metadata{
		Region: config.Platform.Azure.Region,
		Identifier: []map[string]string{{
			fmt.Sprintf("kubernetes.io.cluster.%s", infraID): "owned",
		}, {
			"openshiftClusterID": clusterID,
		}},
	}
}
