package aws

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/lbconfig"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// PlatformStages are the stages to run to provision the infrastructure in AWS.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		"aws",
		"vpc",
		[]providers.Provider{providers.AWS},
		stages.WithCustomExtractLBConfig(extractAWSLBConfig),
	),
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

// extractAWSLBConfig extracts the load balancer information from the terraform outputs, generates the
// Load Balancer Config file, regenerates the bootstrap ignition, and updates the terraform variables file.
func extractAWSLBConfig(s stages.SplitStage, directory string, terraformDir string, file *asset.File, tfvarsFile *asset.File) (string, error) {
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
	apiLBIpRaw, ok := outputs["cluster_public_ips"]
	if !ok {
		return "", fmt.Errorf("failed to read External API LB DNS Name from terraform outputs")
	}
	apiIntLBIpRaw, ok := outputs["cluster_internal_ips"]
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

	apiLBIpList := make([]string, len(apiLBIpRaw.([]interface{})))
	for i, ip := range apiLBIpRaw.([]interface{}) {
		apiLBIpList[i] = ip.(string)
	}
	apiLBIps := strings.Join(apiLBIpList[:], ",")

	apiIntLBIpList := make([]string, len(apiIntLBIpRaw.([]interface{})))
	for i, ip := range apiIntLBIpRaw.([]interface{}) {
		apiIntLBIpList[i] = ip.(string)
	}
	apiIntLBIps := strings.Join(apiIntLBIpList[:], ",")

	lbConfig, err := lbconfig.GenerateLBConfigOverride(apiIntLBIps, apiLBIps)
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
