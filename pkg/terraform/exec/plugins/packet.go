package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-packet/packet"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: packet.Provider,
		})
	}
	KnownPlugins["terraform-provider-packet"] = exec
}
