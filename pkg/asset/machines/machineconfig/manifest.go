package machineconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"

	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"

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
	data := []byte{}
	for _, c := range configs {
		if c == nil {
			continue
		}
		configData, err := yaml.Marshal(c)
		if err != nil {
			return nil, err
		}
		data = append(data, []byte("---\n")...)
		data = append(data, configData...)
	}
	if len(data) == 0 {
		return nil, nil
	}
	return []*asset.File{
		{
			Filename: filepath.Join(directory, fmt.Sprintf(machineConfigFileName, role)),
			Data:     data,
		},
	}, nil
}

// IsManifest tests whether the specified filename is a MachineConfig manifest.
func IsManifest(role, filename string) bool {
	return fmt.Sprintf(machineConfigFileName, role) == filename
}

// Load loads the MachineConfig manifests.
func Load(f asset.FileFetcher, role, directory string) ([]*asset.File, error) {
	file, err := f.FetchByName(filepath.Join(directory, fmt.Sprintf(machineConfigFileName, role)))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return []*asset.File{file}, nil
}
