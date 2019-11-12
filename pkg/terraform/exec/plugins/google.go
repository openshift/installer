package plugins

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-google/v2/google"
)

func init() {
	googleProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: google.Provider,
		})
	}
	KnownPlugins["terraform-provider-google"] = googleProvider
}
