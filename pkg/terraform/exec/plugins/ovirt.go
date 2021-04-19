package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/ovirt/terraform-provider-ovirt/ovirt"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: ovirt.ProviderContext(),
		})
	}
	KnownPlugins["terraform-provider-ovirt"] = exec
}
