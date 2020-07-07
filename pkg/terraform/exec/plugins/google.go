package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-google-beta/google-beta"
)

func init() {
	googleProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: google.Provider,
		})
	}
	KnownPlugins["terraform-provider-google"] = googleProvider
}
