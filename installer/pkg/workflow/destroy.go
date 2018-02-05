package workflow

import (
	"log"
	"os"
	"os/exec"

	"github.com/coreos/tectonic-installer/installer/pkg/tectonic"
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
			terraformPrepareStep,
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
	tfDestroy := exec.Command("terraform", "destroy", "-force", tectonic.FindTemplatesForType("aws")) // TODO: get from cluster config
	tfDestroy.Dir = m.statePath
	tfDestroy.Stdin = os.Stdin
	tfDestroy.Stdout = os.Stdout
	tfDestroy.Stderr = os.Stderr
	err := tfDestroy.Run()
	if err != nil {
		return err
	}
	return nil
}
