// custom marshaling/unmarshaling that will allow to populate nmstate config
// as plain yaml without the need to generate a golang struct following [1]
//
// The yaml parser we use here first convert yaml to json and then manage
// the data with the standard golang json package.
//
// [1] https://github.com/nmstate/nmstate/blob/master/libnmstate/schemas/operational-state.yaml

package v1beta1

import (
	yaml "sigs.k8s.io/yaml"
)

// We are using behind the scenes the golang encode/json so we need to return
// json here for golang to work well, the upper yaml parser will convert it
// to yaml making nmstate yaml transparent to kubernetes-nmstate
func (t NetConfig) MarshalJSON() (output []byte, err error) {
	return yaml.YAMLToJSON([]byte(t.Raw))
}

// Bypass State parsing and directly store it as yaml string to later on
// pass it to namestatectl using it as transparent data at kubernetes-nmstate
func (t *NetConfig) UnmarshalJSON(b []byte) error {
	output, err := yaml.JSONToYAML(b)
	if err != nil {
		return err
	}
	*t = NetConfig{Raw: output}
	return nil
}

// Simple stringer for State
func (t NetConfig) String() string {
	return string(t.Raw)
}
