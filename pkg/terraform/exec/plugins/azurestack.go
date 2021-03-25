package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-azurestack/azurestack"
)

func init() {
	azurestackProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: azurestack.Provider,
		})
	}
	KnownPlugins["terraform-provider-azurestack"] = azurestackProvider
}
