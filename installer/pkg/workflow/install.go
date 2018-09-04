package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/installer/installer/pkg/config-generator"
)

// InstallWorkflow creates new instances of the 'install' workflow,
// responsible for running the actions necessary to install a new cluster.
func InstallWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []step{
			refreshConfigStep,
			generateClusterConfigMaps,
			readClusterConfigStep,
			generateTLSConfigStep,
			generateClusterConfigMaps,
			installAssetsStep,
			generateIgnConfigStep,
			installTopologyStep,
			installBootstrapStep,
			installMachinesStep,
		},
	}
}

func refreshConfigStep(m *metadata) error {
	if err := readClusterConfigStep(m); err != nil {
		return err
	}
	return generateTerraformVariablesStep(m)
}

func installAssetsStep(m *metadata) error {
	return runInstallStep(m, assetsStep)
}

func installTopologyStep(m *metadata) error {
	return runInstallStep(m, topologyStep)
}

func installBootstrapStep(m *metadata) error {
	if !clusterIsBootstrapped(m.clusterDir) {
		return runInstallStep(m, bootstrapStep)
	}
	return nil
}

func installMachinesStep(m *metadata) error {
	return runInstallStep(m, machinesStep)
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
	return c.GenerateIgnConfig(m.clusterDir)
}

func generateTLSConfigStep(m *metadata) error {
	if err := os.MkdirAll(filepath.Join(m.clusterDir, tlsPath), os.ModeDir|0755); err != nil {
		return fmt.Errorf("failed to create TLS directory at %s", tlsPath)
	}

	c := configgenerator.New(m.cluster)
	return c.GenerateTLSConfig(m.clusterDir)
}
