package terraform

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/types"
)

// Stage is an individual stage of terraform infrastructure provisioning.
type Stage interface {
	// Name is the name of the stage.
	Name() string

	// Platform is the name of the platform.
	Platform() string

	// StateFilename is the name of the terraform state file.
	StateFilename() string

	// OutputsFilename is the name of the outputs file for the stage.
	OutputsFilename() string

	// Providers is the list of providers that are used for the stage.
	Providers() []providers.Provider

	// DestroyWithBootstrap is true if the stage should be destroyed when destroying the bootstrap resources.
	DestroyWithBootstrap() bool

	// Destroy destroys the resources created in the stage. This should only be called if the stage should be destroyed
	// when destroying the bootstrap resources.
	Destroy(directory string, terraformDir string, varFiles []string) error

	// ExtractHostAddresses extracts the IPs of the bootstrap and control plane machines.
	ExtractHostAddresses(directory string, config *types.InstallConfig) (bootstrap string, port int, masters []string, err error)

	// ExtractLBConfig extracts the LB DNS Names of the internal and external API LBs.
	ExtractLBConfig(directory string, terraformDir string, file *asset.File, tfvarsFile *asset.File) (ignition string, err error)
}
