package plugins

import (
	"github.com/abhinavdahiya/terraform-provider-azurerm-srv/srv"
	"github.com/hashicorp/terraform/plugin"
)

func init() {
	azurermSrvProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: srv.Provider,
		})
	}
	KnownPlugins["terraform-provider-azurerm-srv"] = azurermSrvProvider
}
