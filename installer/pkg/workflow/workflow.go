package workflow

import "github.com/coreos/tectonic-installer/installer/pkg/config"

// metadata is the state store of the current workflow execution.
// It is meant to carry state for one step to another.
// When creating a new workflow, initial state from external parameters
// is also injected by when initializing the metadata object.
// Steps taked their inputs from the metadata object and persist
// results onto it for later consumption.
type metadata struct {
	cluster        config.Cluster
	configFilePath string
	clusterDir     string
}

// Step is the entrypoint of a workflow step implementation.
// To add a new step, put your logic in a function that matches this signature.
// Next, add a reference to this new function in a Workflow's steps list.
type Step func(*metadata) error

// Workflow is a high-level representation
// of a set of actions performed in a predictable order.
type Workflow struct {
	metadata metadata
	steps    []Step
}

// Execute runs all steps in order.
func (w Workflow) Execute() error {
	for _, step := range w.steps {
		if err := step(&w.metadata); err != nil {
			return err
		}
	}

	return nil
}
