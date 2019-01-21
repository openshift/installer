package plugins

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-google/google"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: google.Provider,
		})
	}
	KnownPlugins["terraform-provider-google"] = exec
}
