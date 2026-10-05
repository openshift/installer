package vsphere

import (
	"fmt"

	"github.com/openshift/installer/pkg/types"
	typesvsphere "github.com/openshift/installer/pkg/types/vsphere"
)

// Metadata converts an install configuration to vSphere metadata.
func Metadata(config *types.InstallConfig) (*typesvsphere.Metadata, error) {
	terraformPlatform := "vsphere"

	metadata := &typesvsphere.Metadata{
		TerraformPlatform: terraformPlatform,
	}

	vcenterList := []typesvsphere.VCenters{}
	for _, vcenter := range config.VSphere.VCenters {
		credentials, err := config.VSphere.CredentialsForVCenter(vcenter.Server)
		if err != nil {
			return nil, fmt.Errorf("failed to get credentials for vcenter %q: %w", vcenter.Server, err)
		}
		vcenterDef := typesvsphere.VCenters{
			VCenter:  vcenter.Server,
			Username: credentials.User,
			Password: credentials.Password,
		}
		vcenterList = append(vcenterList, vcenterDef)
	}
	metadata.VCenters = vcenterList

	return metadata, nil
}
