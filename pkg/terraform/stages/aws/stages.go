package aws

import (
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// PlatformStages are the stages to run to provision the infrastructure in AWS.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		"aws",
		"cluster",
		[]providers.Provider{providers.AWS},
	),
	stages.NewStage(
		"aws",
		"bootstrap",
		[]providers.Provider{providers.AWS},
		stages.WithCustomBootstrapDestroy(customBootstrapDestroy),
	),
}

func customBootstrapDestroy(s stages.SplitStage, directory string, terraformDir string, varFiles []string) error {
	opts := make([]tfexec.DestroyOption, 0, len(varFiles)+1)
	for _, varFile := range varFiles {
		opts = append(opts, tfexec.VarFile(varFile))
	}
	// The bootstrap destroy will no longer refresh state. This was added as a change to counteract
	// the upgrade to the aws terraform provider v5.4.0 where the state changes were causing unsupported
	// operation errors when removing security group rules in sc2s regions.
	logrus.Debugf("aws bootstrap destroy stage will not refresh terraform state")
	opts = append(opts, tfexec.Refresh(false))
	return errors.Wrap(
		terraform.Destroy(directory, awstypes.Name, s, terraformDir, opts...),
		"failed to destroy bootstrap",
	)
}
