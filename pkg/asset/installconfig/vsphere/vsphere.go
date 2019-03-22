// Package vsphere collects vSphere-specific configuration.
package vsphere

import (
	"github.com/openshift/installer/pkg/types/vsphere"
)

// Platform collects vSphere-specific configuration.
func Platform() (*vsphere.Platform, error) {
	return &vsphere.Platform{}, nil
}
