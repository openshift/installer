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
	"github.com/openshift/installer/pkg/terraform/exec/plugins"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
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

	target := terraform.TargetBootstrap
	tfPlugin, err := plugins.GetPlugin(platform)
	for _, item := range tfPlugin.Resources {
		if item == terraform.TargetCompat {
			target = terraform.TargetCompat
			break
		}
	}

	tfPlatformVarsFileName := fmt.Sprintf(cluster.TfPlatformVarsFileName, platform)

	tempDir, err := ioutil.TempDir("", "openshift-install-bootstrap-")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary directory for Terraform execution")
	}
	defer os.RemoveAll(tempDir)

	extraArgs := []string{}
	stateFileName := terraform.GetStateFileName(target)
	for _, filename := range []string{stateFileName, cluster.TfVarsFileName, tfPlatformVarsFileName} {
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
		// First remove the bootstrap node from the load balancers to avoid race condition.
		_, err = terraform.Apply(tempDir, platform, target, append(extraArgs, "-var=gcp_bootstrap_lb=false")...)
		if err != nil {
			return errors.Wrap(err, "failed disabling bootstrap load balancing")
		}

		// Then destory the bootstrap instance and instance group so destroy runs cleanly.
		// First remove the bootstrap from LB target and its instance so that bootstrap module is cleanly destroyed.
		_, err = terraform.Apply(tempDir, platform, target, append(extraArgs, "-var=gcp_bootstrap_enabled=false")...)
		if err != nil {
			return errors.Wrap(err, "failed disabling bootstrap")
		}
	case libvirt.Name:
		// First remove the bootstrap node from DNS
		_, err = terraform.Apply(tempDir, platform, target, append(extraArgs, "-var=bootstrap_dns=false")...)
		if err != nil {
			return errors.Wrap(err, "Terraform apply")
		}
	case openstack.Name:
		imageName := metadata.InfraID + "-ignition"
		err = osp.DeleteGlanceImage(imageName, metadata.OpenStack.Cloud)
		if err != nil {
			return errors.Wrapf(err, "Failed to delete glance image %s", imageName)
		}
	case ovirt.Name:
		extraArgs = append(extraArgs, "-target=module.template.ovirt_vm.tmp_import_vm")
		extraArgs = append(extraArgs, "-target=module.template.ovirt_image_transfer.releaseimage")
	}

	if target == terraform.TargetCompat {
		extraArgs = append(extraArgs, "-target=module.bootstrap")
	}

	err = terraform.Destroy(tempDir, platform, target, extraArgs...)
	if err != nil {
		return errors.Wrap(err, "Terraform destroy")
	}

	tempStateFilePath := filepath.Join(dir, stateFileName+".new")
	err = copy(filepath.Join(tempDir, stateFileName), tempStateFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to copy %s from the temporary directory", stateFileName)
	}
	return os.Rename(tempStateFilePath, filepath.Join(dir, stateFileName))
}

func copy(from string, to string) error {
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(to, data, 0666)
}
