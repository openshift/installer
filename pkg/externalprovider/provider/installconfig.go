package provider

import (
	"github.com/openshift/installer/pkg/types"
)

// InstallConfigExternalProvider describes the methods required for the installconfig asset.
type InstallConfigExternalProvider interface {
	// AddToInstallConfigPlatform adds the current platform to the installconfig.
	AddToInstallConfigPlatform(p *types.Platform) error
}
