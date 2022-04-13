package gcp

import (
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

// PlatformStages are the stages to run to provision the infrastructure in GCP.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		"gcp",
		"cluster",
		[]providers.Provider{providers.Google},
	),
	stages.NewStage(
		"gcp",
		"bootstrap",
		[]providers.Provider{providers.Google, providers.Ignition},
		stages.WithNormalBootstrapDestroy(),
	),
	stages.NewStage(
		"gcp",
		"post-bootstrap",
		[]providers.Provider{providers.Google},
		stages.WithCustomBootstrapDestroy(removeFromLoadBalancers),
	),
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
