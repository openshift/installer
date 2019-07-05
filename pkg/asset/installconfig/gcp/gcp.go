package gcp

import (
	"github.com/openshift/installer/pkg/types/gcp"
)

// Platform collects GCP-specific configuration.
func Platform() (*gcp.Platform, error) {
	return &gcp.Platform{}, nil
}
