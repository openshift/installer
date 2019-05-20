package terraform

import (
	"encoding/json"

	"github.com/pkg/errors"

	tfexec "github.com/openshift/installer/pkg/terraform/exec"
)

// State in local sparse representation of terraform state that includes
// the fields important to installer.
type State struct {
	Resources []StateResource `json:"resources"`
}

// StateResource is local sparse representation of terraform state resource that includes
// the fields most important to installer.
type StateResource struct {
	Module    string                  `json:"module"`
	Name      string                  `json:"name"`
	Type      string                  `json:"type"`
	Instances []StateResourceInstance `json:"instances"`
}

// StateResourceInstance is an instance of terraform state resource.
type StateResourceInstance struct {
	Attributes map[string]interface{} `json:"attributes"`
}

// ErrResourceNotFound is an error that instructs that requested resource was not found.
var ErrResourceNotFound = errors.New("resource not found")

// LookupResource finds a resource for a given module, type and name from the state.
// If module is "root", it is treated as ""
// If no resource is found for the triplet, ErrResourceNotFound error is returned.
func LookupResource(state *State, module, t, name string) (*StateResource, error) {
	if module == "root" {
		module = ""
	}
	for idx, r := range state.Resources {
		if module == r.Module && t == r.Type && name == r.Name {
			return &state.Resources[idx], nil
		}
	}
	return nil, ErrResourceNotFound
}

// ReadState returns that terraform state from the file.
func ReadState(file string) (*State, error) {
	sfRaw, err := tfexec.ReadState(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read %q", file)
	}

	var tfstate State
	if err := json.Unmarshal(sfRaw, &tfstate); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal %q", file)
	}
	return &tfstate, nil
}
