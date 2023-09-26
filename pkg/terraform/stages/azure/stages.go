package azure

import (
	"os"

	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	typesazure "github.com/openshift/installer/pkg/types/azure"
	"github.com/pkg/errors"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	terraformDir, err := terraform.Initialize(installDir)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error initializing terraform")
	}

	// PlatformStages are the stages to run to provision the infrastructure in Azure.
	var platformStages = []infrastructure.Stage{
		stages.NewStage(
			typesazure.Name,
			"vnet",
			installDir,
			terraformDir,
			[]providers.Provider{providers.AzureRM},
		),
		stages.NewStage(
			typesazure.Name,
			"bootstrap",
			installDir,
			terraformDir,
			[]providers.Provider{providers.AzureRM, providers.Ignition, providers.Local},
			stages.WithNormalBootstrapDestroy(),
		),
		stages.NewStage(
			typesazure.Name,
			"cluster",
			installDir,
			terraformDir,
			[]providers.Provider{providers.AzureRM, providers.Time},
		),
	}
	// It would be nice to not need to repeat this for each platform but at this stage
	// Perfect is the enemy of good
	terraform.UnpackTerraform(terraformDir, platformStages)

	cleanup := func() error {
		return os.RemoveAll(terraformDir)
	}

	return platformStages, cleanup, nil
}

func InitializeStackProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	terraformDir, err := terraform.Initialize(installDir)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error initializing terraform")
	}
	// StackPlatformStages are the stages to run to provision the infrastructure in Azure Stack.
	var stackPlatformStages = []infrastructure.Stage{
		stages.NewStage(
			typesazure.StackTerraformName,
			"vnet",
			installDir,
			terraformDir,
			[]providers.Provider{providers.AzureStack},
		),
		stages.NewStage(
			typesazure.StackTerraformName,
			"bootstrap",
			installDir,
			terraformDir,
			[]providers.Provider{providers.AzureStack, providers.Ignition, providers.Local},
			stages.WithNormalBootstrapDestroy(),
		),
		stages.NewStage(
			typesazure.StackTerraformName,
			"cluster",
			installDir,
			terraformDir,
			[]providers.Provider{providers.AzureStack},
		),
	}
	// It would be nice to not need to repeat this for each platform but at this stage
	// Perfect is the enemy of good
	terraform.UnpackTerraform(terraformDir, stackPlatformStages)

	cleanup := func() error {
		return os.RemoveAll(terraformDir)
	}

	return stackPlatformStages, cleanup, nil
}
