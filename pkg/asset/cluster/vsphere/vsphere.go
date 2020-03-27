package vsphere

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// Metadata converts an install configuration to vSphere metadata.
func Metadata(config *types.InstallConfig) *vsphere.Metadata {
	return &vsphere.Metadata{
		VCenter:  config.VSphere.VCenter,
		Username: config.VSphere.Username,
		Password: config.VSphere.Password,
	}
}
