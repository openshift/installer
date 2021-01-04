package defaultprovider

import (
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
