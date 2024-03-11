package baremetal

import (
	"encoding/json"
	"fmt"
	"github.com/openshift/installer/pkg/asset"
	tfvarsAsset "github.com/openshift/installer/pkg/asset/cluster/tfvars"
	"github.com/openshift/installer/pkg/tfvars"
	baremetaltfvars "github.com/openshift/installer/pkg/tfvars/baremetal"
)

const (
	tfVarsFileName         = "terraform.tfvars.json"
	tfPlatformVarsFileName = "terraform.platform.auto.tfvars.json"
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

func GetConfig(parents asset.Parents) (baremetalConfig, error) {
	config := baremetalConfig{}

	terraformVariables := &tfvarsAsset.TerraformVariables{}
	parents.Get(terraformVariables)

	clusterConfig := &tfvars.Config{}
	clusterBaremetalConfig := &baremetaltfvars.Config{}

	for _, file := range terraformVariables.Files() {
		switch file.Filename {
		case tfVarsFileName:
			if err := json.Unmarshal(file.Data, clusterConfig); err != nil {
				return config, err
			}
		case tfPlatformVarsFileName:
			if err := json.Unmarshal(file.Data, clusterBaremetalConfig); err != nil {
				return config, err
			}
		}
	}

	config.ClusterID = clusterConfig.ClusterID
	config.BootstrapOSImage = clusterBaremetalConfig.BootstrapOSImage
	config.LibvirtURI = clusterBaremetalConfig.LibvirtURI
	config.IgnitionBootstrap = clusterConfig.IgnitionBootstrap

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
