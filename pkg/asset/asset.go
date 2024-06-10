package asset

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// ClusterCreationError is the error when terraform fails, implying infrastructure failures
	ClusterCreationError = "failed to create cluster"
	// InstallConfigError wraps all configuration errors in one single error
	InstallConfigError = "failed to create install config"
)

// Asset used to install OpenShift.
type Asset interface {
	// Dependencies returns the assets upon which this asset directly depends.
	Dependencies() []Asset

	// Generate generates this asset given the states of its parent assets.
	Generate(context.Context, Parents) error

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

// WritableRuntimeAsset is a WriteableAsset that has files that can be written to disk,
// in addition to a manifest file that contains the runtime object.
type WritableRuntimeAsset interface {
	WritableAsset

	// RuntimeFiles returns the manifest files along with their
	// instantiated runtime object.
	RuntimeFiles() []*RuntimeFile
}

// File is a file for an Asset.
type File struct {
	// Filename is the name of the file.
	Filename string
	// Data is the contents of the file.
	Data []byte
}

// RuntimeFile is a file that contains a manifest file and a runtime object.
type RuntimeFile struct {
	File

	Object client.Object `json:"-"`
}

// PersistToFile writes all of the files of the specified asset into the specified
// directory.
func PersistToFile(asset WritableAsset, directory string) error {
	for _, f := range asset.Files() {
		if f == nil {
			panic("asset.Files() returned nil")
		}
		path := filepath.Join(directory, f.Filename)
		if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
			return errors.Wrap(err, "failed to create dir")
		}
		if err := os.WriteFile(path, f.Data, 0o640); err != nil { //nolint:gosec // no sensitive info
			return errors.Wrap(err, "failed to write file")
		}
	}
	return nil
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

// SortManifestFiles sorts the specified files by file name.
func SortManifestFiles(files []*RuntimeFile) {
	sort.Slice(files, func(i, j int) bool { return files[i].Filename < files[j].Filename })
}
