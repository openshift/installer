package powervs

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-exec/tfexec"

	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	terraformDir, err := terraform.Initialize(installDir)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing terraform: %w", err)
	}

	// PlatformStages are the stages to run to provision the infrastructure in PowerVS.
	var platformStages = []infrastructure.Stage{
		stages.NewStage("powervs",
			"cluster",
			installDir,
			terraformDir,
			[]providers.Provider{providers.IBM, providers.Ignition, providers.Time}),
		stages.NewStage("powervs",
			"bootstrap",
			installDir,
			terraformDir,
			[]providers.Provider{providers.IBM, providers.Ignition, providers.Time},
			stages.WithNormalBootstrapDestroy()),
		stages.NewStage("powervs",
			"bootstrap-routing",
			installDir,
			terraformDir,
			[]providers.Provider{providers.IBM},
			stages.WithCustomBootstrapDestroy(removeFromLoadBalancers)),
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
	opts = append(opts, tfexec.Var("powervs_expose_bootstrap=false"))
	return fmt.Errorf(
		"failed disabling bootstrap load balancing: %w",
		terraform.Apply(directory, powervstypes.Name, s, terraformDir, opts...),
	)
}
