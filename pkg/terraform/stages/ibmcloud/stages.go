package ibmcloud

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	ibmcloudtfvars "github.com/openshift/installer/pkg/tfvars/ibmcloud"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
)

// PlatformStages are the stages to run to provision the infrastructure in IBM Cloud.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		"ibmcloud",
		"network",
		[]providers.Provider{providers.IBM},
	),
	stages.NewStage(
		"ibmcloud",
		"bootstrap",
		[]providers.Provider{providers.IBM},
		stages.WithCustomBootstrapDestroy(customBootstrapDestroy),
	),
	stages.NewStage(
		"ibmcloud",
		"master",
		[]providers.Provider{providers.IBM},
	),
}

func customBootstrapDestroy(s stages.SplitStage, directory string, terraformDir string, varFiles []string) error {
	opts := make([]tfexec.DestroyOption, 0, len(varFiles)+1)
	for _, varFile := range varFiles {
		opts = append(opts, tfexec.VarFile(varFile))
	}

	// If these is a endpoint override JSON file in the terraformDir's parent directory (terraformDir isn't available during JSON file creation),
	// we want to inject that file into the Terraform variables so IBM Cloud Service endpoints are overridden.
	terraformParentDir := filepath.Dir(terraformDir)
	endpointOverrideFile := filepath.Join(terraformParentDir, ibmcloudtfvars.IBMCloudEndpointJSONFileName)
	if _, err := os.Stat(endpointOverrideFile); err == nil {
		// Set variable to use private endpoints (overrides) from JSON file, via the IBM Cloud Terraform variable: 'ibmcloud_endpoints_json_file'.
		opts = append(opts, tfexec.Var(fmt.Sprintf("ibmcloud_endpoints_json_file=%s", endpointOverrideFile)))
		logrus.Debugf("configuring terraform bootstrap destroy with ibm endpoint overrides: %s", endpointOverrideFile)
	}
	err := terraform.Destroy(directory, ibmcloudtypes.Name, s, terraformDir, opts...)
	if err != nil {
		return fmt.Errorf("failed to destroy bootstrap: %w", err)
	}

	return nil
}
