// +build baremetal

package plugins

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/openshift-metalkube/terraform-provider-ironic/ironic"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: ironic.Provider,
		})
	}
	KnownPlugins["terraform-provider-ironic"] = exec
}
