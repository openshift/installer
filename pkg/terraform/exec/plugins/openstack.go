package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-openstack/openstack"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: openstack.Provider,
		})
	}
	KnownPlugins["terraform-provider-openstack"] = exec
}
