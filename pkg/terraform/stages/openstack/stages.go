package openstack

import (
	"os"

	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	"github.com/pkg/errors"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	terraformDir, err := terraform.Initialize(installDir)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error initializing terraform")
	}

	// PlatformStages are the stages to run to provision the infrastructure in
	// OpenStack.
	var platformStages = []infrastructure.Stage{
		stages.NewStage(
			"openstack",
			"masters",
			installDir,
			terraformDir,
			[]providers.Provider{providers.OpenStack, providers.Ignition},
		),
		stages.NewStage(
			"openstack",
			"bootstrap",
			installDir,
			terraformDir,
			[]providers.Provider{providers.OpenStack},
			stages.WithNormalBootstrapDestroy(),
		),
	}

	// It would be nice to not need to repeat this for each platform but at this stage
	// Perfect is the enemy of good
	terraform.UnpackTerraform(terraformDir, platformStages)

	cleanup := func() error {
		return os.RemoveAll(terraformDir)
	}

	return platformStages, cleanup, nil
}
