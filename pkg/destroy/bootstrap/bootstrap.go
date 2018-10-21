// Package bootstrap uses Terraform to remove bootstrap resources.
package bootstrap

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Destroy uses Terraform to remove bootstrap resources.
func Destroy(dir string) (err error) {
	metadata, err := cluster.LoadMetadata(dir)
	if err != nil {
		return err
	}

	platform := metadata.Platform()
	if platform == "" {
		return errors.New("no platform configured in metadata")
	}

	tempDir, err := ioutil.TempDir("", "openshift-install-")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary directory for Terraform execution")
	}
	defer os.RemoveAll(tempDir)

	for _, filename := range []string{terraform.StateFileName, cluster.TfVarsFileName} {
		err = copy(filepath.Join(dir, filename), filepath.Join(tempDir, filename))
		if err != nil {
			return errors.Wrapf(err, "failed to copy %s to the temporary directory", filename)
		}
	}

	logrus.Infof("Using Terraform to destroy bootstrap resources...")
	err = terraform.Destroy(tempDir, platform, "-target=module.bootstrap")
	if err != nil {
		return errors.Wrap(err, "failed to run terraform")
	}

	tempStateFilePath := filepath.Join(dir, terraform.StateFileName+".new")
	err = copy(filepath.Join(tempDir, terraform.StateFileName), tempStateFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to copy %s from the temporary directory", terraform.StateFileName)
	}
	return os.Rename(tempStateFilePath, filepath.Join(dir, terraform.StateFileName))
}

func copy(from string, to string) error {
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(to, data, 0666)
}
