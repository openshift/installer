package aws

import (
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// PlatformStages are the stages to run to provision the infrastructure in AWS.
var PlatformStages = []terraform.Stage{
	stages.NewStage("aws", "cluster", stages.WithCustomDestroy(destroyBootstrapInstanceProfile)),
	stages.NewStage("aws", "bootstrap", stages.WithNormalDestroy()),
}

func destroyBootstrapInstanceProfile(s stages.SplitStage, directory string, extraArgs []string) error {
	_, err := terraform.Apply(directory, awstypes.Name, s, append(extraArgs, "-var=destroy_bootstrap=true")...)
	return errors.Wrap(err, "failed destroying bootstrap instance profile")
}
