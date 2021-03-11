package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: vsphere.Provider,
		})
	}
	KnownPlugins["terraform-provider-vsphere"] = &TFPlugin{
		Name:    "vsphere",
		Exec:    exec,
		Version: GetVsphereVersion(),
	}
}
