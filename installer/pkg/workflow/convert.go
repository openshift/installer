package workflow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/openshift/installer/pkg/types/config"
)

// ConvertWorkflow creates new instances of the 'convert' workflow,
// responsible for converting an old cluster config.
func ConvertWorkflow(configFilePath string) Workflow {
	return Workflow{
		metadata: metadata{configFilePath: configFilePath},
		steps: []step{
			readTFVarsConfigStep,
			printYAMLConfigStep,
		},
	}
}

func readTFVarsConfigStep(m *metadata) error {
	data, err := ioutil.ReadFile(m.configFilePath)
	if err != nil {
		return err
	}

	m.cluster = config.Cluster{}

	return json.Unmarshal([]byte(data), &m.cluster)
}

func printYAMLConfigStep(m *metadata) error {
	yaml, err := m.cluster.YAML()
	if err != nil {
		return err
	}

	fmt.Println(yaml)

	return nil
}
