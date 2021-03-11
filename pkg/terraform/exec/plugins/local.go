package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-local/local"
)

func init() {
	localProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: local.Provider,
		})
	}
	KnownPlugins["terraform-provider-local"] = &TFPlugin{
		Name:    "local",
		Exec:    localProvider,
		Version: GetLocalVersion(),
	}
}
