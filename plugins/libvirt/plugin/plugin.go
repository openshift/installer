package plugin

import (
	"github.com/openshift/installer/plugins"

	"github.com/dmacvicar/terraform-provider-libvirt/libvirt"
	tfplugin "github.com/hashicorp/terraform-plugin-sdk/plugin"
)

type LibvirtPlugin struct {}

// LibvirtPlugin implements the Plugin interface
var _ plugins.Plugin = &LibvirtPlugin{}

func (p *LibvirtPlugin) Init() {
	defer libvirt.CleanupLibvirtConnections()

	tfplugin.Serve(&tfplugin.ServeOpts{
		ProviderFunc: libvirt.Provider,
	})
}
