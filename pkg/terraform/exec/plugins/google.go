package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-provider-google/google"
)

func init() {
	googleProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: google.Provider,
		})
	}
	KnownPlugins["terraform-provider-google"] = googleProvider
}
