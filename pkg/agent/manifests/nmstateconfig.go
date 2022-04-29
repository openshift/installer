package manifests

import (
	"context"
	"fmt"
	"net"
	"reflect"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/assisted-service/pkg/staticnetworkconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// NMConfig represents a set of configs from which the Node0 IP can be obtained
type NMConfig struct {
	nmConfig func() ([]aiv1beta1.NMStateConfig, error)
}

// NewNMConfig creates a new NMConfig
func NewNMConfig() NMConfig {
	return NMConfig{nmConfig: getNMStateConfig}
}

type nmStateConfig struct {
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
}

type nmStateConfigYamlDecoder int

func (d *nmStateConfigYamlDecoder) newDecodedYaml(yamlDecoder *yaml.YAMLToJSONDecoder) (interface{}, error) {
	decodedData := new(aiv1beta1.NMStateConfig)
	err := yamlDecoder.Decode(&decodedData)

	return decodedData, err
}

// Get a list of nmStateConfig objects from the manifest file
func getNMStateConfig() ([]aiv1beta1.NMStateConfig, error) {
	var decoder nmStateConfigYamlDecoder
	yamlList, err := getFileMultipleYamls("nmstateconfig.yaml", &decoder)

	var nmStateConfigList []aiv1beta1.NMStateConfig
	for i := range yamlList {
		nmStateConfigList = append(nmStateConfigList, *yamlList[i].(*aiv1beta1.NMStateConfig))
	}

	if err != nil {
		err = fmt.Errorf("error reading nmstateconfig file %w", err)
		return nil, err
	}

	return nmStateConfigList, nil
}

func validateNMStateConfigAndInfraEnv(nmStateConfig aiv1beta1.NMStateConfig, infraEnv aiv1beta1.InfraEnv) error {
	if len(nmStateConfig.ObjectMeta.Labels) == 0 {
		return errors.Errorf("NMStateConfig does not have any labels set")
	}

	if len(infraEnv.Spec.NMStateConfigLabelSelector.MatchLabels) == 0 {
		return errors.Errorf("Infra env does not have any labels set with NMStateConfigLabelSelector.MatchLabels")
	}

	if !reflect.DeepEqual(infraEnv.Spec.NMStateConfigLabelSelector.MatchLabels, nmStateConfig.ObjectMeta.Labels) {
		return errors.Errorf("Infra env and NMStateConfig labels do not match")
	}

	return nil
}

func buildMacInterfaceMap(nmStateConfig aiv1beta1.NMStateConfig) models.MacInterfaceMap {
	macInterfaceMap := make(models.MacInterfaceMap, 0, len(nmStateConfig.Spec.Interfaces))
	for _, cfg := range nmStateConfig.Spec.Interfaces {
		logrus.Println("adding MAC interface map to host static network config - Name: ", cfg.Name, " MacAddress:", cfg.MacAddress)
		macInterfaceMap = append(macInterfaceMap, &models.MacInterfaceMapItems0{
			MacAddress:     cfg.MacAddress,
			LogicalNicName: cfg.Name,
		})
	}
	return macInterfaceMap
}

// GetNMIgnitionFiles returns the list of NetworkManager configuration files
func GetNMIgnitionFiles(staticNetworkConfig []*models.HostStaticNetworkConfig) ([]staticnetworkconfig.StaticNetworkConfigData, error) {
	log := logrus.New()
	staticNetworkConfigGenerator := staticnetworkconfig.New(log.WithField("pkg", "manifests"), staticnetworkconfig.Config{MaxConcurrentGenerations: 2})

	// Validate the network config
	if err := staticNetworkConfigGenerator.ValidateStaticConfigParams(context.Background(), staticNetworkConfig); err != nil {
		err = fmt.Errorf("staticNetwork configuration is not valid: %w", err)
		return nil, err
	}

	networkConfigStr, err := staticNetworkConfigGenerator.FormatStaticNetworkConfigForDB(staticNetworkConfig)
	if err != nil {
		err = fmt.Errorf("error marshalling StaticNetwork configuration: %w", err)
		return nil, err
	}

	filesList, err := staticNetworkConfigGenerator.GenerateStaticNetworkConfigData(context.Background(), networkConfigStr)
	if err != nil {
		err = fmt.Errorf("failed to create StaticNetwork config data: %w", err)
		return nil, err
	}

	return filesList, err
}

func processNMStateConfig(infraEnv aiv1beta1.InfraEnv) ([]*models.HostStaticNetworkConfig, error) {

	nmStateConfigList, err := getNMStateConfig()

	if err != nil {
		err = fmt.Errorf("error with nmstateconfig file: %w", err)
		return nil, err
	}

	var staticNetworkConfig []*models.HostStaticNetworkConfig
	for _, nmStateConfig := range nmStateConfigList {

		err = validateNMStateConfigAndInfraEnv(nmStateConfig, infraEnv)
		if err != nil {
			return nil, err
		}

		staticNetworkConfig = append(staticNetworkConfig, &models.HostStaticNetworkConfig{
			MacInterfaceMap: buildMacInterfaceMap(nmStateConfig),
			NetworkYaml:     string(nmStateConfig.Spec.NetConfig.Raw),
		})
	}
	return staticNetworkConfig, nil
}

// GetNodeZeroIP retrieves the first IP from the user provided nmStateConfig yaml file to set as node0 IP
func (n NMConfig) GetNodeZeroIP() string {
	configList, err := n.nmConfig()
	if err != nil {
		panic(err)
	}

	var nmStateConfig nmStateConfig
	// Use entry for first host
	err = yaml.Unmarshal(configList[0].Spec.NetConfig.Raw, &nmStateConfig)
	if err != nil {
		panic(err)
	}

	var nodeZeroIP string
	if nmStateConfig.Interfaces != nil {
		if nmStateConfig.Interfaces[0].IPV4.Address != nil {
			nodeZeroIP = nmStateConfig.Interfaces[0].IPV4.Address[0].IP
		}
		if nmStateConfig.Interfaces[0].IPV6.Address != nil {
			nodeZeroIP = nmStateConfig.Interfaces[0].IPV6.Address[0].IP

		}
		if net.ParseIP(nodeZeroIP) == nil {
			panic("invalid YAML - NMStateconfig")
		}
	} else {
		panic("invalid YAML - NMStateconfig: No valid interfaces set")
	}

	return nodeZeroIP
}
