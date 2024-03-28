package baremetal

import (
	"fmt"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/types"
)

type Provider struct{}

func InitializeProvider() infrastructure.Provider {
	return Provider{}
}

func (a Provider) Provision(dir string, parents asset.Parents) ([]*asset.File, error) {
	config, err := GetConfig(parents)
	if err != nil {
		return []*asset.File{}, fmt.Errorf("failed to get baremetal platform config: %v", err)
	}

	err = createBootstrap(config)
	if err != nil {
		return []*asset.File{}, fmt.Errorf("failed to create bootstrap: %v", err)
	}

	return []*asset.File{}, nil
}

func (a Provider) DestroyBootstrap(dir string) error {
	return nil
}

func (a Provider) ExtractHostAddresses(dir string, ic *types.InstallConfig, ha *infrastructure.HostAddresses) error {
	return nil
}
