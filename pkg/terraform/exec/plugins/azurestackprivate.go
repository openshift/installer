package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/openshift/installer/pkg/terraform/exec/plugins/azurestackprivate"
)

func init() {
	azurestackPrivateProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: azurestackprivate.Provider,
		})
	}
	KnownPlugins["terraform-provider-azurestackprivate"] = azurestackPrivateProvider
}
