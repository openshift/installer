package asset

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
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
type WritableAsset interface {
	Asset

	// Files returns the files to write.
	Files() []*File
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
func PersistToFile(asset WritableAsset, directory string) error {
	for _, f := range asset.Files() {
		path := filepath.Join(directory, f.Filename)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return errors.Wrap(err, "failed to create dir")
		}
		if err := ioutil.WriteFile(path, f.Data, 0644); err != nil {
			return errors.Wrap(err, "failed to write file")
		}
	}
	return nil
}
