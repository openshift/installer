package workflow

import (
	"log"

	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

// Workflow is a high-level representation
// of a set of actions performed in a predictable order.
type Workflow interface {
	Execute() error
}

// metadata is the state store of the current workflow execution.
// It is meant to carry state for one step to another.
// When creating a new workflow, initial state from external parameters
// is also injected by when initializing the metadata object.
// Steps taked thier inputs from the metadata object and persist
// results onto it for later consumption.
type metadata struct {
	config.Cluster
	configFile string
	statePath  string
}

// Step is the entrypoint of a workflow step implementation.
// To add a new step, put your logic in a function that matches this signature.
// Next, add a refrence to this new function in a Workflow's steps list.
type Step func(*metadata) error

type simpleWorkflow struct {
	metadata metadata
	steps    []Step
}

func (w simpleWorkflow) Execute() error {
	var err error
	for _, step := range w.steps {
		err = step(&w.metadata)
		if err != nil {
			log.Fatal(err) // TODO: actually do proper error handling
		}
	}
	return nil
}
