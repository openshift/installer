package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: kubevirt.Provider})
	}
	KnownPlugins["terraform-provider-kubevirt"] = &TFPlugin{
		Name:    "kubevirt",
		Exec:    exec,
		Version: GetKubevirtVersion(),
	}
}
