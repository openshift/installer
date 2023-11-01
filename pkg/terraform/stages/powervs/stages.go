package powervs

import (
	"fmt"

	"github.com/hashicorp/terraform-exec/tfexec"

	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
)

// PlatformStages are the stages to run to provision the infrastructure in PowerVS.
var PlatformStages = []terraform.Stage{
	stages.NewStage("powervs",
		"cluster",
		[]providers.Provider{providers.IBM, providers.Ignition, providers.Time}),
	stages.NewStage("powervs",
		"bootstrap",
		[]providers.Provider{providers.IBM, providers.Ignition, providers.Time},
		stages.WithNormalBootstrapDestroy()),
	stages.NewStage("powervs",
		"bootstrap-routing",
		[]providers.Provider{providers.IBM},
		stages.WithCustomBootstrapDestroy(removeFromLoadBalancers)),
}

func removeFromLoadBalancers(s stages.SplitStage, directory string, terraformDir string, varFiles []string) error {
	opts := make([]tfexec.ApplyOption, 0, len(varFiles)+1)
	for _, varFile := range varFiles {
		opts = append(opts, tfexec.VarFile(varFile))
	}
	opts = append(opts, tfexec.Var("powervs_expose_bootstrap=false"))
	err := terraform.Apply(directory, powervstypes.Name, s, terraformDir, opts...)
	if err == nil {
		return nil
	}
	return fmt.Errorf(
		"failed disabling bootstrap load balancing: %w",
		err,
	)
}
