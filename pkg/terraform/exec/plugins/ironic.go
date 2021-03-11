// +build baremetal

package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/openshift-metal3/terraform-provider-ironic/ironic"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: ironic.Provider,
		})
	}
	KnownPlugins["terraform-provider-ironic"] = &TFPlugin{
		Name:    "ironic",
		Exec:    exec,
		Version: GetIronicVersion(),
	}
}
