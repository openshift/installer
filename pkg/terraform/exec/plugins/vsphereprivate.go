package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/openshift/installer/pkg/terraform/exec/plugins/vsphereprivate"
)

func init() {
	vspherePrivateProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: vsphereprivate.Provider})
	}
	KnownPlugins["terraform-provider-vsphereprivate"] = vspherePrivateProvider
}
