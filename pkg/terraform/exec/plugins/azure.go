package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm"
)

func init() {
	azurermProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: azurerm.Provider,
		})
	}
	KnownPlugins["terraform-provider-azurerm"] = &TFPlugin{
		Name:    "azurerm",
		Exec:    azurermProvider,
		Version: GetAzurermVersion(),
	}
}
