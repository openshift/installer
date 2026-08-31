package baremetal

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"

	baremetalhost "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/tfvars"
	baremetaltfvars "github.com/openshift/installer/pkg/tfvars/baremetal"
	"github.com/openshift/installer/pkg/types"
)

const (
	tfVarsFileName         = "terraform.tfvars.json"
	tfPlatformVarsFileName = "terraform.platform.auto.tfvars.json"
	// MastersFileName is the file where we store networking data for the control plane.
	MastersFileName = ".masters.json"
)

type baremetalConfig struct {
	ClusterID         string
	IgnitionBootstrap string
	baremetaltfvars.Config
}

func getConfig(dir string) (baremetalConfig, error) {
	config := baremetalConfig{}
	clusterConfig := &tfvars.Config{}
	clusterBaremetalConfig := &baremetaltfvars.Config{}

	data, err := os.ReadFile(filepath.Join(dir, tfVarsFileName))
	if err == nil {
		err = json.Unmarshal(data, clusterConfig)
	}
	if err != nil {
		return config, fmt.Errorf("failed to load cluster terraform variables: %w", err)
	}

	config.ClusterID = clusterConfig.ClusterID
	config.IgnitionBootstrap = clusterConfig.IgnitionBootstrap

	data, err = os.ReadFile(filepath.Join(dir, tfPlatformVarsFileName))
	if err == nil {
		err = json.Unmarshal(data, clusterBaremetalConfig)
	}
	if err != nil {
		return config, fmt.Errorf("failed to load cluster terraform variables: %w", err)
	}

	config.Config = *clusterBaremetalConfig

	return config, nil
}

func getMasterAddresses(dir string, machineNetworks []types.MachineNetworkEntry) ([]string, error) {
	logrus.Debug("baremetal: getting master addresses")
	masters := []string{}
	if len(machineNetworks) == 0 {
		return masters, fmt.Errorf("failed to get master addresses: machine network is not configured")
	}

	primaryIsIPv4 := isIPv4(machineNetworks[0].CIDR.IP)
	primaryNetworks := make([]types.MachineNetworkEntry, 0, len(machineNetworks))
	for _, machineNetwork := range machineNetworks {
		if isIPv4(machineNetwork.CIDR.IP) == primaryIsIPv4 {
			primaryNetworks = append(primaryNetworks, machineNetwork)
		}
	}

	data, err := os.ReadFile(filepath.Join(dir, MastersFileName))
	if err != nil {
		return masters, fmt.Errorf("failed to read masters.json (this can happen when bootstrap didn't run): %w", err)
	}

	hosts := map[string]baremetalhost.BareMetalHost{}

	err = json.Unmarshal(data, &hosts)
	if err != nil {
		return masters, err
	}

	for _, bmh := range hosts {
		logrus.Debug("  bmh:", bmh.Name)

		if bmh.Status.HardwareDetails == nil {
			logrus.Warnf("baremetal: no hardware details found for master %q, skipping", bmh.Name)
			continue
		}

		address, err := addressForMaster(primaryIsIPv4, bmh.Status.HardwareDetails.NIC, primaryNetworks)
		if err != nil {
			logrus.Warnf("baremetal: primary network error %v, %q, skipping", err, bmh.Name)
			continue
		}
		masters = append(masters, address)
	}

	return masters, nil
}

func addressForMaster(primaryIsIPv4 bool, nics []baremetalhost.NIC, machineNetworks []types.MachineNetworkEntry) (string, error) {
	for _, nic := range nics {
		ip := net.ParseIP(nic.IP)
		if ip == nil || ip.IsLinkLocalUnicast() || isIPv4(ip) != primaryIsIPv4 {
			continue
		}

		for _, machineNetwork := range machineNetworks {
			if machineNetwork.CIDR.Contains(ip) {
				return nic.IP, nil
			}
		}
	}
	return "", fmt.Errorf("no primary machine network found")
}

func isIPv4(ip net.IP) bool {
	return ip.To4() != nil
}
