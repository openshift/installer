package imagebuilder

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/vincent-petithory/dataurl"

	data "github.com/openshift-agent-team/fleeting/data/data/agent"
)

type ConfigBuilder struct {
}

func (c ConfigBuilder) Ignition() ([]byte, error) {
	var err error

	config := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
	}

	config.Storage.Files, err = c.getFiles()
	if err != nil {
		return nil, err
	}

	config.Systemd.Units, err = c.getUnits()
	if err != nil {
		return nil, err
	}

	return json.Marshal(config)
}

func (c ConfigBuilder) getFiles() ([]igntypes.File, error) {
	var readDir func(dirPath string, files []igntypes.File) ([]igntypes.File, error)
	files := make([]igntypes.File, 0)

	readDir = func(dirPath string, files []igntypes.File) ([]igntypes.File, error) {
		entries, err := data.IgnitionData.ReadDir(path.Join("files", dirPath))
		if err != nil {
			return files, fmt.Errorf("Failed to open file dir \"%s\": %w", dirPath, err)
		}
		for _, e := range entries {
			fullPath := path.Join(dirPath, e.Name())
			if e.IsDir() {
				files, err = readDir(fullPath, files)
				if err != nil {
					return files, err
				}
			} else {
				contents, err := data.IgnitionData.ReadFile(path.Join("files", fullPath))
				if err != nil {
					return files, fmt.Errorf("Failed to read file %s: %w", fullPath, err)
				}
				info, err := e.Info()
				if err != nil {
					return files, fmt.Errorf("Failed to get file %s info: %w", fullPath, err)
				}
				mode := int(info.Mode())
				file := igntypes.File{
					Node: igntypes.Node{
						Path:      fullPath,
						Overwrite: ignutil.BoolToPtr(true),
					},
					FileEmbedded1: igntypes.FileEmbedded1{
						Mode: &mode,
						Contents: igntypes.Resource{
							Source: ignutil.StrToPtr(dataurl.EncodeBytes(contents)),
						},
					},
				}
				files = append(files, file)
			}
		}
		return files, nil
	}

	return readDir("/", files)
}

func (c ConfigBuilder) getUnits() ([]igntypes.Unit, error) {
	units := make([]igntypes.Unit, 0)
	basePath := "systemd/units"

	entries, err := data.IgnitionData.ReadDir(basePath)
	if err != nil {
		return units, fmt.Errorf("Failed to read systemd units: %w", err)
	}

	for _, e := range entries {
		contents, err := data.IgnitionData.ReadFile(path.Join(basePath, e.Name()))
		if err != nil {
			return units, fmt.Errorf("Failed to read unit %s: %w", e.Name(), err)
		}

		unit := igntypes.Unit{
			Name:     strings.TrimSuffix(e.Name(), ".template"),
			Enabled:  ignutil.BoolToPtr(true),
			Contents: ignutil.StrToPtr(string(contents)),
		}
		units = append(units, unit)
	}

	return units, nil
}
