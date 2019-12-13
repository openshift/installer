package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: vsphere.Provider,
		})
	}
	KnownPlugins["terraform-provider-vsphere"] = exec
}
