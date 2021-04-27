package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/ovirt/terraform-provider-ovirt/ovirt"

	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
)

func init() {
	ovirtProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: ovirt.ProviderContext(),
		})
	}
	KnownPlugins["terraform-provider-ovirt"] = &TFPlugin{
		Name:      ovirttypes.Name,
		Exec:      ovirtProvider,
		Resources: []string{"compat"},
	}
}
