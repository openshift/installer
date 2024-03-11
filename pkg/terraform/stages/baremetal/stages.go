package baremetal

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	"github.com/openshift/installer/pkg/types"
)

// PlatformStages are the stages to run to provision the infrastructure in
// Bare Metal.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		"baremetal",
		"bootstrap",
		[]providers.Provider{providers.Libvirt},
		stages.WithNormalBootstrapDestroy(),
	),
}

func extractOutputHostAddresses(s stages.SplitStage, directory string, config *types.InstallConfig) (bootstrap string, port int, masters []string, err error) {
	port = 22
	bootstrap = config.Platform.BareMetal.BootstrapProvisioningIP

	// masters.tfvars.json
	// 1:{"control_plane_interfaces":[[{"ip":"fd00:1101::47ee:13f6:d3bd:baba","mac":"00:8c:53:d4:b5:2a","name":"enp1s0"},{"ip":"fd2e:6f44:5dd8:c956::14","mac":"00:8c:53:d4:b5:2c","name":"enp2s0"}],[{"ip":"fd00:1101::134c:e458:6e20:f0c2","mac":"00:8c:53:d4:b5:2e","name":"enp1s0"},{"ip":"fd2e:6f44:5dd8:c956::15","mac":"00:8c:53:d4:b5:30","name":"enp2s0"}],[{"ip":"fd00:1101::deab:b4ba:884:26ca","mac":"00:8c:53:d4:b5:32","name":"enp1s0"},{"ip":"fd2e:6f44:5dd8:c956::16","mac":"00:8c:53:d4:b5:34","name":"enp2s0"}]]}

	outputsFilePath := filepath.Join(directory, s.OutputsFilename())
	if _, err := os.Stat(outputsFilePath); err != nil {
		return "", 0, nil, errors.Wrapf(err, "could not find outputs file %q", outputsFilePath)
	}

	outputsFile, err := os.ReadFile(outputsFilePath)
	if err != nil {
		return "", 0, nil, errors.Wrapf(err, "failed to read outputs file %q", outputsFilePath)
	}

	outputs := map[string]interface{}{}
	if err := json.Unmarshal(outputsFile, &outputs); err != nil {
		return "", 0, nil, errors.Wrapf(err, "could not unmarshal outputs file %q", outputsFilePath)
	}

	// control_plane_interfaces are the interfaces for the control plane
	// hosts, retrieved via Ironic introspection data. We prefer to get an
	// IP in a machine network, but will fallback to any valid IP returned
	// by the introspection data.
	if ifacesRaw, ok := outputs["control_plane_interfaces"]; ok {
		ifacesSlice, ok := ifacesRaw.([]interface{})
		if !ok {
			return "", 0, nil, errors.Wrapf(err, "could not read control plane interfaces from outputs file %q", outputsFilePath)
		}

		for idx, ifaceRaw := range ifacesSlice {
			ifaceSliceRaw, ok := ifaceRaw.([]interface{})
			if !ok {
				return "", 0, nil, errors.Wrapf(err, "could not unmarshal raw interface slice")
			}

			masterIP := ""
			var ips []string

			for _, s := range ifaceSliceRaw {
				m, ok := s.(map[string]interface{})
				if !ok {
					return "", 0, nil, errors.Wrapf(err, "could not unmarshal interface")
				}

				ipString, ok := m["ip"].(string)
				if !ok {
					continue
				}

				// Look at all interfaces -- if we find an IP in one of the
				// machine networks, we've got the best IP. Otherwise,
				// collect all the IP's and pick the first found.
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
				return "", 0, nil, errors.Wrapf(err, "could not get ip for master-%d", idx)
			}
		}

	}

	return bootstrap, port, masters, nil
}
