// Package bootstrap uses Terraform to remove bootstrap resources.
package bootstrap

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset/cluster"
	osp "github.com/openshift/installer/pkg/destroy/openstack"
	"github.com/openshift/installer/pkg/terraform"
	platformstages "github.com/openshift/installer/pkg/terraform/stages/platform"
	typesazure "github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/openstack"
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

	if platform == openstack.Name {
		imageName := metadata.InfraID + "-ignition"
		if err := osp.DeleteGlanceImage(imageName, metadata.OpenStack.Cloud); err != nil {
			return errors.Wrapf(err, "Failed to delete glance image %s", imageName)
		}
	}

	// Azure Stack uses the Azure platform but has its own Terraform configuration.
	if platform == typesazure.Name && metadata.Azure.CloudName == typesazure.StackCloud {
		platform = typesazure.StackTerraformName
	}

	varFiles := []string{cluster.TfVarsFileName, cluster.TfPlatformVarsFileName}
	tfStages := platformstages.StagesForPlatform(platform)
	for _, stage := range tfStages {
		varFiles = append(varFiles, stage.OutputsFilename())
	}

	for i := len(tfStages) - 1; i >= 0; i-- {
		stage := tfStages[i]

		if !stage.DestroyWithBootstrap() {
			continue
		}

		tempDir, err := ioutil.TempDir("", fmt.Sprintf("openshift-install-%s-", stage.Name()))
		if err != nil {
			return errors.Wrap(err, "failed to create temporary directory for Terraform execution")
		}
		defer os.RemoveAll(tempDir)

		stateFilePathInInstallDir := filepath.Join(dir, stage.StateFilename())
		stateFilePathInTempDir := filepath.Join(tempDir, terraform.StateFilename)
		if err := copy(stateFilePathInInstallDir, stateFilePathInTempDir); err != nil {
			return errors.Wrap(err, "failed to copy state file to the temporary directory")
		}

		targetVarFiles := make([]string, 0, len(varFiles))
		for _, filename := range varFiles {
			sourcePath := filepath.Join(dir, filename)
			targetPath := filepath.Join(tempDir, filename)
			if err := copy(sourcePath, targetPath); err != nil {
				// platform may not need platform-specific Terraform variables
				if filename == cluster.TfPlatformVarsFileName {
					if os.IsNotExist(err) && err.(*os.PathError).Path == sourcePath {
						continue
					}
				}
				return errors.Wrapf(err, "failed to copy %s to the temporary directory", filename)
			}
			targetVarFiles = append(targetVarFiles, targetPath)
		}

		if err := stage.Destroy(tempDir, targetVarFiles); err != nil {
			return err
		}

		if err := copy(stateFilePathInTempDir, stateFilePathInInstallDir); err != nil {
			return errors.Wrap(err, "failed to copy state file from the temporary directory")
		}
	}

	return nil
}

func copy(from string, to string) error {
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(to, data, 0666)
}
