package azure

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	typesazure "github.com/openshift/installer/pkg/types/azure"
)

// PlatformStages are the stages to run to provision the infrastructure in Azure.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		typesazure.Name,
		"vnet",
		[]providers.Provider{providers.AzureRM},
	),
	stages.NewStage(
		typesazure.Name,
		"bootstrap",
		[]providers.Provider{providers.AzureRM, providers.Ignition, providers.Local},
		stages.WithNormalBootstrapDestroy(),
	),
	stages.NewStage(
		typesazure.Name,
		"cluster",
		[]providers.Provider{providers.AzureRM, providers.Time},
	),
}

// StackPlatformStages are the stages to run to provision the infrastructure in Azure Stack.
var StackPlatformStages = []terraform.Stage{
	stages.NewStage(
		typesazure.StackTerraformName,
		"vnet",
		[]providers.Provider{providers.AzureStack},
	),
	stages.NewStage(
		typesazure.StackTerraformName,
		"bootstrap",
		[]providers.Provider{providers.AzureStack, providers.Ignition, providers.Local},
		stages.WithNormalBootstrapDestroy(),
	),
	stages.NewStage(
		typesazure.StackTerraformName,
		"cluster",
		[]providers.Provider{providers.AzureStack},
	),
}
