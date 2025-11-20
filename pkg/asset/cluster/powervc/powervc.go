package powervc

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervc"
)

// Metadata converts an install configuration to PowerVC metadata.
func Metadata(infraID string, config *types.InstallConfig) *powervc.Metadata {
	return &powervc.Metadata{
		Cloud: config.Platform.PowerVC.Cloud,
		Identifier: map[string]string{
			"openshiftClusterID": infraID,
		},
	}
}
