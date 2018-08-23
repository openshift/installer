package workflow

// DestroyWorkflow creates new instances of the 'destroy' workflow,
// responsible for running the actions required to remove resources
// of an existing cluster and clean up any remaining artefacts.
func DestroyWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []step{
			refreshConfigStep,
			destroyJoinMastersStep,
			destroyJoinWorkersStep,
			destroyBootstrapStep,
			destroyTopologyStep,
			destroyAssetsStep,
		},
	}
}

func destroyAssetsStep(m *metadata) error {
	return runDestroyStep(m, assetsStep)
}

func destroyBootstrapStep(m *metadata) error {
	return runDestroyStep(m, bootstrapStep)
}

func destroyTopologyStep(m *metadata) error {
	return runDestroyStep(m, topologyStep)
}

func destroyJoinWorkersStep(m *metadata) error {
	return runDestroyStep(m, joinWorkersStep)
}

func destroyJoinMastersStep(m *metadata) error {
	return runDestroyStep(m, mastersStep)
}

func runDestroyStep(m *metadata, step string, extraArgs ...string) error {
	if !hasStateFile(m.clusterDir, step) {
		// there is no statefile, therefore nothing to destroy for this step
		return nil
	}
	templateDir, err := findStepTemplates(step, m.cluster.Platform)
	if err != nil {
		return err
	}

	return tfDestroy(m.clusterDir, step, templateDir, extraArgs...)
}
