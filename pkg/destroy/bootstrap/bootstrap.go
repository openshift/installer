// Package bootstrap uses Terraform to remove bootstrap resources.
package bootstrap

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/pkg/errors"
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

	copyNames := []string{terraform.StateFileName, cluster.TfVarsFileName}

	if platform == "libvirt" {
		err = ioutil.WriteFile(filepath.Join(dir, "disable-bootstrap.tfvars"), []byte(`{
  "bootstrap_dns": false
}
`), 0666)
		if err != nil {
			return err
		}
		copyNames = append(copyNames, "disable-bootstrap.tfvars")
	}

	tempDir, err := ioutil.TempDir("", "openshift-install-")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary directory for Terraform execution")
	}
	defer os.RemoveAll(tempDir)

	for _, filename := range copyNames {
		err = copy(filepath.Join(dir, filename), filepath.Join(tempDir, filename))
		if err != nil {
			return errors.Wrapf(err, "failed to copy %s to the temporary directory", filename)
		}
	}

	if platform == "libvirt" {
		_, err = terraform.Apply(tempDir, platform, fmt.Sprintf("-var-file=%s", filepath.Join(tempDir, "disable-bootstrap.tfvars")))
		if err != nil {
			return errors.Wrap(err, "Terraform apply")
		}
	}

	err = terraform.Destroy(tempDir, platform, "-target=module.bootstrap")
	if err != nil {
		return errors.Wrap(err, "Terraform destroy")
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
