package workflow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/openshift/installer/installer/pkg/config-generator"
	"github.com/openshift/installer/pkg/types/config"
)

// InstallWorkflow creates new instances of the 'install' workflow,
// responsible for running the actions necessary to install a new cluster.
func InstallWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []step{
			readClusterConfigStep,
			generateTerraformVariablesStep,
			generateTLSConfigStep,
			generateClusterConfigMaps,
			generateIgnConfigStep,
			installAssetsStep,
			installInfraStep,
		},
	}
}

func installAssetsStep(m *metadata) error {
	return runInstallStep(m, assetsStep)
}

func installInfraStep(m *metadata) error {
	return runInstallStep(m, infraStep)
}

func runInstallStep(m *metadata, step string, extraArgs ...string) error {
	templateDir, err := findStepTemplates(step, m.cluster.Platform)
	if err != nil {
		return err
	}
	if err := tfInit(m.clusterDir, templateDir); err != nil {
		return err
	}
	return tfApply(m.clusterDir, step, templateDir, extraArgs...)
}

func generateIgnConfigStep(m *metadata) error {
	c := configgenerator.New(m.cluster)
	masterIgns, workerIgn, err := c.GenerateIgnConfig(m.clusterDir)
	if err != nil {
		return fmt.Errorf("failed to generate ignition configs: %v", err)
	}

	terraformVariablesFilePath := filepath.Join(m.clusterDir, terraformVariablesFileName)
	data, err := ioutil.ReadFile(terraformVariablesFilePath)
	if err != nil {
		return fmt.Errorf("failed to read terraform.tfvars: %v", err)
	}

	var cluster config.Cluster
	if err := json.Unmarshal(data, &cluster); err != nil {
		return fmt.Errorf("failed to unmarshal terraform.tfvars: %v", err)
	}

	cluster.IgnitionMasters = masterIgns
	cluster.IgnitionWorker = workerIgn

	data, err = json.MarshalIndent(&cluster, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal terraform.tfvars: %v", err)
	}

	return ioutil.WriteFile(terraformVariablesFilePath, data, 0666)
}

func generateTLSConfigStep(m *metadata) error {
	if err := os.MkdirAll(filepath.Join(m.clusterDir, tlsPath), os.ModeDir|0755); err != nil {
		return fmt.Errorf("failed to create TLS directory at %s", tlsPath)
	}

	c := configgenerator.New(m.cluster)
	return c.GenerateTLSConfig(m.clusterDir)
}
