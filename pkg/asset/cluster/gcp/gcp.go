// Package gcp extracts GCP metadata from install configurations.
package gcp

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

// Metadata converts an install configuration to GCP metadata.
func Metadata(config *types.InstallConfig) *gcp.Metadata {
	return &gcp.Metadata{
		Region:    config.Platform.GCP.Region,
		ProjectID: config.Platform.GCP.ProjectID,
	}
}
