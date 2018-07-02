package workflow

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ImportWorkflow creates new instances of the 'import' workflow,
// responsible for importing a cluster from a config file.
func ImportWorkflow(configFilePath string) Workflow {
	return Workflow{
		metadata: metadata{configFilePath: configFilePath},
		steps: []Step{
			importWorspaceStep,
			refreshConfigStep,
		},
	}
}

func importWorspaceStep(m *metadata) error {
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
	if err := copyFile(m.configFilePath, configFilePath); err != nil {
		return fmt.Errorf("failed to create cluster config at %q: %v", clusterDir, err)
	}

	// generate the internal config file under the clusterDir folder
	return buildInternalConfig(clusterDir)
}
