package agentconfig

import "sigs.k8s.io/yaml"

func unmarshalJSON(b []byte) []byte {
	output, _ := yaml.JSONToYAML(b)
	return output
}
