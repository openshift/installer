package workflow

// NewDestroyWorkflow creates new instances of the 'destroy' workflow,
// responsible for running the actions required to remove resources
// of an existing cluster and clean up any remaining artefacts.
func NewDestroyWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []Step{
			readClusterConfigStep,
			destroyJoinMastersStep,
			destroyJoinWorkersStep,
			destroyEtcdStep,
			destroyBootstrapStep,
			destroyTNCDNSStep,
			destroyTopologyStep,
			destroyAssetsStep,
		},
	}
}

func destroyAssetsStep(m *metadata) error {
	return runDestroyStep(m.clusterDir, assetsStep)
}

func destroyEtcdStep(m *metadata) error {
	return runDestroyStep(m.clusterDir, etcdStep)
}

func destroyBootstrapStep(m *metadata) error {
	return runDestroyStep(m.clusterDir, bootstrapStep)
}

func destroyTNCDNSStep(m *metadata) error {
	return destroyTNCDNS(m.clusterDir)
}

func destroyTopologyStep(m *metadata) error {
	return runDestroyStep(m.clusterDir, topologyStep)
}

func destroyJoinWorkersStep(m *metadata) error {
	return runDestroyStep(m.clusterDir, joinWorkersStep)
}

func destroyJoinMastersStep(m *metadata) error {
	return runDestroyStep(m.clusterDir, joinMastersStep)
}

func runDestroyStep(clusterDir, step string, extraArgs ...string) error {
	if !hasStateFile(clusterDir, step) {
		// there is no statefile, therefore nothing to destroy for this step
		return nil
	}
	templateDir, err := findTemplates(step)
	if err != nil {
		return err
	}

	return tfDestroy(clusterDir, step, templateDir, extraArgs...)
}
