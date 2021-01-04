package externalprovider

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig/aws"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/externalprovider/provider/ovirt"
	"github.com/openshift/installer/pkg/types"
)

var providerRegistry = NewRegistry()

func init() {
	providerRegistry.Register(ovirt.NewOvirtProvider())
}

// AddToInstallConfigPlatform adds the current platform to the installconfig.
func AddToInstallConfigPlatform(
	ProviderName Name,
	p *types.Platform,
) error {
	provider := providerRegistry.MustGet(string(ProviderName))
	return provider.AddToInstallConfigPlatform(p)
}

// ValidateInstallConfig validates the install config.
func ValidateInstallConfig(
	ProviderName Name,
	Config *types.InstallConfig,
	File *asset.File,
	AWS *aws.Metadata,
	Azure *icazure.Metadata,
) error {
	provider := providerRegistry.MustGet(string(ProviderName))
	return provider.ValidateInstallConfig(
		Config,
		File,
		AWS,
		Azure,
	)
}

// PlatformCredsCheck validates the platform credentials.
func PlatformCredsCheck(
	ProviderName Name,
	Config *types.InstallConfig,
	File *asset.File,
	AWS *aws.Metadata,
	Azure *icazure.Metadata,
) error {
	provider := providerRegistry.MustGet(string(ProviderName))
	return provider.PlatformCredsCheck(
		Config,
		File,
		AWS,
		Azure,
	)
}

// PlatformPermsCheck validates the platform permissions.
func PlatformPermsCheck(
	ProviderName Name,
	Config *types.InstallConfig,
	File *asset.File,
	AWS *aws.Metadata,
	Azure *icazure.Metadata,
) error {
	provider := providerRegistry.MustGet(string(ProviderName))
	return provider.PlatformPermsCheck(
		Config,
		File,
		AWS,
		Azure,
	)
}

// PlatformProvisionCheck validates the if provisioning can commence on the platform.
func PlatformProvisionCheck(
	ProviderName Name,
	Config *types.InstallConfig,
	File *asset.File,
	AWS *aws.Metadata,
	Azure *icazure.Metadata,
) error {
	provider := providerRegistry.MustGet(string(ProviderName))
	return provider.PlatformProvisionCheck(
		Config,
		File,
		AWS,
		Azure,
	)
}
