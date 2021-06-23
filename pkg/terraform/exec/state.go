package exec

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/states/statefile"
	"github.com/pkg/errors"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// ReadState reads the terraform state from file and returns the contents in bytes
// It returns an error if reading the state was unsuccessful
// ReadState utilizes the terraform's internal wiring to upconvert versions of terraform state to return
// the state it currently recognizes.
func ReadState(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open %q", file)
	}
	defer f.Close()

	sf, err := statefile.Read(f)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read statefile from %q", file)
	}

	out := bytes.Buffer{}
	if err := statefile.Write(sf, &out); err != nil {
		return nil, errors.Wrapf(err, "failed to write statefile")
	}
	return out.Bytes(), nil
}

// Outputs reads the terraform state file and returns the outputs of the stage as json.
func Outputs(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open %q", file)
	}
	defer f.Close()

	sf, err := statefile.Read(f)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read statefile from %q", file)
	}

	rootModule := sf.State.Module(addrs.RootModuleInstance)

	outputs := make(map[string]interface{}, len(rootModule.OutputValues))
	for key, value := range rootModule.OutputValues {
		outputs[key] = ctyjson.SimpleJSONValue{Value: value.Value}
	}

	data, err := json.Marshal(outputs)
	return data, errors.Wrap(err, "could not marshal outputs")
}
