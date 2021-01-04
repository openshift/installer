package defaultprovider

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/types"
)

// DefaultProvider is a default implementation of most functions required by externalprovider.Provider. This struct
// should be embedded in all providers to make adding new APIs easier.
type DefaultProvider struct {
}

// AddToInstallConfigPlatform adds the current platform to the installconfig.
func (d *DefaultProvider) AddToInstallConfigPlatform(_ *types.Platform) error {
	return nil
}

// ValidateInstallConfig validates the install config.
func (d *DefaultProvider) ValidateInstallConfig(
	_ *types.InstallConfig,
	_ *asset.File,
	_ *aws.Metadata,
	_ *azure.Metadata,
) error {
	return nil
}

// PlatformCredsCheck validates the platform credentials.
func (d *DefaultProvider) PlatformCredsCheck(
	_ *types.InstallConfig,
	_ *asset.File,
	_ *aws.Metadata,
	_ *azure.Metadata,
) error {
	return nil
}
