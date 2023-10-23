package infrastructure

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
)

// Provider defines the interface to be used for provisioning
// and working with cloud infrastructure.
type Provider interface {
	// Provision creates the infrastructure resources for the stage.
	// dir: the path of the install dir
	// vars: cluster configuration input variables, such as terraform variables files
	// returns a slice of File assets, which will be appended to the cluster asset file list.
	Provision(dir string, vars []*asset.File) ([]*asset.File, error)

	// DestroyBootstrap destroys the temporary bootstrap resources.
	DestroyBootstrap(dir string) error

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
