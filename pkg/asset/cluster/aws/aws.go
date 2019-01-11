// Package aws extracts AWS metadata from install configurations.
package aws

import (
	"fmt"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// Metadata converts an install configuration to AWS metadata.
func Metadata(clusterID string, config *types.InstallConfig) *aws.Metadata {
	return &aws.Metadata{
		Region: config.Platform.AWS.Region,
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
