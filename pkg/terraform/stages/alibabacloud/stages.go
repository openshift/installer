package alibabacloud

import (
	"os"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	alitypes "github.com/openshift/installer/pkg/types/alibabacloud"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	terraformDir, err := terraform.Initialize(installDir)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error initializing terraform")
	}

	// PlatformStages are the stages to run to provision the infrastructure in Alibaba Cloud.
	var platformStages = []infrastructure.Stage{
		stages.NewStage(
			"alibabacloud",
			"cluster",
			installDir,
			terraformDir,
			[]providers.Provider{providers.AliCloud},
		),
		stages.NewStage(
			"alibabacloud",
			"bootstrap",
			installDir,
			terraformDir,
			[]providers.Provider{providers.AliCloud},
			stages.WithNormalBootstrapDestroy(),
		),
		// This stage has been added due to an issue that occurs when the bootstrap stage performs
		// the teardown of the bootstrap assets. The teardown of the bootstrap SLB attachment causes
		// all of the controlplane backends to be removed from the SLB. This stage attaches the
		// master VMs and the bootstrap VMs as backend servers to the SLBs on create. Later,
		// on bootstrap destroy, this stages removes only the bootstrap VM from the backend servers.
		stages.NewStage(
			"alibabacloud",
			"slb_attach_controlplane",
			installDir,
			terraformDir,
			[]providers.Provider{providers.AliCloud},
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
	opts = append(opts, tfexec.Var("ali_bootstrap_lb=false"))
	return errors.Wrap(
		terraform.Apply(directory, alitypes.Name, s, terraformDir, opts...),
		"failed disabling bootstrap load balancing",
	)
}
