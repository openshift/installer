package gcp

import (
	"github.com/pkg/errors"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

// PlatformStages are the stages to run to provision the infrastructure in GCP.
var PlatformStages = []terraform.Stage{
	stages.NewStage("gcp", "cluster"),
	stages.NewStage("gcp", "bootstrap", stages.WithNormalBootstrapDestroy()),
	stages.NewStage("gcp", "post-bootstrap", stages.WithCustomBootstrapDestroy(removeFromLoadBalancers)),
}

// XXX: extraOpts is of type tfexec.DestroyOption, but really should be
// tfexec.ApplyOption. Maybe we should provide custom Init/Apply/Destroy
// function options?
func removeFromLoadBalancers(s stages.SplitStage, directory string, extraOpts []tfexec.DestroyOption) error {
	_, err := terraform.Apply(directory, gcptypes.Name, s, tfexec.Var("gcp_bootstrap_lb=false"))
	return errors.Wrap(err, "failed disabling bootstrap load balancing")
}
