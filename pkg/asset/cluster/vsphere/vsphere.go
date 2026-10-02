package vsphere

import (
	"github.com/openshift/installer/pkg/types"
	typesvsphere "github.com/openshift/installer/pkg/types/vsphere"
)

// Metadata converts an install configuration to vSphere metadata.
func Metadata(config *types.InstallConfig) *typesvsphere.Metadata {
	terraformPlatform := "vsphere"

	metadata := &typesvsphere.Metadata{
		TerraformPlatform: terraformPlatform,
	}

	vcenterList := []typesvsphere.VCenters{}
	for _, vcenter := range config.VSphere.VCenters {
		credentials, _ := config.VSphere.CredentialsForVCenter(vcenter.Server)
		vcenterDef := typesvsphere.VCenters{
			VCenter:  vcenter.Server,
			Username: credentials.User,
			Password: credentials.Password,
		}
		vcenterList = append(vcenterList, vcenterDef)
	}
	metadata.VCenters = vcenterList

	return metadata
}
