package terraform

import (
	"fmt"
	"os"
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
// It outputs the tfstate file.
func Apply(clusterDir string, state string, templateDir string, extraArgs ...string) error {
	defaultArgs := []string{
		"apply",
		"-auto-approve",
		fmt.Sprintf("-state=%s.tfstate", state),
	}
	extraArgs = append(extraArgs, templateDir)
	args := append(defaultArgs, extraArgs...)
	return terraformExec(clusterDir, args...)
}

// Destroy runs "terraform destory" with the given clusterDir, templateDir, state file
// and extra args.
func Destroy(clusterDir, state, templateDir string, extraArgs ...string) error {
	defaultArgs := []string{
		"destroy",
		"-force",
		fmt.Sprintf("-state=%s.tfstate", state),
	}
	extraArgs = append(extraArgs, templateDir)
	args := append(defaultArgs, extraArgs...)
	return terraformExec(clusterDir, args...)
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
