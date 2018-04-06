package workflow

import "github.com/coreos/tectonic-installer/installer/pkg/config-generator"

// NewInstallFullWorkflow creates new instances of the 'install' workflow,
// responsible for running the actions necessary to install a new cluster.
func NewInstallFullWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []Step{
			readClusterConfigStep,
			installAssetsStep,
			generateKubeConfigStep,
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

// NewInstallAssetsWorkflow creates new instances of the 'assets' workflow,
// responsible for running the actions necessary to generate cluster assets.
func NewInstallAssetsWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []Step{
			readClusterConfigStep,
			installAssetsStep,
			generateKubeConfigStep,
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
			readClusterConfigStep,
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
			readClusterConfigStep,
			installJoinMastersStep,
			installJoinWorkersStep,
		},
	}
}

func installAssetsStep(m *metadata) error {
	return runInstallStep(m.clusterDir, assetsStep)
}

func installTopologyStep(m *metadata) error {
	return runInstallStep(m.clusterDir, topologyStep)
}

func installBootstrapStep(m *metadata) error {
	return runInstallStep(m.clusterDir, bootstrapStep)
}

func installTNCCNAMEStep(m *metadata) error {
	if !clusterIsBootstrapped(m.clusterDir) {
		return createTNCCNAME(m.clusterDir)
	}
	return nil
}

func installTNCARecordStep(m *metadata) error {
	return createTNCARecord(m.clusterDir)
}

func installEtcdStep(m *metadata) error {
	return runInstallStep(m.clusterDir, etcdStep)
}

func installJoinMastersStep(m *metadata) error {
	// TODO: import will fail after a first run, error is ignored for now
	importAutoScalingGroup(m)
	return runInstallStep(m.clusterDir, joinMastersStep)
}

func installJoinWorkersStep(m *metadata) error {
	return runInstallStep(m.clusterDir, joinWorkersStep)
}

func runInstallStep(clusterDir, step string, extraArgs ...string) error {
	templateDir, err := findTemplates(step)
	if err != nil {
		return err
	}
	if err := tfInit(clusterDir, templateDir); err != nil {
		return err
	}
	return tfApply(clusterDir, step, templateDir, extraArgs...)
}

func generateIgnConfigStep(m *metadata) error {
	c := configgenerator.New(m.cluster)
	return c.GenerateIgnConfig(m.clusterDir)
}
