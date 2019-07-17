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
	"github.com/openshift/installer/pkg/types/gcp"
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

	tempDir, err := ioutil.TempDir("", "openshift-install-")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary directory for Terraform execution")
	}
	defer os.RemoveAll(tempDir)

	extraArgs, err := setup(dir, tempDir, platform)
	if err != nil {
		return err
	}

	switch platform {
	case gcp.Name:
		// Then destory the bootstrap instance and instance group so destroy runs cleanly.
		_, err = terraform.Apply(tempDir, platform, append(extraArgs, "-var=gcp_bootstrap_enabled=false")...)
		if err != nil {
			return errors.Wrap(err, "Terraform apply")
		}
	case libvirt.Name:
		// First remove the bootstrap node from DNS
		_, err = terraform.Apply(tempDir, platform, append(extraArgs, "-var=bootstrap_dns=false")...)
		if err != nil {
			return errors.Wrap(err, "Terraform apply")
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

// Unlink uses Terraform to remove bootstrap resources from Load balancing.
// Platforms like GCP require bootstrap to be remove from LB backends for bootstrapping to be finished.
func Unlink(dir string) error {
	metadata, err := cluster.LoadMetadata(dir)
	if err != nil {
		return err
	}

	platform := metadata.Platform()
	if platform != gcp.Name {
		return nil
	}

	tempDir, err := ioutil.TempDir("", "openshift-install-")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary directory for Terraform execution")
	}
	defer os.RemoveAll(tempDir)

	extraArgs, err := setup(dir, tempDir, platform)
	if err != nil {
		return err
	}

	switch platform {
	case gcp.Name:
		// First remove the bootstrap node from the load balancers to avoid race condition.
		_, err = terraform.Apply(tempDir, platform, append(extraArgs, "-var=gcp_bootstrap_lb=false")...)
		if err != nil {
			return errors.Wrap(err, "Terraform apply")
		}
	}

	tempStateFilePath := filepath.Join(dir, terraform.StateFileName+".new")
	err = copy(filepath.Join(tempDir, terraform.StateFileName), tempStateFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to copy %s from the temporary directory", terraform.StateFileName)
	}
	return os.Rename(tempStateFilePath, filepath.Join(dir, terraform.StateFileName))
}

func setup(src, dst string, platform string) ([]string, error) {
	retArgs := []string{}
	tfPlatformVarsFileName := fmt.Sprintf(cluster.TfPlatformVarsFileName, platform)
	for _, filename := range []string{terraform.StateFileName, cluster.TfVarsFileName, tfPlatformVarsFileName} {
		sourcePath := filepath.Join(src, filename)
		targetPath := filepath.Join(dst, filename)
		err := copy(sourcePath, targetPath)
		if err != nil {
			if os.IsNotExist(err) && err.(*os.PathError).Path == sourcePath && filename == tfPlatformVarsFileName {
				continue // platform may not need platform-specific Terraform variables
			}
			return retArgs, errors.Wrapf(err, "failed to copy %s to the temporary directory", filename)
		}
		if strings.HasSuffix(filename, ".tfvars.json") {
			retArgs = append(retArgs, fmt.Sprintf("-var-file=%s", targetPath))
		}
	}
	return retArgs, nil
}

func copy(from string, to string) error {
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(to, data, 0666)
}
