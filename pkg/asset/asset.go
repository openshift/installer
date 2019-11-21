package asset

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
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

// ConfigDriveImage create config drive iso with ignition-config file
// refer to https://coreos.com/os/docs/latest/config-drive.html
func ConfigDriveImage(directory string, asset WritableAsset) error {
	// On the s390x/s390 platform, the igniton-config file will be converted to config iso
	if strings.Contains(asset.Name(), "Ignition Config") {
		name := strings.ToLower(asset.Name())
		filename := strings.Split(name, " ")[0] + ".ign"
		isoname := strings.Split(name, " ")[0] + ".iso"

		logrus.Infof("Starting create config-drive image for %s", asset.Name())

		logrus.Debugf("Creating tmp dir /tmp/new-drive/openstack/latest for %s", filename)
		if err := os.MkdirAll("/tmp/new-drive/openstack/latest", 0755); err != nil {
			return err
		}

		// copy the config file as /openstack/latest/user_data
		if err := os.Chdir(directory); err != nil {
			return err
		}
		src, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer src.Close()
		dst, err := os.Create("/tmp/new-drive/openstack/latest/user_data")
		if err != nil {
			return err
		}
		defer dst.Close()
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		logrus.Debugf("Creating config-drive image for %s", filename)
		// mkisofs will create an ISO 9660 filesystem labeled 'config-2'
		cmd := exec.Command("mkisofs", "-R", "-V", "config-2", "-o", "/tmp/"+isoname, "/tmp/new-drive/")
		_, err = cmd.Output()
		if err != nil {
			return err
		}
		logrus.Debugf("remove tmp dir /tmp/new-drive/openstack")
		if err := os.RemoveAll("/tmp/new-drive/openstack"); err != nil {
			return err
		}
		logrus.Debugf("Config drive image for %s created", filename)
		return nil
	}
	return nil
}
