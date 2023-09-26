package nutanix

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

	// PlatformStages are the stages to run to provision the infrastructure in Nutanix.
	var platformStages = []infrastructure.Stage{
		stages.NewStage("nutanix", "cluster", installDir, terraformDir, []providers.Provider{providers.Nutanix, providers.Ignition}),
		stages.NewStage("nutanix", "bootstrap", installDir, terraformDir, []providers.Provider{providers.Nutanix}, stages.WithNormalBootstrapDestroy()),
	}

	// It would be nice to not need to repeat this for each platform but at this stage
	// Perfect is the enemy of good
	terraform.UnpackTerraform(terraformDir, platformStages)

	cleanup := func() error {
		return os.RemoveAll(terraformDir)
	}

	return platformStages, cleanup, nil
}
