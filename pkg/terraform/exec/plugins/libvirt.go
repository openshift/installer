// +build libvirt

package plugins

import (
	"github.com/dmacvicar/terraform-provider-libvirt/libvirt"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	libvirttypes "github.com/openshift/installer/pkg/types/libvirt"
)

func init() {
	libvirtProvider := func() {
		defer libvirt.CleanupLibvirtConnections()

		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: libvirt.Provider,
		})
	}
	KnownPlugins["terraform-provider-libvirt"] = &TFPlugin{
		Name:      libvirttypes.Name,
		Exec:      libvirtProvider,
		Resources: []string{"compat"},
	}
}
