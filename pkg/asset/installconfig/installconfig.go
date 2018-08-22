package installconfig

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/pborman/uuid"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
)

// InstallConfig generates the install-config.yml file.
type installConfig struct {
	assetStock  Stock
	directory   string
	inputReader *bufio.Reader
}

var _ asset.Asset = (*installConfig)(nil)

// Dependencies returns all of the dependencies directly needed by an
// installConfig asset.
func (a *installConfig) Dependencies() []asset.Asset {
	return []asset.Asset{
		a.assetStock.EmailAddress(),
		a.assetStock.Password(),
		a.assetStock.BaseDomain(),
		a.assetStock.ClusterName(),
		a.assetStock.License(),
		a.assetStock.PullSecret(),
		a.assetStock.Platform(),
	}
}

// Generate generates the install-config.yml file.
func (a *installConfig) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	clusterName := string(dependencies[a.assetStock.ClusterName()].Contents[0].Data)
	platform := string(dependencies[a.assetStock.Platform()].Contents[0].Data)

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

	if err := os.MkdirAll(a.directory, 0755); err != nil {
		return nil, err
	}
	assetPath := filepath.Join(a.directory, "install-config.yml")
	if err := ioutil.WriteFile(assetPath, data, 0644); err != nil {
		return nil, err
	}

	return &asset.State{
		Contents: []asset.Content{
			{
				Name: assetPath,
				Data: data,
			},
		},
	}, nil
}
