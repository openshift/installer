package manifests

import (
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

// Read a Yaml file and unmarshal the contents
func GetFileData(fileName string, output interface{}) error {

	path := filepath.Join("./manifests", fileName)

	contents, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading file %s: %w", path, err)
	} else if err = yaml.Unmarshal(contents, output); err != nil {
		err = fmt.Errorf("Error unmarshalling contents of %s: %w", path, err)
	}

	return err
}
