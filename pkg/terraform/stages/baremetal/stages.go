package baremetal

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	"github.com/openshift/installer/pkg/types"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	terraformDir, err := terraform.Initialize(installDir)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error initializing terraform")
	}

	// PlatformStages are the stages to run to provision the infrastructure in
	// Bare Metal.
	var platformStages = []infrastructure.Stage{
		stages.NewStage(
			"baremetal",
			"bootstrap",
			installDir,
			terraformDir,
			[]providers.Provider{providers.Ironic, providers.Libvirt},
			stages.WithNormalBootstrapDestroy(),
		),
		stages.NewStage(
			"baremetal",
			"masters",
			installDir,
			terraformDir,
			[]providers.Provider{providers.Ironic},
			stages.WithCustomExtractHostAddresses(extractOutputHostAddresses),
		),
	}
	// It would be nice to not need to repeat this for each platform but at this stage
	// Perfect is the enemy of good
	terraform.UnpackTerraform(terraformDir, platformStages)

	cleanup := func() error {
		return os.RemoveAll(terraformDir)
	}

	return platformStages, cleanup, nil
}

func extractOutputHostAddresses(s stages.SplitStage, directory string, config *types.InstallConfig) (bootstrap string, port int, masters []string, err error) {
	port = 22
	bootstrap = config.Platform.BareMetal.BootstrapProvisioningIP

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
