package validation

import (
	"fmt"
	"net/url"
	"strings"

	libvirt "github.com/digitalocean/go-libvirt"
	"github.com/sirupsen/logrus"
)

// libvirtInterfaceValidator fetches the valid interface names from a particular libvirt instance, and returns a closure
// to validate if an interface is found among them
func libvirtInterfaceValidator(libvirtURI string) (func(string) error, error) {
	// Connect to libvirt and obtain a list of interface names
	interfaces := make(map[string]struct{})
	var exists = struct{}{}

	uri, err := url.Parse(libvirtURI)
	if err != nil {
		return nil, err
	}

	virt, err := libvirt.ConnectToURI(uri)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := virt.Disconnect(); err != nil {
			logrus.Errorln("error disconnecting from libvirt:", err)
		}
	}()

	networks, _, err := virt.ConnectListAllNetworks(1, libvirt.ConnectListNetworksActive)
	if err != nil {
		return nil, fmt.Errorf("could not list libvirt networks: %w", err)
	}
	for _, network := range networks {
		bridgeName, err := virt.NetworkGetBridgeName(network)
		if err == nil && bridgeName == network.Name {
			interfaces[network.Name] = exists
		}
	}
	bridges, _, err := virt.ConnectListAllInterfaces(1, libvirt.ConnectListInterfacesActive)
	if err != nil {
		return nil, fmt.Errorf("could not list libvirt interfaces: %w", err)
	}

	for _, bridge := range bridges {
		interfaces[bridge.Name] = exists
	}
	interfaceNames := make([]string, len(interfaces))
	idx := 0
	for key := range interfaces {
		interfaceNames[idx] = key
		idx++
	}

	// Return a closure to validate if any particular interface is found among interfaceNames
	return func(interfaceName string) error {
		if len(interfaceNames) == 0 {
			return fmt.Errorf("no interfaces found")
		} else {
			for _, foundInterface := range interfaceNames {
				if foundInterface == interfaceName {
					return nil
				}
			}

			return fmt.Errorf("could not find interface %q, valid interfaces are %s", interfaceName, strings.Join(interfaceNames, ", "))
		}
	}, nil
}
