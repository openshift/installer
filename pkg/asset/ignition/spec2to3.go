package ignition

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ConvertSpec2ToSpec3 converts Ignition spec2 to ignition spec3
func ConvertSpec2ToSpec3(spec2data []byte) ([]byte, error) {
	// Unmarshal
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal(spec2data, &jsonMap)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Marshal Ignition config")
	}

	// Replace ignition.version
	ign := jsonMap["ignition"].(map[string]interface{})
	ign["version"] = "3.0.0"

	// ignition.config.append -> ignition.config.merge
	config := ign["config"].(map[string]interface{})
	if val, ok := config["append"]; ok {
		config["merge"] = val
		delete(config, "append")
	}
	ign["config"] = config
	jsonMap["ignition"] = ign

	// Delete networkd section
	if _, ok := jsonMap["networkd"]; ok {
		delete(jsonMap, "networkd")
	}

	// Remove filesystem in storage.files
	if sval, ok := jsonMap["storage"]; ok {
		storage := sval.(map[string]interface{})

		if fval, ok := storage["files"]; ok {
			files := fval.([]interface{})

			updatedFiles := make([]interface{}, 0)

			for i := range files {
				file := files[i].(map[string]interface{})
				if _, ok := file["filesystem"]; ok {
					delete(file, "filesystem")
				}
				updatedFiles = append(updatedFiles, file)
			}
			storage["files"] = updatedFiles
		}
		jsonMap["storage"] = storage
	}

	// Convert to bytes
	spec3data, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Marshal Ignition config")
	}
	return spec3data, nil
}
