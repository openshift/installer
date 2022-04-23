package alibabacloud

import (
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	alitypes "github.com/openshift/installer/pkg/types/alibabacloud"
)

// PlatformStages are the stages to run to provision the infrastructure in Alibaba Cloud.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		"alibabacloud",
		"cluster",
		[]providers.Provider{providers.AliCloud},
	),
	stages.NewStage(
		"alibabacloud",
		"bootstrap",
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
		[]providers.Provider{providers.AliCloud},
		stages.WithCustomBootstrapDestroy(removeFromLoadBalancers),
	),
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
