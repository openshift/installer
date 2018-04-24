package workflow

import "github.com/coreos/tectonic-installer/installer/pkg/config-generator"

// NewInstallFullWorkflow creates new instances of the 'install' workflow,
// responsible for running the actions necessary to install a new cluster.
func NewInstallFullWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []Step{
			refreshConfigStep,
			generateClusterConfigMaps,
			installTLSAssetsStep,
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

// NewInstallTLSWorkflow creates the TLS assets, previously created by the
// "assets" step
func NewInstallTLSWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []Step{
			readClusterConfigStep,
			generateClusterConfigMaps,
			installTLSAssetsStep,
		},
	}
}

// NewInstallAssetsWorkflow creates new instances of the 'assets' workflow,
// responsible for running the actions necessary to generate cluster assets.
func NewInstallAssetsWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []Step{
			refreshConfigStep,
			generateClusterConfigMaps,
			installAssetsStep,
			generateIgnConfigStep,
		},
	}
}

// NewInstallBootstrapWorkflow creates new instances of the 'bootstrap' workflow,
// responsible for running the actions necessary to generate a single bootstrap machine cluster.
func NewInstallBootstrapWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []Step{
			refreshConfigStep,
			installTopologyStep,
			installTNCCNAMEStep,
			installBootstrapStep,
			installTNCARecordStep,
			installEtcdStep,
		},
	}
}

// NewInstallJoinWorkflow creates new instances of the 'join' workflow,
// responsible for running the actions necessary to scale the machines of the cluster.
func NewInstallJoinWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []Step{
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

func installTLSAssetsStep(m *metadata) error {
	return runInstallStep(m, tlsStep)

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
