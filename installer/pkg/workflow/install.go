package workflow

// NewInstallFullWorkflow creates new instances of the 'install' workflow,
// responsible for running the actions necessary to install a new cluster.
func NewInstallFullWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []Step{
			readClusterConfigStep,
			installAssetsStep,
			generateClusterConfigStep,
			installBootstrapStep,
			installJoinStep,
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
			generateClusterConfigStep,
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
			installBootstrapStep,
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
			installJoinStep,
		},
	}
}

func installAssetsStep(m *metadata) error {
	return runInstallStep(m.clusterDir, assetsStep)
}

func installBootstrapStep(m *metadata) error {
	if err := runInstallStep(m.clusterDir, bootstrapStep); err != nil {
		return err
	}

	if err := waitForNCG(m); err != nil {
		return err
	}

	if err := destroyCNAME(m.clusterDir); err != nil {
		return err
	}

	return nil
}

func installJoinStep(m *metadata) error {
	// TODO: import will fail after a first run, error is ignored for now
	importAutoScalingGroup(m)

	return runInstallStep(m.clusterDir, joinStep)
}

func runInstallStep(clusterDir, step string) error {
	templateDir := findTemplatesForStep(step)
	if err := tfInit(clusterDir, templateDir); err != nil {
		return err
	}

	return tfApply(clusterDir, step, templateDir)
}
