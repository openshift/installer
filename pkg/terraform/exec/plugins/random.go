package plugins

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-random/v2/random"
)

func init() {
	randomProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: random.Provider,
		})
	}
	KnownPlugins["terraform-provider-random"] = randomProvider
}
