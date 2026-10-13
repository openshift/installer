package baremetal

import (
	"context"
	"fmt"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/types"
)

// Provider is the baremetal platform provider.
type Provider struct{}

// InitializeProvider initializes an empty Provider.
func InitializeProvider() infrastructure.Provider {
	return Provider{}
}

// Provision creates a baremetal platform bootstrap node.
func (a Provider) Provision(ctx context.Context, dir string, parents asset.Parents) ([]*asset.File, error) {
	config, err := getConfig(dir)
	if err != nil {
		return []*asset.File{}, fmt.Errorf("failed to get baremetal platform config: %w", err)
	}

	err = createBootstrap(config)
	if err != nil {
		return []*asset.File{}, fmt.Errorf("failed to create bootstrap: %w", err)
	}

	return []*asset.File{}, nil
}

// DestroyBootstrap destroys the temporary bootstrap resources.
func (a Provider) DestroyBootstrap(ctx context.Context, dir string) error {
	config, err := getConfig(dir)
	if err != nil {
		return fmt.Errorf("failed to get baremetal platform config: %w", err)
	}
	err = destroyBootstrap(config)
	if err != nil {
		return fmt.Errorf("failed to create bootstrap: %w", err)
	}

	return nil
}

// ExtractHostAddresses extracts the IPs of the bootstrap and control plane machines.
func (a Provider) ExtractHostAddresses(dir string, ic *types.InstallConfig, ha *infrastructure.HostAddresses) error {
	ha.Bootstrap = ic.Platform.BareMetal.BootstrapProvisioningIP

	masters, err := getMasterAddresses(dir)
	if err != nil {
		return err
	}

	ha.Masters = masters

	return nil
}
