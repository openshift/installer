package aws

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	"github.com/openshift/installer/pkg/types"
)

// PlatformStages are the stages to run to provision the infrastructure in AWS.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		"aws",
		"cluster",
		[]providers.Provider{providers.AWS},
		stages.WithCustomBootstrapDestroy(extractAWSLBAddresses),
	),
	stages.NewStage(
		"aws",
		"bootstrap",
		[]providers.Provider{providers.AWS},
		stages.WithNormalBootstrapDestroy(),
	),
}

func extractAWSLBAddresses(s stages.SplitStage, directory string, terraformDir string, varFiles []string) error {
	tfstate, err := stages.GetTerraformOutputs(s, types.InstallDir)
	if err != nil {
		return err
	}

	loadBalancerData := make(map[string]interface{})
	if apiLBName, ok := tfstate["aws_lb_api_external_dns_name"]; ok {
		loadBalancerData["api_external_dns_names"] = []string{apiLBName.(string)}
	}
	if apiIntLBName, ok := tfstate["aws_lb_api_internal_dns_name"]; ok {
		loadBalancerData["api_internal_dns_names"] = []string{apiIntLBName.(string)}
	}

	return terraform.CreatePersistentTerraformData(loadBalancerData)
}
