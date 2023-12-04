package azure

import (
	"encoding/json"
	"fmt"
	"os"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/lbconfig"
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
		stages.WithCustomExtractLBConfig(extractAzureLBConfig),
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

// extractAzureLBConfig extracts the load balancer information from the terraform outputs, generates the
// Load Balancer Config file, regenerates the bootstrap ignition, and updates the terraform variables file.
func extractAzureLBConfig(s stages.SplitStage, directory string, terraformDir string, file *asset.File, tfvarsFile *asset.File) (string, error) {
	// Convert the terraform outputs file into json to extract LB data
	outputs := map[string]interface{}{}
	err := json.Unmarshal(file.Data, &outputs)
	if err != nil {
		return "", err
	}

	userConfiguredDNSRaw, ok := outputs["user_provisioned_dns"]
	if !ok {
		return "", fmt.Errorf("failed to read cluster hosted dns from terraform inputs")
	}
	if !userConfiguredDNSRaw.(bool) {
		return "", nil
	}

	// Extract the Load Balancer ip addresses from the terraform output.
	apiLBIpRaw, ok := outputs["public_lb_pip_v4_ip_address"]
	if !ok {
		return "", fmt.Errorf("failed to read External API LB DNS Name from terraform outputs")
	}
	apiIntLBIpRaw, ok := outputs["internal_lb_ip_v4_address"]
	if !ok {
		return "", fmt.Errorf("failed to read Internal API LB DNS Name from terraform outputs")
	}

	// Parse the terraform input values. Determine if the install is using a user configured dns solution.
	tfvarData := map[string]interface{}{}
	err = json.Unmarshal(tfvarsFile.Data, &tfvarData)
	if err != nil {
		return "", err
	}

	ignitionBootstrap, ok := tfvarData["ignition_bootstrap"]
	if !ok {
		return "", fmt.Errorf("failed to read ignition bootstrap from tfvars")
	}

	ignData := igntypes.Config{}
	err = json.Unmarshal([]byte(ignitionBootstrap.(string)), &ignData)
	if err != nil {
		return "", err
	}

	lbConfig, err := lbconfig.GenerateLBConfigOverride(apiIntLBIpRaw.(string), apiLBIpRaw.(string))
	if err != nil {
		return "", err
	}
	if err := asset.NewDefaultFileWriter(lbConfig).PersistToFile(directory); err != nil {
		return "", fmt.Errorf("failed to save %s to state file: %w", lbConfig.Name(), err)
	}
	path := fmt.Sprintf("/opt/openshift/manifests/%s", lbconfig.ConfigName)
	ignData.Storage.Files = append(ignData.Storage.Files, ignition.FileFromString(path, "root", 0644, string(lbConfig.File.Data)))

	ignitionOutput, err := json.Marshal(ignData)
	if err != nil {
		return "", err
	}

	// Update the ignition bootstrap variable to include the lbconfig.
	tfvarData["ignition_bootstrap"] = string(ignitionOutput)

	// Convert the bootstrap data and write the data back to a file. This will overwrite the original tfvars file.
	jsonBootstrap, err := json.Marshal(tfvarData)
	if err != nil {
		return "", fmt.Errorf("failed to convert bootstrap ignition to bytes: %w", err)
	}
	tfvarsFile.Data = jsonBootstrap

	// update the value on disk to match
	if err := os.WriteFile(fmt.Sprintf("%s/%s", directory, tfvarsFile.Filename), jsonBootstrap, 0o600); err != nil {
		return "", fmt.Errorf("failed to rewrite %s: %w", tfvarsFile.Filename, err)
	}

	return "", nil
}
