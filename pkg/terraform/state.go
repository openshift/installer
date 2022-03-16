package terraform

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

// StateFilename is the default name of the terraform state file.
const StateFilename = "terraform.tfstate"

// Outputs reads the terraform state file and returns the outputs of the stage as json.
func Outputs(dir string, terraformDir string) ([]byte, error) {
	tf, err := newTFExec(dir, terraformDir)
	if err != nil {
		return nil, err
	}

	tfoutput, err := tf.Output(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to read terraform state file")
	}

	outputs := make(map[string]interface{}, len(tfoutput))
	for key, value := range tfoutput {
		outputs[key] = value.Value
	}

	data, err := json.Marshal(outputs)
	return data, errors.Wrap(err, "could not marshal outputs")
}
