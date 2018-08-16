package asset

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/pborman/uuid"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/types"
)

// InstallConfig generates the install-config.yml file.
type InstallConfig struct {
	assetStock *Stock
}

var _ Asset = (*InstallConfig)(nil)

// Dependencies returns all of the dependencies directly needed by an
// InstallConfig asset.
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

// Generate generates the install-config.yml file.
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

	assetPath := filepath.Join(a.assetStock.directory, "install-config.yml")
	a.assetStock.createAssetDirectory()
	if err := ioutil.WriteFile(assetPath, data, 0644); err != nil {
		return nil, err
	}

	return &State{
		Contents: []Content{
			{
				Name: assetPath,
				Data: data,
			},
		},
	}, nil
}
