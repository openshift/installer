package gcp

import (
	"os"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	terraformDir, err := terraform.Initialize(installDir)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error initializing terraform")
	}

	// PlatformStages are the stages to run to provision the infrastructure in GCP.
	var platformStages = []infrastructure.Stage{
		stages.NewStage(
			"gcp",
			"cluster",
			installDir,
			terraformDir,
			[]providers.Provider{providers.Google},
		),
		stages.NewStage(
			"gcp",
			"bootstrap",
			installDir,
			terraformDir,
			[]providers.Provider{providers.Google, providers.Ignition},
			stages.WithNormalBootstrapDestroy(),
		),
		stages.NewStage(
			"gcp",
			"post-bootstrap",
			installDir,
			terraformDir,
			[]providers.Provider{providers.Google},
			stages.WithCustomBootstrapDestroy(removeFromLoadBalancers),
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

func removeFromLoadBalancers(s stages.SplitStage, directory string, terraformDir string, varFiles []string) error {
	opts := make([]tfexec.ApplyOption, 0, len(varFiles)+1)
	for _, varFile := range varFiles {
		opts = append(opts, tfexec.VarFile(varFile))
	}
	opts = append(opts, tfexec.Var("gcp_bootstrap_lb=false"))
	return errors.Wrap(
		terraform.Apply(directory, gcptypes.Name, s, terraformDir, opts...),
		"failed disabling bootstrap load balancing",
	)
}
