package manifests

import (
	"context"
	"fmt"
	"reflect"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/assisted-service/pkg/staticnetworkconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func getNMStateConfig() (aiv1beta1.NMStateConfig, error) {
	var nmStateConfig aiv1beta1.NMStateConfig
	if err := GetFileData("nmstateconfig.yaml", &nmStateConfig); err != nil {
		// This file is optional if DHCP addressing is used
		return nmStateConfig, err
	}

	return nmStateConfig, nil
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

func ProcessNMStateConfig(infraEnv aiv1beta1.InfraEnv) ([]*models.HostStaticNetworkConfig, error) {

	nmStateConfig, err := getNMStateConfig()

	if err != nil {
		fmt.Println("A nmstateconfig file has not been provided, no static addresses set in iso")
		return nil, nil
	}

	err = validateNMStateConfigAndInfraEnv(nmStateConfig, infraEnv)
	if err != nil {
		return nil, err
	}

	var staticNetworkConfig []*models.HostStaticNetworkConfig
	staticNetworkConfig = append(staticNetworkConfig, &models.HostStaticNetworkConfig{
		MacInterfaceMap: buildMacInterfaceMap(nmStateConfig),
		NetworkYaml:     string(nmStateConfig.Spec.NetConfig.Raw),
	})
	return staticNetworkConfig, nil
}
