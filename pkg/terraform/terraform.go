package terraform

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
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

// Apply runs "terraform apply" with the given clusterDir, templateDir and extra args.
// It returns the absolute path of the tfstate file.
func Apply(clusterDir string, templateDir string, extraArgs ...string) (string, error) {
	stateFileName := "terraform.tfstate"
	defaultArgs := []string{
		"apply",
		"-auto-approve",
		fmt.Sprintf("-state=%s", stateFileName),
	}
	extraArgs = append(extraArgs, templateDir)
	args := append(defaultArgs, extraArgs...)

	return path.Join(clusterDir, stateFileName), terraformExec(clusterDir, args...)
}

// Init runs "terraform init" with the given clusterDir and templateDir.
func Init(clusterDir, templateDir string) error {
	return terraformExec(clusterDir, "init", templateDir)
}

// HasStateFile returns true if the stateFile exists under stateDir.
func HasStateFile(stateDir string, stateFile string) bool {
	stepStateFile := filepath.Join(stateDir, fmt.Sprintf("%s.tfstate", stateFile))
	_, err := os.Stat(stepStateFile)
	return !os.IsNotExist(err)
}
