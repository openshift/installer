package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/coreos/tectonic-installer/installer/pkg/terraform-generator"
)

const (
	kubeSystemPath             = "generated/manifests"
	kubeSystemFileName         = "cluster-config.yaml"
	tectonicSystemPath         = "generated/tectonic"
	tectonicSystemFileName     = "cluster-config.yaml"
	terraformVariablesFileName = "terraform.tfvars"
)

// NewInitWorkflow creates new instances of the 'init' workflow,
// responsible for initializing a new cluster.
func NewInitWorkflow(configFilePath string) Workflow {
	return Workflow{
		metadata: metadata{configFilePath: configFilePath},
		steps: []Step{
			readClusterConfigStep,
			prepareWorspaceStep,
			generateTerraformVariablesStep,
		},
	}
}

func generateTerraformVariablesStep(m *metadata) error {
	terraformGenerator := terraformgenerator.New(m.cluster)

	vars, err := terraformGenerator.TFVars()
	if err != nil {
		return err
	}

	terraformVariablesFilePath := filepath.Join(m.clusterDir, terraformVariablesFileName)
	return writeFile(terraformVariablesFilePath, vars)
}

func prepareWorspaceStep(m *metadata) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Failed to get current directory because: %s", err)
	}

	m.clusterDir = filepath.Join(dir, m.cluster.Name)
	if stat, err := os.Stat(m.clusterDir); err == nil && stat.IsDir() {
		return fmt.Errorf("cluster directory already exists at %s", m.clusterDir)
	}

	if err := os.MkdirAll(m.clusterDir, os.ModeDir|0755); err != nil {
		return fmt.Errorf("Failed to create cluster directory at %s", m.clusterDir)
	}

	configFilePath := filepath.Join(m.clusterDir, configFileName)
	if err := copyFile(m.configFilePath, configFilePath); err != nil {
		return err
	}
	m.configFilePath = configFilePath

	return nil
}
