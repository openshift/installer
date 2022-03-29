package manifests

import (
	"fmt"
	"os"
	"reflect"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	"github.com/pkg/errors"
)

func getNMStateConfig() aiv1beta1.NMStateConfig {
	var nmStateConfig aiv1beta1.NMStateConfig
	if err := GetFileData("nmstateconfig.yaml", &nmStateConfig); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return nmStateConfig
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

func ProcessNMStateConfig(infraEnv aiv1beta1.InfraEnv) ([]*models.HostStaticNetworkConfig, error) {

	nmStateConfig := getNMStateConfig()

	err := validateNMStateConfigAndInfraEnv(nmStateConfig, infraEnv)
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
