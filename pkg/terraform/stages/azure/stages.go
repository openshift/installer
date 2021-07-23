package azure

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in Azure.
var PlatformStages = []terraform.Stage{
	stages.NewStage("azure", "vnet"),
	stages.NewStage("azure", "bootstrap", stages.WithNormalDestroy()),
	stages.NewStage("azure", "cluster"),
}
