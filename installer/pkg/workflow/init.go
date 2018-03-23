package workflow

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	configgenerator "github.com/coreos/tectonic-installer/installer/pkg/config-generator"
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
			buildInternalStep,
			generateTerraformVariablesStep,
		},
	}
}

func buildInternalStep(m *metadata) error {
	if m.clusterDir == "" {
		return errors.New("no clusterDir path set in metadata")
	}

	// fill the internal struct
	clusterID, err := configgenerator.GenerateClusterID(16)
	if err != nil {
		return err
	}
	m.cluster.Internal.ClusterID = clusterID

	// store the content
	yamlContent, err := yaml.Marshal(m.cluster.Internal)
	internalFileContent := []byte("# Do not touch, auto-generated\n")
	internalFileContent = append(internalFileContent, yamlContent...)
	if err != nil {
		return err
	}
	return writeFile(filepath.Join(m.clusterDir, internalFileName), string(internalFileContent))
}

func generateTerraformVariablesStep(m *metadata) error {
	vars, err := m.cluster.TFVars()
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
	return copyFile(m.configFilePath, configFilePath)
}
