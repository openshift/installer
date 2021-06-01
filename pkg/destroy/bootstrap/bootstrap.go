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
	platformstages "github.com/openshift/installer/pkg/terraform/stages/platform"
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

	tfPlatformVarsFileName := fmt.Sprintf(cluster.TfPlatformVarsFileName, platform)
	varFiles := []string{cluster.TfVarsFileName, tfPlatformVarsFileName}
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

		if err := copyToTemp(stage.StateFilename(), dir, tempDir, false); err != nil {
			return err
		}

		extraArgs := make([]string, len(varFiles))
		for i, filename := range varFiles {
			allowMissing := filename == tfPlatformVarsFileName // platform may not need platform-specific Terraform variables
			if err := copyToTemp(filename, dir, tempDir, allowMissing); err != nil {
				return err
			}
			extraArgs[i] = fmt.Sprintf("-var-file=%s", filepath.Join(tempDir, filename))
		}

		if err := stage.Destroy(tempDir, extraArgs); err != nil {
			return err
		}

		tempStateFilePath := filepath.Join(dir, stage.StateFilename()+".new")
		err = copy(filepath.Join(tempDir, stage.StateFilename()), tempStateFilePath)
		if err != nil {
	}

	return nil
}

func copyToTemp(filename, sourceDir, tempDir string, allowMissing bool) error {
	sourcePath := filepath.Join(sourceDir, filename)
	targetPath := filepath.Join(tempDir, filename)
	err := copy(sourcePath, targetPath)
	if err != nil {
		if os.IsNotExist(err) && err.(*os.PathError).Path == sourcePath && allowMissing {
			return nil
		}
		return errors.Wrapf(err, "failed to copy %s to the temporary directory", filename)
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
