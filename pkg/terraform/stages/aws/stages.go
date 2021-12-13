package aws

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in AWS.
var PlatformStages = []terraform.Stage{
	stages.NewStage("aws", "cluster"),
	stages.NewStage("aws", "bootstrap", stages.WithNormalBootstrapDestroy()),
}
