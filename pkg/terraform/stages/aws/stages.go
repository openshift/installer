package aws

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
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
		stages.WithNormalBootstrapDestroy(),
	),
}
