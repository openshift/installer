// Package aws extracts AWS metadata from install configurations.
package aws

import (
	"fmt"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// Metadata converts an install configuration to AWS metadata.
func Metadata(clusterID, infraID string, config *types.InstallConfig) *aws.Metadata {
	return &aws.Metadata{
		Region:          config.Platform.AWS.Region,
		CustomEndpoints: config.Platform.AWS.FetchCustomEndpoints(),
		Identifier: []map[string]string{{
			fmt.Sprintf("kubernetes.io/cluster/%s", infraID): "owned",
		}, {
			"openshiftClusterID": clusterID,
		}},
	}
}
