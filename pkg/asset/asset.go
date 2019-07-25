package asset

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Asset used to install OpenShift.
type Asset interface {
	// Dependencies returns the assets upon which this asset directly depends.
	Dependencies() []Asset

	// Generate generates this asset given the states of its parent assets.
	Generate(Parents) error

	// Name returns the human-friendly name of the asset.
	Name() string
}

// WritableAsset is an Asset that has files that can be written to disk.
// It can also be loaded from disk.
type WritableAsset interface {
	Asset

	// Files returns the files to write.
	Files() []*File

	// Load returns the on-disk asset if it exists.
	// The asset object should be changed only when it's loaded successfully.
	Load(FileFetcher) (found bool, err error)
}

// File is a file for an Asset.
type File struct {
	// Filename is the name of the file.
	Filename string
	// Data is the contents of the file.
	Data []byte
}

// PersistToFile writes all of the files of the specified asset into the specified
// directory.
func PersistToFile(asset WritableAsset, directory string, fcos bool) error {
	for _, f := range asset.Files() {
		path := filepath.Join(directory, f.Filename)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return errors.Wrap(err, "failed to create dir")
		}
		if fcos && strings.HasSuffix(f.Filename, ".ign") {
			// Run transpiler here
			spec3data, err := convertSpec2ToSpec3(f.Data)
			if err != nil {
				return errors.Wrap(err, "failed to convert spec2 to spec3")
			}
			f.Data = spec3data
		}
		if err := ioutil.WriteFile(path, f.Data, 0644); err != nil {
			return errors.Wrap(err, "failed to write file")
		}
	}
	return nil
}

func convertSpec2ToSpec3(spec2data []byte) ([]byte, error) {
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

// DeleteAssetFromDisk removes all the files for asset from disk.
// this is function is not safe for calling concurrently on the same directory.
func DeleteAssetFromDisk(asset WritableAsset, directory string) error {
	logrus.Debugf("Purging asset %q from disk", asset.Name())
	for _, f := range asset.Files() {
		path := filepath.Join(directory, f.Filename)
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return errors.Wrap(err, "failed to remove file")
		}

		dir := filepath.Dir(path)
		ok, err := isDirEmpty(dir)
		if err != nil && !os.IsNotExist(err) {
			return errors.Wrap(err, "failed to read directory")
		}
		if ok {
			if err := os.Remove(dir); err != nil {
				return errors.Wrap(err, "failed to remove directory")
			}
		}
	}
	return nil
}

func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// SortFiles sorts the specified files by file name.
func SortFiles(files []*File) {
	sort.Slice(files, func(i, j int) bool { return files[i].Filename < files[j].Filename })
}
