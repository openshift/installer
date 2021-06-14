package terraform

import (
	"github.com/openshift/installer/pkg/types"
)

// Stage is an individual stage of terraform infrastructure provisioning.
type Stage interface {
	// Name is the name of the stage.
	Name() string
	// StateFilename is the name of the terraform state file.
	StateFilename() string
	// OutputsFilename is the name of the outputs file for the stage.
	OutputsFilename() string
	// DestroyWithBootstrap is true if the stage should be destroyed when destroying the bootstrap resources.
	DestroyWithBootstrap() bool
	// Destroy destroys the resources created in the stage. This should only be called if the stage should be destroyed
	// when destroying the bootstrap resources.
	Destroy(directory string, extraArgs []string) error
	// ExtractHostAddresses extracts the IPs of the bootstrap and control plane machines.
	ExtractHostAddresses(directory string, config *types.InstallConfig) (bootstrap string, port int, masters []string, err error)
}
