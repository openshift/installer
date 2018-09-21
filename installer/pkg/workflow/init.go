package workflow

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	configgenerator "github.com/openshift/installer/installer/pkg/config-generator"
	"github.com/openshift/installer/pkg/types/config"
)

const (
	generatedPath              = "generated"
	kcoConfigFileName          = "kco-config.yaml"
	maoConfigFileName          = "mao-config.yaml"
	kubeSystemPath             = "generated/manifests"
	kubeSystemFileName         = "cluster-config.yaml"
	tectonicSystemPath         = "generated/tectonic"
	tlsPath                    = "generated/tls"
	tectonicSystemFileName     = "cluster-config.yaml"
	terraformVariablesFileName = "terraform.tfvars"
)

// InitWorkflow creates new instances of the 'init' workflow,
// responsible for initializing a new cluster.
func InitWorkflow(configFilePath string) Workflow {
	return Workflow{
		metadata: metadata{configFilePath: configFilePath},
		steps: []step{
			prepareWorspaceStep,
		},
	}
}

func buildInternalConfig(clusterDir string) error {
	if clusterDir == "" {
		return errors.New("no cluster dir given for building internal config")
	}

	// fill the internal struct
	clusterID, err := configgenerator.GenerateClusterID(16)
	if err != nil {
		return err
	}
	internalCfg := config.Internal{
		ClusterID: clusterID,
	}

	// store the content
	yamlContent, err := yaml.Marshal(internalCfg)
	internalFileContent := []byte("# Do not touch, auto-generated\n")
	internalFileContent = append(internalFileContent, yamlContent...)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(clusterDir, internalFileName), internalFileContent, 0666)
}

func generateTerraformVariablesStep(m *metadata) error {
	vars, err := m.cluster.TFVars()
	if err != nil {
		return err
	}

	terraformVariablesFilePath := filepath.Join(m.clusterDir, terraformVariablesFileName)
	return ioutil.WriteFile(terraformVariablesFilePath, []byte(vars), 0666)
}

func prepareWorspaceStep(m *metadata) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	if m.configFilePath == "" {
		return errors.New("a path to a config file is required")
	}

	// load initial cluster config to get cluster.Name
	cluster, err := readClusterConfig(m.configFilePath, "")
	if err != nil {
		return fmt.Errorf("failed to get configuration from file %q: %v", m.configFilePath, err)
	}

	if err := cluster.ValidateAndLog(); err != nil {
		return err
	}

	if cluster.Platform == config.PlatformLibvirt {
		if err := cluster.Libvirt.UseCachedImage(); err != nil {
			return err
		}
	}

	// generate clusterDir folder
	clusterDir := filepath.Join(dir, cluster.Name)
	m.clusterDir = clusterDir
	if stat, err := os.Stat(clusterDir); err == nil && stat.IsDir() {
		return fmt.Errorf("cluster directory already exists at %q", clusterDir)
	}

	if err := os.MkdirAll(clusterDir, os.ModeDir|0755); err != nil {
		return fmt.Errorf("failed to create cluster directory at %q", clusterDir)
	}

	// put config file under the clusterDir folder
	configFilePath := filepath.Join(clusterDir, configFileName)
	configContent, err := yaml.Marshal(cluster)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(configFilePath, configContent, 0666); err != nil {
		return fmt.Errorf("failed to create cluster config at %q: %v", clusterDir, err)
	}

	// generate the internal config file under the clusterDir folder
	return buildInternalConfig(clusterDir)
}
