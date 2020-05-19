// Package baremetal contains utilities that help gather Baremetal specific
// information from terraform state.
package baremetal

import (
	"net"

	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/types"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

// ControlPlaneIPs returns the ip addresses for control plane hosts, retrieved via Ironic introspection data. We prefer
// to get an IP in a machine network, but will fallback to any valid IP returned by introspection data.
func ControlPlaneIPs(config *types.InstallConfig, tfs *terraform.State) ([]string, error) {
	mrs, err := terraform.LookupResource(tfs, "module.masters", "ironic_introspection", "openshift-master-introspection")
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup masters introspection data")
	}

	var errs []error
	var masters []string

	for idx, inst := range mrs.Instances {
		interfaces, _, err := unstructured.NestedSlice(inst.Attributes, "interfaces")
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "could not get interfaces for master-%d", idx))
			continue
		}

		masterIP := ""
		var ips []string

		// Look at all interfaces -- if we find an IP in one of the machine networks, we've got the best IP. Otherwise,
		// collect all the IP's and pick the first found.
		for _, iface := range interfaces {
			ipString, _, err := unstructured.NestedString(iface.(map[string]interface{}), "ip")
			if err != nil {
				continue
			}

			if ip := net.ParseIP(ipString); ip != nil {
				for _, network := range config.MachineNetwork {
					if network.CIDR.Contains(ip) {
						masterIP = ipString
						break
					} else {
						ips = append(ips, ipString)
					}
				}
			}
		}

		if masterIP != "" {
			masters = append(masters, masterIP)
		} else if len(ips) > 0 {
			masters = append(masters, ips[0])
		} else {
			errs = append(errs, errors.Wrapf(err, "could not get ip for master-%d", idx))
		}
	}

	return masters, utilerrors.NewAggregate(errs)
}
