package terraform

import (
	"fmt"
	"path"
)

func terraformExec(clusterDir string, args ...string) error {
	// Create an executor
	ex, err := newExecutor()
	if err != nil {
		return fmt.Errorf("could not create Terraform executor: %s", err)
	}

	err = ex.execute(clusterDir, args...)
	if err != nil {
		return fmt.Errorf("failed to run Terraform: %s", err)
	}
	return nil
}

// Apply runs "terraform apply" in the given directory.
// It returns the absolute path of the tfstate file.
func Apply(dir string, extraArgs ...string) (string, error) {
	stateFileName := "terraform.tfstate"
	defaultArgs := []string{
		"apply",
		"-auto-approve",
		fmt.Sprintf("-state=%s", stateFileName),
	}
	args := append(defaultArgs, extraArgs...)

	return path.Join(dir, stateFileName), terraformExec(dir, args...)
}

// Init runs "terraform init" in the given directory.
func Init(dir string) error {
	return terraformExec(dir, "init")
}
