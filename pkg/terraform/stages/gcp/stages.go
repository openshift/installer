package gcp

import (
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

// PlatformStages are the stages to run to provision the infrastructure in GCP.
var PlatformStages = []terraform.Stage{
	stages.NewStage("gcp", "cluster"),
	stages.NewStage("gcp", "bootstrap", stages.WithNormalDestroy()),
	stages.NewStage("gcp", "post-bootstrap", stages.WithCustomDestroy(removeFromLoadBalancers)),
}

func removeFromLoadBalancers(s stages.SplitStage, directory string, extraArgs []string) error {
	_, err := terraform.Apply(directory, gcptypes.Name, s, append(extraArgs, "-var=gcp_bootstrap_lb=false")...)
	return errors.Wrap(err, "failed disabling bootstrap load balancing")
}
