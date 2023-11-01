package aws

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

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
		"cluster",
		[]providers.Provider{providers.AWS},
	),
	stages.NewStage(
		"aws",
		"dnsConfig",
		[]providers.Provider{providers.AWS},
		stages.WithCustomAddLBConfig(addAWSLBConfig),
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

func addAWSLBConfig(s stages.SplitStage, directory string, terraformDir string, varFiles []string) (string, error) {
	outputsFilePath := filepath.Join(directory, s.OutputsFilename())
	if _, err := os.Stat(outputsFilePath); err != nil {
		return "", errors.Wrapf(err, "could not find outputs file %q", outputsFilePath)
	}

	outputsFile, err := os.ReadFile(outputsFilePath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read outputs file %q", outputsFilePath)
	}

	outputs := map[string]interface{}{}
	if err := json.Unmarshal(outputsFile, &outputs); err != nil {
		return "", errors.Wrapf(err, "could not unmarshal outputs file %q", outputsFilePath)
	}

	var apiLBDNSNames []string
	if apiLBDNSNamesRaw, ok := outputs["aws_external_api_lb_dns_name"]; ok {
		apiLBDNSNames[0], ok = apiLBDNSNamesRaw.(string)
		if !ok {
			return "", errors.New("could not read External API LB DNS Name from terraform outputs")
		}
	}

	var apiIntLBDNSNames []string
	if apiIntLBDNSNamesRaw, ok := outputs["aws_internal_api_lb_dns_name"]; ok {
		apiIntLBDNSNames[0], ok = apiIntLBDNSNamesRaw.(string)
		if !ok {
			return "", errors.New("could not read API Internal LB DNS Name from terraform outputs")
		}
	}

	// Extract bootstrap ignition from terraform outputs
	var bootstrapIgnition string
	if bootstrapIgnitionRaw, ok := outputs["aws_bootstrap_ignition"]; ok {
		bootstrapIgnition, ok = bootstrapIgnitionRaw.(string)
		if !ok {
			return "", errors.New("could not read bootstrap ignition from terraform outputs")
		}
	}

	// Create openshift-lbconfigfordns configmap and inject into bootstrap ignition
	bootstrapIgnWithLBConfig, err := lbconfig.InjectLBInfo([]byte(bootstrapIgnition), apiLBDNSNames, apiIntLBDNSNames)
	if err != nil {
		return "", errors.Wrap(err, "unable to inject LB ConfigMap into Bootstrap Ignition")
	}

	return bootstrapIgnWithLBConfig, nil
}
