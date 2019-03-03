// Package bootstrap uses Terraform to remove bootstrap resources.
package bootstrap

import (
	"fmt"
	"github.com/openshift/installer/pkg/asset/byo"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/types/libvirt"
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

	tfPlatformVarsFileName := fmt.Sprintf(cluster.TfPlatformVarsFileName, platform)
	copyNames := []string{terraform.StateFileName, cluster.TfVarsFileName, tfPlatformVarsFileName}

	if platform == libvirt.Name {
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

	if metadata.BYO {
		copyNames = append(copyNames, "main.tf", fmt.Sprintf("variables-%s.tf", metadata.Platform()))

		var byoFiles []string
		modules, err := filepath.Glob(filepath.Join(dir, "**/*.tf"))
		if err != nil {
			return err
		}
		byoFiles = append(byoFiles, modules...)

		plugins, err := filepath.Glob(filepath.Join(dir, byo.PluginsDir, "**/*"))
		if err != nil {
			return err
		}
		byoFiles = append(byoFiles, plugins...)

		for _, filename := range byoFiles {
			copyNames = append(copyNames, strings.Replace(filename, filepath.ToSlash(dir), "", 1))
		}
	}

	extraArgs := []string{}
	for _, filename := range copyNames {
		sourcePath := filepath.Join(dir, filename)
		targetPath := filepath.Join(tempDir, filename)

		path := filepath.Dir(targetPath)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
		}

		perm := 0666
		if strings.HasPrefix(filename, filepath.FromSlash("/plugins/")) {
			perm = 0755
		}

		err = copy(sourcePath, targetPath, perm)
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

	if platform == libvirt.Name {
		_, err = terraform.Apply(tempDir, platform, metadata.BYO, extraArgs...)
		if err != nil {
			return errors.Wrap(err, "Terraform apply")
		}
	}

	extraArgs = append(extraArgs, "-target=module.bootstrap")
	err = terraform.Destroy(tempDir, platform, metadata.BYO, extraArgs...)
	if err != nil {
		return errors.Wrap(err, "Terraform destroy")
	}

	tempStateFilePath := filepath.Join(dir, terraform.StateFileName+".new")
	err = copy(filepath.Join(tempDir, terraform.StateFileName), tempStateFilePath, 0666)
	if err != nil {
		return errors.Wrapf(err, "failed to copy %s from the temporary directory", terraform.StateFileName)
	}
	return os.Rename(tempStateFilePath, filepath.Join(dir, terraform.StateFileName))
}

func copy(from string, to string, perm int) error {
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(to, data, os.FileMode(perm))
}
