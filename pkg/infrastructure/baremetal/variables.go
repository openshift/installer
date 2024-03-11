package baremetal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	baremetalhost "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/openshift/installer/pkg/tfvars"
	baremetaltfvars "github.com/openshift/installer/pkg/tfvars/baremetal"
	"github.com/sirupsen/logrus"
)

const (
	tfVarsFileName         = "terraform.tfvars.json"
	tfPlatformVarsFileName = "terraform.platform.auto.tfvars.json"
	MastersFileName        = ".masters.json"
)

type bridge struct {
	Name string
	MAC  string
}

type baremetalConfig struct {
	ClusterID         string
	BootstrapOSImage  string
	IgnitionBootstrap string
	LibvirtURI        string
	Bridges           []bridge
}

func GetConfig(dir string) (baremetalConfig, error) {
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

	config.BootstrapOSImage = clusterBaremetalConfig.BootstrapOSImage
	config.LibvirtURI = clusterBaremetalConfig.LibvirtURI

	for _, bridgeMap := range clusterBaremetalConfig.Bridges {
		mac, ok := bridgeMap["mac"]
		if !ok {
			return config, fmt.Errorf("bridge is missng a MAC address")
		}

		name, ok := bridgeMap["name"]
		if !ok {
			return config, fmt.Errorf("bridge is missng a name")
		}

		b := bridge{
			Name: name,
			MAC:  mac,
		}

		config.Bridges = append(config.Bridges, b)
	}

	return config, nil

}

func getMasterAddresses(dir string) ([]string, error) {
	logrus.Debug("baremetal: getting master addresses")
	masters := []string{}

	data, err := os.ReadFile(filepath.Join(dir, MastersFileName))
	if err != nil {
		return masters, err
	}

	hosts := map[string]baremetalhost.BareMetalHost{}

	err = json.Unmarshal(data, &hosts)
	if err != nil {
		return masters, err
	}

	for _, bmh := range hosts {
		for _, nic := range bmh.Status.HardwareDetails.NIC {
			masters = append(masters, nic.IP)
		}
	}

	return masters, nil
}
