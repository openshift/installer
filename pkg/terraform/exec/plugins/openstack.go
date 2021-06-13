package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-provider-openstack/terraform-provider-openstack/openstack"

	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
)

func init() {
	openstackProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: openstack.Provider,
		})
	}
	KnownPlugins["terraform-provider-openstack"] = &TFPlugin{
		Name: openstacktypes.Name,
		Exec: openstackProvider,
	}
}
