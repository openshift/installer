package vsphere

import (
	"github.com/openshift/installer/pkg/types"
	typesvsphere "github.com/openshift/installer/pkg/types/vsphere"
)

// Metadata converts an install configuration to vSphere metadata.
func Metadata(config *types.InstallConfig) *typesvsphere.Metadata {
	terraformPlatform := "vsphere"

	if vsphere := config.Platform.VSphere; vsphere != nil {
		if len(vsphere.FailureDomains) != 0 {
			terraformPlatform = typesvsphere.ZoningTerraformName
		}
	}

	return &typesvsphere.Metadata{
		VCenter:           config.VSphere.VCenter,
		Username:          config.VSphere.Username,
		Password:          config.VSphere.Password,
		TerraformPlatform: terraformPlatform,
	}
}
