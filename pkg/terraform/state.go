package terraform

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
)

// Outputs reads the terraform state file and returns the outputs of the stage as json.
func Outputs(dir string, file string) ([]byte, error) {
	tf, err := newTFExec(dir)
	if err != nil {
		return nil, err
	}

	tfoutput, err := tf.Output(context.Background(), tfexec.State(file))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read terraform state file %q", file)
	}

	outputs := make(map[string]interface{}, len(tfoutput))
	for key, value := range tfoutput {
		outputs[key] = value.Value
	}

	data, err := json.Marshal(outputs)
	return data, errors.Wrap(err, "could not marshal outputs")
}
