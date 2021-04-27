package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere"

	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

func init() {
	vsphereProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: vsphere.Provider,
		})
	}
	KnownPlugins["terraform-provider-vsphere"] = &TFPlugin{
		Name:      vspheretypes.Name,
		Exec:      vsphereProvider,
		Resources: []string{"compat"},
	}
}
