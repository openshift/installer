package workflow

import (
	"log"
	"os"
)

// NewDestroyWorkflow creates new instances of the 'destroy' workflow,
// responsible for running the actions required to remove resources
// of an existing cluster and clean up any remaining artefacts.
func NewDestroyWorkflow(buildPath string) Workflow {
	pathStat, err := os.Stat(buildPath)
	// TODO: add deeper checking of the path for having cluster state
	if os.IsNotExist(err) || !pathStat.IsDir() {
		log.Fatalf("Provided path %s is not valid cluster state location.", buildPath)
	} else if err != nil {
		log.Fatalf("%v encountered while validating build location.", err)
	}
	return simpleWorkflow{
		metadata: metadata{
			statePath: buildPath,
		},
		steps: []Step{
			tectonicPrepareStep,
			terraformInitStep,
			terraformDestroyStep,
		},
	}
}

func terraformDestroyStep(m *metadata) error {
	if m.statePath == "" {
		log.Fatalf("Invalid build location - cannot destroy.")
	}

	log.Printf("Destroying cluster from %s...", m.statePath)

	return terraformExec(m, "destroy", "-force")
}
