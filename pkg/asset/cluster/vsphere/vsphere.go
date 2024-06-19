package vsphere

import (
	"github.com/openshift/installer/pkg/types"
	typesvsphere "github.com/openshift/installer/pkg/types/vsphere"
)

// Metadata converts an install configuration to vSphere metadata.
func Metadata(config *types.InstallConfig) *typesvsphere.Metadata {
	terraformPlatform := "vsphere"

	// Since currently we only support a single vCenter
	// just use the first entry in the VCenters slice.
	return &typesvsphere.Metadata{
		VCenter:           config.VSphere.VCenters[0].Server,
		Username:          config.VSphere.VCenters[0].Username,
		Password:          config.VSphere.VCenters[0].Password,
		TerraformPlatform: terraformPlatform,
	}
}
