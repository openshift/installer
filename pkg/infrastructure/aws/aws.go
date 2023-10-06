package aws

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/types"

	"github.com/sirupsen/logrus"
)

const clusterTagValue = "owned"

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	awsProvisionStage := []infrastructure.Stage{
		AWSInfraProvider{
			"aws",
			"cluster",
			true,
			normalAWSDestroy,
			normalAWSExtractHostAddresses,
			normalAWSProvision,
		},
	}

	noop := func() error { return nil }
	return awsProvisionStage, noop, nil
}

type AWSInfraProvider struct {
	platform             string
	name                 string
	destroyWithBootstrap bool
	destroy              DestroyFunc
	extractHostAddresses ExtractFunc
	provision            ProvisionFunc
}

// DestroyFunc is a function for destroying the stage.
type DestroyFunc func(a AWSInfraProvider, directory string, varFiles []string) error

// ExtractFunc is a function for extracting host addresses.
type ExtractFunc func(a AWSInfraProvider, directory string, ic *types.InstallConfig) (string, int, []string, error)

// Provision is a function for creating cloud resources.
type ProvisionFunc func(a AWSInfraProvider, tfVars, fileList []*asset.File) (*asset.File, *asset.File, error)

// Name implements pkg/infrastructure/Stage.Name
func (a AWSInfraProvider) Name() string {
	return a.name
}

// Platform implements pkg/infrastructure/Stage.Platform
func (a AWSInfraProvider) Platform() string {
	return a.platform
}

// Provision implements pkg/infrastructure/Stage.Provision
func (a AWSInfraProvider) Provision(tfvars, fileList []*asset.File) (*asset.File, *asset.File, error) {
	return a.provision(a, tfvars, fileList)
}

// DestroyWithBootstrap implements pkg/infrastructure/Stage.DestroyWithBootstrap
func (a AWSInfraProvider) DestroyWithBootstrap() bool {
	return a.destroyWithBootstrap
}

// Destroy implements pkg/infrastructure/Stage.Destroy
func (a AWSInfraProvider) Destroy(directory string, varFiles []string) error {
	return a.destroy(a, directory, varFiles)
}

// ExtractHostAddresses implements pkg/infrastructure/Stage.ExtractHostAddresses
func (a AWSInfraProvider) ExtractHostAddresses(directory string, ic *types.InstallConfig) (string, int, []string, error) {
	return a.extractHostAddresses(a, directory, ic)
}

func normalAWSDestroy(a AWSInfraProvider, directory string, varFiles []string) error {
	//panic("not implemented")
	logrus.Errorln("Pretending to destroy bootstrap resources...")
	return nil
}

func normalAWSExtractHostAddresses(a AWSInfraProvider, directory string, ic *types.InstallConfig) (string, int, []string, error) {
	//panic("not implemented")
	logrus.Errorln("Pretending to return bootstrap host addresses")
	return "", 0, nil, nil
}
