package plugins

import (
	metal "github.com/equinix/terraform-provider-equinix-metal/packet"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: metal.Provider,
		})
	}
	// TODO(displague) update to equinix-metal when TF 0.13+ sdk can be used
	KnownPlugins["terraform-provider-packet"] = exec
}
