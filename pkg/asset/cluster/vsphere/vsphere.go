package vsphere

import (
	"context"

	"github.com/openshift/installer/pkg/asset/installconfig"
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

// PreTerraform performs any infrastructure initialization which must
// happen before Terraform creates the remaining infrastructure.
func PreTerraform(ctx context.Context, infraID string, installConfig *installconfig.InstallConfig) error {
	// TODO: create VM Template using cachedImage
	return nil
}
