package libvirt

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

	// PlatformStages are the stages to run to provision the infrastructure in libvirt.
	var platformStages = []infrastructure.Stage{
		stages.NewStage(
			"libvirt",
			"cluster",
			installDir,
			terraformDir,
			[]providers.Provider{providers.Libvirt},
		),
		stages.NewStage(
			"libvirt",
			"bootstrap",
			installDir,
			terraformDir,
			[]providers.Provider{providers.Libvirt},
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
