package provider

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig/aws"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/types"
)

// InstallConfigExternalProvider describes the methods required for the installconfig asset.
type InstallConfigExternalProvider interface {
	// AddToInstallConfigPlatform adds the current platform to the installconfig.
	AddToInstallConfigPlatform(p *types.Platform) error

	// ValidateInstallConfig validates the install config.
	ValidateInstallConfig(
		Config *types.InstallConfig,
		File *asset.File,
		AWS *aws.Metadata,
		Azure *icazure.Metadata,
	) error

	// PlatformCredsCheck validates the platform credentials.
	PlatformCredsCheck(
		Config *types.InstallConfig,
		File *asset.File,
		AWS *aws.Metadata,
		Azure *icazure.Metadata,
	) error

	// PlatformPermsCheck validates the platform permissions.
	PlatformPermsCheck(
		Config *types.InstallConfig,
		File *asset.File,
		AWS *aws.Metadata,
		Azure *icazure.Metadata,
	) error

	// PlatformProvisionCheck validates the if provisioning can commence on the platform.
	PlatformProvisionCheck(
		Config *types.InstallConfig,
		File *asset.File,
		AWS *aws.Metadata,
		Azure *icazure.Metadata,
	) error
}
