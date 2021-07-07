package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	// "github.com/aliyun/terraform-provider-alicloud/alicloud"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			// TODO AlibabaCloud: There is a multi-version dependency problem with k8s.io/client-go v11.0.0+incompatible
			// ProviderFunc: alicloud.Provider,
		})
	}
	KnownPlugins["terraform-provider-alicloud"] = exec
}
