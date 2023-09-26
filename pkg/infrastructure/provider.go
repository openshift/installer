package infrastructure

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
)

// Stage is an individual stage of infrastructure provisioning.
type Stage interface {
	// Name is the name of the stage.
	Name() string

	// Provision creates the infrastructure resources for the stage.
	// tfVars: the terraform variables files
	// fileList: the file list for the cluster asset.
	// returns tfvars output file and a state file
	Provision(tfVars, fileList []*asset.File) (*asset.File, *asset.File, error)

	// DestroyWithBootstrap is true if the stage should be destroyed when destroying the bootstrap resources.
	DestroyWithBootstrap() bool

	// Destroy destroys the resources created in the stage. This should only be called if the stage should be destroyed
	// when destroying the bootstrap resources.
	Destroy(directory string, varFiles []string) error

	// ExtractHostAddresses extracts the IPs of the bootstrap and control plane machines.
	ExtractHostAddresses(directory string, config *types.InstallConfig) (bootstrap string, port int, masters []string, err error)
}
