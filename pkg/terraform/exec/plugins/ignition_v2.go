// +build !okd

package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-ignition/ignition"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: ignition.Provider,
		})
	}
	KnownPlugins["terraform-provider-ignition"] = exec
}
