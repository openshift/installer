package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/installer/installer/pkg/config-generator"
)

// InstallFullWorkflow creates new instances of the 'install' workflow,
// responsible for running the actions necessary to install a new cluster.
func InstallFullWorkflow(clusterDir string) Workflow {
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
			installTNCCNAMEStep,
			installBootstrapStep,
			installTNCARecordStep,
			installEtcdStep,
			installJoinMastersStep,
			installJoinWorkersStep,
		},
	}
}

// InstallTLSNewWorkflow generates the TLS certificates using go, instead of TF
func InstallTLSNewWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []step{
			refreshConfigStep,
			generateClusterConfigMaps,
			generateTLSConfigStep,
		},
	}
}

// InstallAssetsWorkflow creates new instances of the 'assets' workflow,
// responsible for running the actions necessary to generate cluster assets.
func InstallAssetsWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []step{
			refreshConfigStep,
			generateClusterConfigMaps,
			installAssetsStep,
			generateIgnConfigStep,
		},
	}
}

// InstallBootstrapWorkflow creates new instances of the 'bootstrap' workflow,
// responsible for running the actions necessary to generate a single bootstrap machine cluster.
func InstallBootstrapWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []step{
			refreshConfigStep,
			installTopologyStep,
			installTNCCNAMEStep,
			installBootstrapStep,
			installTNCARecordStep,
			installEtcdStep,
		},
	}
}

// InstallJoinWorkflow creates new instances of the 'join' workflow,
// responsible for running the actions necessary to scale the machines of the cluster.
func InstallJoinWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []step{
			refreshConfigStep,
			installJoinMastersStep,
			installJoinWorkersStep,
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
		return runInstallStep(m, mastersStep, []string{bootstrapOn}...)
	}
	return nil
}

func installTNCCNAMEStep(m *metadata) error {
	if !clusterIsBootstrapped(m.clusterDir) {
		return createTNCCNAME(m)
	}
	return nil
}

func installTNCARecordStep(m *metadata) error {
	return createTNCARecord(m)
}

func installEtcdStep(m *metadata) error {
	return runInstallStep(m, etcdStep)
}

func installJoinMastersStep(m *metadata) error {
	return runInstallStep(m, mastersStep, []string{bootstrapOff}...)
}

func installJoinWorkersStep(m *metadata) error {
	return runInstallStep(m, joinWorkersStep)
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
