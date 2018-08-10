package asset

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/pborman/uuid"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/types"
)

type InstallConfig struct {
	assetStock *Stock
}

var _ Asset = (*InstallConfig)(nil)

func (a *InstallConfig) Dependencies() []Asset {
	return []Asset{
		a.assetStock.emailAddress,
		a.assetStock.password,
		a.assetStock.baseDomain,
		a.assetStock.clusterName,
		a.assetStock.license,
		a.assetStock.pullSecret,
		a.assetStock.platform,
	}
}

func (a *InstallConfig) Generate(dependencies map[Asset]*State) (*State, error) {
	clusterName := string(dependencies[a.assetStock.clusterName].Contents[0].Data)
	platform := string(dependencies[a.assetStock.platform].Contents[0].Data)

	installConfig := types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterName,
		},
		ClusterID: uuid.NewUUID(),
	}

	switch platform {
	case AWSPlatformType:
		installConfig.AWS = &types.AWSPlatform{}
	case LibvirtPlatformType:
		installConfig.Libvirt = &types.LibvirtPlatform{}
	default:
		return nil, fmt.Errorf("unknown platform type %q", platform)
	}

	data, err := yaml.Marshal(installConfig)
	if err != nil {
		return nil, err
	}

	return &State{
		Contents: []Content{
			{Data: data},
		},
	}, nil
}
