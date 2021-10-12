package plugins

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: alicloud.Provider,
		})
	}
	KnownPlugins["terraform-provider-alicloud"] = exec
}
