// Package bootstrap uses Terraform to remove bootstrap resources.
package bootstrap

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/pkg/errors"
)

// Disable uses Terraform to disable bootstrap resources, i.e. unregister from load balancer
func Disable(dir string) (err error) {
	return update(dir, false)
}

// Destroy uses Terraform to remove bootstrap resources.
func Destroy(dir string) (err error) {
	return update(dir, true)
}

func update(dir string, destroy bool) (err error) {
	metadata, err := cluster.LoadMetadata(dir)
	if err != nil {
		return err
	}

	platform := metadata.Platform()
	if platform == "" {
		return errors.New("no platform configured in metadata")
	}

	tfPlatformVarsFileName := fmt.Sprintf(cluster.TfPlatformVarsFileName, platform)
	copyNames := []string{terraform.StateFileName, cluster.TfVarsFileName, tfPlatformVarsFileName}

	err = ioutil.WriteFile(filepath.Join(dir, "disable-bootstrap.tfvars"), []byte(`{
  "bootstrap_dns": false
}
`), 0666)
	if err != nil {
		return err
	}
	copyNames = append(copyNames, "disable-bootstrap.tfvars")

	tempDir, err := ioutil.TempDir("", "openshift-install-")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary directory for Terraform execution")
	}
	defer os.RemoveAll(tempDir)

	extraArgs := []string{}
	for _, filename := range copyNames {
		sourcePath := filepath.Join(dir, filename)
		targetPath := filepath.Join(tempDir, filename)
		err = copy(sourcePath, targetPath)
		if err != nil {
			if os.IsNotExist(err) && err.(*os.PathError).Path == sourcePath && filename == tfPlatformVarsFileName {
				continue // platform may not need platform-specific Terraform variables
			}
			return errors.Wrapf(err, "failed to copy %s to the temporary directory", filename)
		}
		if strings.HasSuffix(filename, ".tfvars") {
			extraArgs = append(extraArgs, fmt.Sprintf("-var-file=%s", targetPath))
		}
	}

	if !destroy || platform == libvirt.Name {
		_, err = terraform.Apply(tempDir, platform, extraArgs...)
		if err != nil {
			return errors.Wrap(err, "Terraform apply")
		}
		if !destroy {
			return
		}
	}

	extraArgs = append(extraArgs, "-target=module.bootstrap")
	err = terraform.Destroy(tempDir, platform, extraArgs...)
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
