package gcp

import (
	"encoding/json"
	"fmt"
	"os"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/lbconfig"
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
		stages.WithCustomExtractLBConfig(extractGCPLBConfig),
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

func removeFromLoadBalancers(s stages.SplitStage, directory string, terraformDir string, varFiles []string) error {
	opts := make([]tfexec.ApplyOption, 0, len(varFiles)+1)
	for _, varFile := range varFiles {
		opts = append(opts, tfexec.VarFile(varFile))
	}
	opts = append(opts, tfexec.Var("gcp_bootstrap_lb=false"))
	return errors.Wrap(
		terraform.Apply(directory, gcptypes.Name, s, terraformDir, opts...),
		"failed disabling bootstrap load balancing",
	)
}

// extractGCPLBConfig extracts the load balancer information from the terraform outputs, generates the
// Load Balancer Config file, regenerates the bootstrap ignition, and updates the terraform variables file.
func extractGCPLBConfig(s stages.SplitStage, directory string, terraformDir string, file *asset.File, tfvarsFile *asset.File) (string, error) {
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
	apiLBIpRaw, ok := outputs["cluster_public_ip"]
	if !ok {
		return "", fmt.Errorf("failed to read External API LB DNS Name from terraform outputs")
	}
	apiIntLBIpRaw, ok := outputs["cluster_ip"]
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

	err = stages.AddLoadBalancersToInfra(gcptypes.Name, &ignData, []string{apiLBIpRaw.(string)}, []string{apiIntLBIpRaw.(string)})
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
