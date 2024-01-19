// Package gcp extracts GCP metadata from install configurations.
package gcp

import (
	"fmt"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

// Metadata converts an install configuration to GCP metadata.
func Metadata(config *types.InstallConfig) *gcp.Metadata {
	// leave the private zone domain blank when not using a pre-created private zone
	privateZoneDomain := fmt.Sprintf("%s.", config.ClusterDomain())
	if config.GCP.Network == "" || config.GCP.NetworkProjectID == "" {
		privateZoneDomain = ""
	}

	return &gcp.Metadata{
		Region:            config.Platform.GCP.Region,
		ProjectID:         config.Platform.GCP.ProjectID,
		NetworkProjectID:  config.Platform.GCP.NetworkProjectID,
		PrivateZoneDomain: privateZoneDomain,
	}
}
