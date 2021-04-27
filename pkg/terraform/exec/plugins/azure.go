package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm"

	azuretypes "github.com/openshift/installer/pkg/types/azure"
)

func init() {
	azurermProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: azurerm.Provider,
		})
	}
	KnownPlugins["terraform-provider-azurerm"] = &TFPlugin{
		Name:      azuretypes.Name,
		Exec:      azurermProvider,
		Resources: []string{"compat"},
	}
}
