package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			// TODO AlibabaCloud: There is a multi-version dependency problem with k8s.io/client-go v11.0.0+incompatible, future support.
			// ProviderFunc: alicloud.Provider,
		})
	}
	KnownPlugins["terraform-provider-alicloud"] = exec
}
