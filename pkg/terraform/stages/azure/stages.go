package azure

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
	typesazure "github.com/openshift/installer/pkg/types/azure"
)

// PlatformStages are the stages to run to provision the infrastructure in Azure.
var PlatformStages = []terraform.Stage{
	stages.NewStage(typesazure.Name, "vnet"),
	stages.NewStage(typesazure.Name, "bootstrap", stages.WithNormalBootstrapDestroy()),
	stages.NewStage(typesazure.Name, "cluster"),
}

// StackPlatformStages are the stages to run to provision the infrastructure in Azure Stack.
var StackPlatformStages = []terraform.Stage{
	stages.NewStage(typesazure.StackTerraformName, "vnet"),
	stages.NewStage(typesazure.StackTerraformName, "bootstrap", stages.WithNormalBootstrapDestroy()),
	stages.NewStage(typesazure.StackTerraformName, "cluster"),
}
