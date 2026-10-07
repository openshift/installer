package machineconfig

import (
	"fmt"
	"path/filepath"

	"sigs.k8s.io/yaml"

	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/asset"
)

const (
	machineConfigFileName = "99_openshift-machineconfig_%s.yaml"
)

var (
	machineConfigFileNamePattern = fmt.Sprintf(machineConfigFileName, "*")
)

// Manifests creates manifest files containing the MachineConfigs.
func Manifests(configs []*mcfgv1.MachineConfig, role, directory string) ([]*asset.File, error) {
	var ret []*asset.File
	for _, c := range configs {
		if c == nil {
			continue
		}
		configData, err := yaml.Marshal(c)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &asset.File{
			// The filename is derived from the MachineConfig name alone, so two
			// configs sharing a name produce one file and the second silently wins.
			// Callers are responsible for keeping those names unique; disk setup
			// does so by rejecting colliding labels during validation.
			Filename: filepath.Join(directory, fmt.Sprintf(machineConfigFileName, c.ObjectMeta.Name)),
			Data:     configData,
		})
	}
	if len(ret) == 0 {
		return nil, nil
	}
	return ret, nil
}

// IsManifest tests whether the specified filename is a MachineConfig manifest.
func IsManifest(filename string) (bool, error) {
	matched, err := filepath.Match(machineConfigFileNamePattern, filename)
	if err != nil {
		return false, err
	}
	return matched, nil
}

// Load loads the MachineConfig manifests.
func Load(f asset.FileFetcher, role, directory string) ([]*asset.File, error) {
	return f.FetchByPattern(filepath.Join(directory, machineConfigFileNamePattern))
}
