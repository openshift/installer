// +build baremetal

package validation

import (
	"fmt"
	"github.com/libvirt/libvirt-go"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"strings"
)

func init() {
	dynamicValidators = append(dynamicValidators, validateInterfaces)
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

	if err := findInterface(p.ProvisioningBridge); err != nil {
		errorList = append(errorList, field.Invalid(fldPath.Child("provisioningBridge"), p.ProvisioningBridge, err.Error()))
	}

	return errorList
}

// interfaceValidator fetches the valid interface names from a particular libvirt instance, and returns a closure
// to validate if an interface is found among them
func interfaceValidator(libvirtURI string) (func(string) error, error) {
	// Connect to libvirt and obtain a list of interface names
	conn, err := libvirt.NewConnect(libvirtURI)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to libvirt")
	}

	interfaces, err := conn.ListAllInterfaces(libvirt.CONNECT_LIST_INTERFACES_ACTIVE)
	if err != nil {
		return nil, errors.Wrap(err, "could not list libvirt interfaces")
	}

	interfaceNames := make([]string, len(interfaces))
	for idx, iface := range interfaces {
		iface, err := iface.GetName()
		if err == nil {
			interfaceNames[idx] = iface
		} else {
			return nil, errors.Wrap(err, "could not get interface name from libvirt")
		}
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
