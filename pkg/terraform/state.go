package terraform

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

// ErrResourceNotFound is an error that instructs that requested resource was not found.
var ErrResourceNotFound = errors.New("resource not found")

// LookupResource finds a resource for a given module, type and name from the state.
// If module is "root", it is treated as ""
// If no resource is found for the triplet, ErrResourceNotFound error is returned.
func LookupResource(state *tfjson.State, module, t, name string) ([]*tfjson.StateResource, error) {
	var err error = ErrResourceNotFound
	var resources []*tfjson.StateResource

	if state.Values != nil {
		rm := state.Values.RootModule
		if rm != nil && module == rm.Address {
			for _, r := range rm.Resources {
				if t == r.Type && name == r.Name {
					resources = append(resources, r)
					err = nil
				}
			}
		}

		for _, cm := range rm.ChildModules {
			if module != cm.Address {
				continue
			}
			for _, r := range cm.Resources {
				if t == r.Type && name == r.Name {
					resources = append(resources, r)
					err = nil
				}
			}
		}
	}

	return resources, err
}

// ReadState returns that terraform state from the file.
func ReadState(file string) (*tfjson.State, error) {
	_, err := os.Stat(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to stat %q", file)
	}

	tfpath := GetTerraformPath()
	datadir := filepath.Dir(file)

	tf, err := tfexec.NewTerraform(datadir, tfpath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find terraform", err)
	}

	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)
	tf.SetLogger(logrus.StandardLogger())

	tfstate, err := tf.ShowStateFile(context.Background(), file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read terraform state file %q", file)
	}

	return tfstate, nil
}

// Outputs reads the terraform state file and returns the outputs of the stage as json.
func Outputs(file string) ([]byte, error) {
	tfstate, err := ReadState(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read terraform state file %q", file)
	}

	if tfstate.Values != nil {
		rm := tfstate.Values.RootModule
		if rm != nil && "" == rm.Address {
			outputs := make(map[string]interface{}, len(tfstate.Values.Outputs))
			for key, value := range tfstate.Values.Outputs {
				outputs[key] = value.Value
			}

			data, err := json.Marshal(outputs)
			if err != nil {
				return nil, errors.Wrap(err, "could not marshal outputs")
			}

			return data, nil
		}
	}
	return nil, nil
}
