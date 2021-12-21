package gcp

import (
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

func removeFromLoadBalancers(s stages.SplitStage, directory string, extraArgs []string) error {
	_, err := terraform.Apply(directory, gcptypes.Name, s, append(extraArgs, "-var=gcp_bootstrap_lb=false")...)
	return errors.Wrap(err, "failed disabling bootstrap load balancing")
}
