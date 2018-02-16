package workflow

// NewDestroyWorkflow creates new instances of the 'destroy' workflow,
// responsible for running the actions required to remove resources
// of an existing cluster and clean up any remaining artefacts.
func NewDestroyWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []Step{
			readClusterConfigStep,
			destroyJoinStep,
			destroyBootstrapStep,
			destroyAssetsStep,
		},
	}
}

func destroyAssetsStep(m *metadata) error {
	return tfDestroy(m.clusterDir, assetsStep, findTemplatesForStep(assetsStep))
}

func destroyBootstrapStep(m *metadata) error {
	return tfDestroy(m.clusterDir, bootstrapStep, findTemplatesForStep(bootstrapStep))
}

func destroyJoinStep(m *metadata) error {
	return tfDestroy(m.clusterDir, joinStep, findTemplatesForStep(joinStep))
}
