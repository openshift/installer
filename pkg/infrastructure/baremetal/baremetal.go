package baremetal

import (
	"context"
	"fmt"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/types"
)

type Provider struct{}

func InitializeProvider() infrastructure.Provider {
	return Provider{}
}

// Provision creates a baremetal platform bootstrap node.
func (a Provider) Provision(ctx context.Context, dir string, parents asset.Parents) ([]*asset.File, error) {
	config, err := GetConfig(dir)
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
	config, err := GetConfig(dir)
	if err != nil {
		return fmt.Errorf("failed to get baremetal platform config: %v", err)
	}
	err = destroyBootstrap(config)
	if err != nil {
		return fmt.Errorf("failed to create bootstrap: %v", err)
	}

	return nil
}

func (a Provider) ExtractHostAddresses(dir string, ic *types.InstallConfig, ha *infrastructure.HostAddresses) error {
	ha.Bootstrap = ic.Platform.BareMetal.BootstrapProvisioningIP

	masters, err := getMasterAddresses(dir)
	if err != nil {
		return err
	}

	ha.Masters = masters

	return nil
}
