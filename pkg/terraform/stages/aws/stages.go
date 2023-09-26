package aws

import (
	"os"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	"github.com/openshift/installer/pkg/types/aws"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	terraformDir, err := terraform.Initialize(installDir)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error initializing terraform")
	}

	platformStages := []infrastructure.Stage{
		stages.NewStage(
			"aws",
			"cluster",
			installDir,
			terraformDir,
			[]providers.Provider{providers.AWS},
		),
		stages.NewStage(
			"aws",
			"bootstrap",
			installDir,
			terraformDir,
			[]providers.Provider{providers.AWS},
			stages.WithCustomBootstrapDestroy(customBootstrapDestroy),
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

func customBootstrapDestroy(s stages.SplitStage, directory string, terraformDir string, varFiles []string) error {
	opts := make([]tfexec.DestroyOption, 0, len(varFiles)+1)
	for _, varFile := range varFiles {
		opts = append(opts, tfexec.VarFile(varFile))
	}
	// The bootstrap destroy will no longer refresh state. This was added as a change to counteract
	// the upgrade to the aws terraform provider v5.4.0 where the state changes were causing unsupported
	// operation errors when removing security group rules in sc2s regions.
	logrus.Debugf("aws bootstrap destroy stage will not refresh terraform state")
	opts = append(opts, tfexec.Refresh(false))
	return errors.Wrap(
		terraform.Destroy(directory, aws.Name, s, terraformDir, opts...),
		"failed to destroy bootstrap",
	)
}
