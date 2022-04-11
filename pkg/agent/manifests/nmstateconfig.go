package manifests

import (
	"context"
	"fmt"

	"io"
	"os"

	yamlV3 "gopkg.in/yaml.v3"

	"github.com/openshift/assisted-service/models"
	"github.com/openshift/assisted-service/pkg/staticnetworkconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type NMStateConfig struct {
	Kind string `yaml:"kind"`
	Spec struct {
		Config struct {
			Interfaces []struct {
				IPV4 struct {
					Address []struct {
						IP string `yaml:"ip,omitempty"`
					} `yaml:"address,omitempty"`
				} `yaml:"ipv4,omitempty"`
				IPV6 struct {
					Address []struct {
						IP string `yaml:"ip,omitempty"`
					} `yaml:"address,omitempty"`
				} `yaml:"ipv6,omitempty"`
			} `yaml:"interfaces,omitempty"`
		} `yaml:"config,omitempty"`
		Interfaces []struct {
			Name       string `yaml:"name,omitempty"`
			MacAddress string `yaml:"macAddress"`
		} `yaml:"interfaces,omitempty"`
	} `yaml:"spec,omitempty"`
}

// Retrieve the all NMStateConfigs from the user provided NMStateConfig yaml file
func getNMStateConfig() []NMStateConfig {
	var nmStateConfig []NMStateConfig
	f, err := os.Open("./manifests/nmstateconfig.yaml")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	d := yamlV3.NewDecoder(f)
	for {
		config := new(NMStateConfig)
		err := d.Decode(&config)
		if config == nil {
			continue
		}
		// break the loop in case of EOF
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		nmStateConfig = append(nmStateConfig, *config)
	}
	return nmStateConfig
}

func buildMacInterfaceMap(nmStateConfig NMStateConfig) models.MacInterfaceMap {
	macInterfaceMap := make(models.MacInterfaceMap, 0, len(nmStateConfig.Spec.Interfaces))
	for _, cfg := range nmStateConfig.Spec.Interfaces {
		// ToDo: Use logging
		// fmt.Println("adding MAC interface map to host static network config - Name: %s, MacAddress: %s ,",
		// 	cfg.Name, cfg.MacAddress)
		macInterfaceMap = append(macInterfaceMap, &models.MacInterfaceMapItems0{
			MacAddress:     cfg.MacAddress,
			LogicalNicName: cfg.Name,
		})
	}
	return macInterfaceMap
}

// Get the NetworkManager configuration files
func GetNMIgnitionFiles(staticNetworkConfig []*models.HostStaticNetworkConfig) ([]staticnetworkconfig.StaticNetworkConfigData, error) {
	log := logrus.New()
	staticNetworkConfigGenerator := staticnetworkconfig.New(log.WithField("pkg", "manifests"), staticnetworkconfig.Config{2})

	// Validate the network config
	if err := staticNetworkConfigGenerator.ValidateStaticConfigParams(context.Background(), staticNetworkConfig); err != nil {
		err = fmt.Errorf("StaticNetwork configuration is not valid: %w", err)
		return nil, err
	}

	networkConfigStr, err := staticNetworkConfigGenerator.FormatStaticNetworkConfigForDB(staticNetworkConfig)
	if err != nil {
		err = fmt.Errorf("Error marshalling StaticNetwork configuration: %w", err)
		return nil, err
	}

	filesList, err := staticNetworkConfigGenerator.GenerateStaticNetworkConfigData(context.Background(), networkConfigStr)
	if err != nil {
		err = fmt.Errorf("Failed to create StaticNetwork config data: %w", err)
		return nil, err
	}

	return filesList, err
}

func GetStaticNetworkConfig() []*models.HostStaticNetworkConfig {
	var staticNetworkConfig []*models.HostStaticNetworkConfig
	nmStateConfigs := getNMStateConfig()
	for _, config := range nmStateConfigs {
		networkYaml, err := yamlV3.Marshal(config.Spec.Config)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		staticNetworkConfig = append(staticNetworkConfig, &models.HostStaticNetworkConfig{
			MacInterfaceMap: buildMacInterfaceMap(config),
			NetworkYaml:     string(networkYaml),
		})

	}

	return staticNetworkConfig
}

// Retrieve the first IP from the user provided NMStateConfig yaml file to set as node0 IP
func GetNodeZeroIP() string {
	config := getNMStateConfig()
	nodeZeroIP := config[0].Spec.Config.Interfaces[0].IPV4.Address[0].IP
	return nodeZeroIP
}
