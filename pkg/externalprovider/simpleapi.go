package externalprovider

import (
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
