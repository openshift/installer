// Package google extracts GCP metadata from install configurations.
package google

import (
	"fmt"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/google"
)

// Metadata converts an install configuration to GCP metadata.
func Metadata(clusterID string, config *types.InstallConfig) *google.Metadata {
	return &google.Metadata{
		Region: config.Platform.GCP.Region,
		Identifier: []map[string]string{
			{
				"openshiftClusterID": clusterID,
			},
			{
				fmt.Sprintf("kubernetes.io/cluster/%s", config.ObjectMeta.Name): "owned",
			},
		},
	}
}
