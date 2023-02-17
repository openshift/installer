package terraform

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/types"
)

var (
	PersistentTerraformFilename = filepath.Join(types.InstallDir, "persistent.tfvars.json")
)

// CreatePersistentTerraformData Creates a terraform file that will persist. The file should
// contain variables that are required to be accessible to stages of the installation.
func CreatePersistentTerraformData(jsonData map[string]interface{}) error {
	if _, err := os.Stat(PersistentTerraformFilename); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	file, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(PersistentTerraformFilename, file, 0644)
}
