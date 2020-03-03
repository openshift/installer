package loader

import (
	"fmt"
	"plugin"

	installerplugins "github.com/openshift/installer/plugins"
)

const libvirtPluginName string = "libvirt-plugin.so"

func LoadPlugin() (installerplugins.Plugin, error) {
	p, err := plugin.Open(libvirtPluginName)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", libvirtPluginName, err)
	}

	factoryFunc, err := p.Lookup("NewPlugin")
	if err != nil {
		return nil, fmt.Errorf("failed to find NewPlugin entrypoint in %s: %v", libvirtPluginName, err)
	}

	return factoryFunc.(func() installerplugins.Plugin)(), nil
}
