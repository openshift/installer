// Package bootstrap uses Terraform to remove bootstrap resources.
package bootstrap

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/openshift/installer/pkg/asset/cluster"
	osp "github.com/openshift/installer/pkg/destroy/openstack"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/openstack"
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

	tempDir, err := ioutil.TempDir("", "openshift-install-")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary directory for Terraform execution")
	}
	defer os.RemoveAll(tempDir)

	extraArgs := []string{}
	for _, filename := range []string{terraform.StateFileName, cluster.TfVarsFileName, tfPlatformVarsFileName} {
		sourcePath := filepath.Join(dir, filename)
		targetPath := filepath.Join(tempDir, filename)
		err = copy(sourcePath, targetPath)
		if err != nil {
			if os.IsNotExist(err) && err.(*os.PathError).Path == sourcePath && filename == tfPlatformVarsFileName {
				continue // platform may not need platform-specific Terraform variables
			}
			return errors.Wrapf(err, "failed to copy %s to the temporary directory", filename)
		}
		if strings.HasSuffix(filename, ".tfvars.json") {
			extraArgs = append(extraArgs, fmt.Sprintf("-var-file=%s", targetPath))
		}
	}

	switch platform {
	case gcp.Name:
		// Remove the bootstrap machine from load balancers first before destroying the bootstrap resources. This is
		// necessary because terraform does not inherently know that the ordering of removing from the load balancers
		// first is necessary.
		_, err = terraform.Apply(tempDir, platform, append(extraArgs, "-var=gcp_bootstrap_included_in_lb=false")...)
		if err != nil {
			return errors.Wrap(err, "failed to remove bootstrap from load balancers")
		}
	case openstack.Name:
		imageName := metadata.InfraID + "-ignition"
		err = osp.DeleteGlanceImage(imageName, metadata.OpenStack.Cloud)
		if err != nil {
			return errors.Wrapf(err, "Failed to delete glance image %s", imageName)
		}
	}

	_, err = terraform.Apply(tempDir, platform, append(extraArgs, "-var=bootstrapping=false")...)
	if err != nil {
		return errors.Wrap(err, "failed disabling bootstrap")
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
