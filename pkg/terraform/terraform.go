package terraform

import (
	"fmt"
	"path/filepath"

	"github.com/openshift/installer/data"
	"github.com/pkg/errors"
)

const (
	// StateFileName is the default name for Terraform state files.
	StateFileName string = "terraform.tfstate"
)

func terraformExec(clusterDir string, args ...string) error {
	// Create an executor
	ex, err := newExecutor()
	if err != nil {
		return errors.Wrap(err, "failed to create Terraform executor")
	}

	err = ex.execute(clusterDir, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute Terraform")
	}
	return nil
}

// Apply unpacks the platform-specific Terraform modules into the
// given directory and then runs 'terraform init' and 'terraform
// apply'.  It returns the absolute path of the tfstate file, rooted
// in the specified directory, along with any errors from Terraform.
func Apply(dir string, platform string, extraArgs ...string) (path string, err error) {
	err = unpackAndInit(dir, platform)
	if err != nil {
		return "", err
	}

	defaultArgs := []string{
		"apply",
		"-auto-approve",
		"-input=false",
		"-no-color",
		fmt.Sprintf("-state=%s", StateFileName),
	}
	args := append(defaultArgs, extraArgs...)

	return filepath.Join(dir, StateFileName), terraformExec(dir, args...)
}

// Destroy unpacks the platform-specific Terraform modules into the
// given directory and then runs 'terraform init' and 'terraform
// destroy'.
func Destroy(dir string, platform string, extraArgs ...string) (err error) {
	err = unpackAndInit(dir, platform)
	if err != nil {
		return err
	}

	defaultArgs := []string{
		"destroy",
		"-auto-approve",
		"-no-color",
		fmt.Sprintf("-state=%s", StateFileName),
	}
	args := append(defaultArgs, extraArgs...)

	return terraformExec(dir, args...)
}

// unpack unpacks the platform-specific Terraform modules into the
// given directory.
func unpack(dir string, platform string) (err error) {
	err = data.Unpack(dir, platform)
	if err != nil {
		return err
	}

	err = data.Unpack(filepath.Join(dir, "config.tf"), "config.tf")
	if err != nil {
		return err
	}

	return nil
}

// unpackAndInit unpacks the platform-specific Terraform modules into
// the given directory and then runs 'terraform init'.
func unpackAndInit(dir string, platform string) (err error) {
	err = unpack(dir, platform)
	if err != nil {
		return errors.Wrap(err, "failed to unpack Terraform modules")
	}

	err = terraformExec(dir, "init", "-input=false", "-no-color")
	if err != nil {
		return errors.Wrap(err, "failed to initialize Terraform")
	}

	return nil
}
