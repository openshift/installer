package infrastructure

import (
	"context"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
)

// Provider defines the interface to be used for provisioning
// and working with cloud infrastructure.
type Provider interface {
	// Provision creates the infrastructure resources for the stage.
	// ctx: parent context
	// dir: the path of the install dir
	// parents: the parent assets, which can be used to obtain any cluser asset dependencies
	// returns a slice of File assets, which will be appended to the cluster asset file list.
	Provision(ctx context.Context, dir string, parents asset.Parents) ([]*asset.File, error)

	// DestroyBootstrap destroys the temporary bootstrap resources.
	DestroyBootstrap(ctx context.Context, dir string) error

	// ExtractHostAddresses extracts the IPs of the bootstrap and control plane machines.
	ExtractHostAddresses(dir string, config *types.InstallConfig, ha *HostAddresses) error
}

// HostAddresses contains the node addresses & ports to be
// used for gather bootsrap debug logs.
type HostAddresses struct {
	Bootstrap string
	Masters   []string
	Port      int
}
