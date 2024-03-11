//go:build baremetal
// +build baremetal

package validation

import (
	"fmt"
	"net/url"
	"strings"

	libvirt "github.com/digitalocean/go-libvirt"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/validate"
)

func init() {
	dynamicProvisioningValidators = append(dynamicProvisioningValidators, validateInterfaces)
}

// validateInterfaces ensures that any interfaces required by the platform exist on the libvirt host.
func validateInterfaces(p *baremetal.Platform, fldPath *field.Path) field.ErrorList {
	errorList := field.ErrorList{}

	findInterface, err := interfaceValidator(p.LibvirtURI)
	if err != nil {
		errorList = append(errorList, field.InternalError(fldPath.Child("libvirtURI"), err))
		return errorList
	}

	if err := findInterface(p.ExternalBridge); err != nil {
		errorList = append(errorList, field.Invalid(fldPath.Child("externalBridge"), p.ExternalBridge, err.Error()))
	}

	if err := validate.MAC(p.ExternalMACAddress); p.ExternalMACAddress != "" && err != nil {
		errorList = append(errorList, field.Invalid(fldPath.Child("externalMACAddress"), p.ExternalMACAddress, err.Error()))
	}

	if err := findInterface(p.ProvisioningBridge); p.ProvisioningNetwork != baremetal.DisabledProvisioningNetwork && err != nil {
		errorList = append(errorList, field.Invalid(fldPath.Child("provisioningBridge"), p.ProvisioningBridge, err.Error()))
	}

	if err := validate.MAC(p.ProvisioningMACAddress); p.ProvisioningMACAddress != "" && err != nil {
		errorList = append(errorList, field.Invalid(fldPath.Child("provisioningMACAddress"), p.ProvisioningMACAddress, err.Error()))
	}

	if p.ProvisioningMACAddress != "" && strings.EqualFold(p.ProvisioningMACAddress, p.ExternalMACAddress) {
		errorList = append(errorList, field.Duplicate(fldPath.Child("provisioningMACAddress"), "provisioning and external MAC addresses may not be identical"))
	}

	return errorList
}

// interfaceValidator fetches the valid interface names from a particular libvirt instance, and returns a closure
// to validate if an interface is found among them

func interfaceValidator(libvirtURI string) (func(string) error, error) {
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
	defer virt.Disconnect()

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
