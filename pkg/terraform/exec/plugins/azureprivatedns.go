package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/openshift/installer/pkg/terraform/exec/plugins/azureprivatedns"
)

func init() {
	azurePrivateDNSProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: azureprivatedns.Provider,
		})
	}
	KnownPlugins["terraform-provider-azureprivatedns"] = &TFPlugin{
		Name:    "azureprivatedns",
		Exec:    azurePrivateDNSProvider,
		Version: GetAzureprivatednsVersion(),
	}
}
