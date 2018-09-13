package installconfig

import (
	"fmt"
	"path/filepath"

	"github.com/ghodss/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// installConfig generates the install-config.yml file.
type installConfig struct {
	assetStock Stock
	directory  string
}

var _ asset.Asset = (*installConfig)(nil)

// Dependencies returns all of the dependencies directly needed by an
// installConfig asset.
func (a *installConfig) Dependencies() []asset.Asset {
	return []asset.Asset{
		a.assetStock.ClusterID(),
		a.assetStock.EmailAddress(),
		a.assetStock.Password(),
		a.assetStock.SSHKey(),
		a.assetStock.BaseDomain(),
		a.assetStock.ClusterName(),
		a.assetStock.PullSecret(),
		a.assetStock.Platform(),
	}
}

// Generate generates the install-config.yml file.
func (a *installConfig) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	clusterID := string(dependencies[a.assetStock.ClusterID()].Contents[0].Data)
	emailAddress := string(dependencies[a.assetStock.EmailAddress()].Contents[0].Data)
	password := string(dependencies[a.assetStock.Password()].Contents[0].Data)
	sshKey := string(dependencies[a.assetStock.SSHKey()].Contents[0].Data)
	baseDomain := string(dependencies[a.assetStock.BaseDomain()].Contents[0].Data)
	clusterName := string(dependencies[a.assetStock.ClusterName()].Contents[0].Data)
	pullSecret := string(dependencies[a.assetStock.PullSecret()].Contents[0].Data)

	installConfig := types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterName,
		},
		ClusterName: clusterName,
		ClusterID:   clusterID,
		Admin: types.Admin{
			Email:    emailAddress,
			Password: password,
			SSHKey:   sshKey,
		},
		BaseDomain: baseDomain,
		PullSecret: pullSecret,
	}

	platformState := dependencies[a.assetStock.Platform()]
	platform := string(platformState.Contents[0].Data)
	switch platform {
	case AWSPlatformType:
		region := string(platformState.Contents[1].Data)
		installConfig.AWS = &types.AWSPlatform{
			Region: region,
		}
	case LibvirtPlatformType:
		uri := string(platformState.Contents[1].Data)
		installConfig.Libvirt = &types.LibvirtPlatform{
			URI: uri,
		}
	default:
		return nil, fmt.Errorf("unknown platform type %q", platform)
	}

	data, err := yaml.Marshal(installConfig)
	if err != nil {
		return nil, err
	}

	return &asset.State{
		Contents: []asset.Content{
			{
				Name: filepath.Join(a.directory, "install-config.yml"),
				Data: data,
			},
		},
	}, nil
}

// GetInstallConfig returns the *types.InstallConfig from the parent asset map.
func GetInstallConfig(installConfig asset.Asset, parents map[asset.Asset]*asset.State) (*types.InstallConfig, error) {
	var cfg types.InstallConfig

	st, ok := parents[installConfig]
	if !ok {
		return nil, fmt.Errorf("failed to find %T in parents", installConfig)
	}

	if err := yaml.Unmarshal(st.Contents[0].Data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the installconfig: %v", err)
	}

	return &cfg, nil
}
