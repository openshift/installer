package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: alicloud.Provider,
		})
	}
	KnownPlugins["terraform-provider-alicloud"] = exec
}
