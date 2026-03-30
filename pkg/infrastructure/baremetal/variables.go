package baremetal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	baremetalhost "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/tfvars"
	baremetaltfvars "github.com/openshift/installer/pkg/tfvars/baremetal"
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

func getMasterAddresses(dir string) ([]string, error) {
	logrus.Debug("baremetal: getting master addresses")
	masters := []string{}

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
			logrus.Debug("    HardwareDetails nil, skipping")
			continue
		}

		for _, nic := range bmh.Status.HardwareDetails.NIC {
			masters = append(masters, nic.IP)
		}
	}

	return masters, nil
}
