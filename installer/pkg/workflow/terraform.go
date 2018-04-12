package workflow

import (
	"fmt"
	"os"
	"path/filepath"
)

func terraformExec(clusterDir string, args ...string) error {
	// Create an executor
	ex, err := newExecutor()
	if err != nil {
		return fmt.Errorf("Could not create Terraform executor: %s", err)
	}

	err = ex.execute(clusterDir, args...)
	if err != nil {
		return fmt.Errorf("Failed to run Terraform: %s", err)
	}
	return nil
}

func tfApply(clusterDir string, state string, templateDir string, extraArgs ...string) error {
	defaultArgs := []string{
		"apply",
		"-auto-approve",
		fmt.Sprintf("-state=%s.tfstate", state),
	}
	extraArgs = append(extraArgs, templateDir)
	args := append(defaultArgs, extraArgs...)
	return terraformExec(clusterDir, args...)
}

func tfDestroy(clusterDir, state, templateDir string, extraArgs ...string) error {
	defaultArgs := []string{
		"destroy",
		"-force",
		fmt.Sprintf("-state=%s.tfstate", state),
	}
	extraArgs = append(extraArgs, templateDir)
	args := append(defaultArgs, extraArgs...)
	return terraformExec(clusterDir, args...)
}

func tfInit(clusterDir, templateDir string) error {
	return terraformExec(clusterDir, "init", templateDir)
}

func hasStateFile(stateDir string, stateName string) bool {
	stepStateFile := filepath.Join(stateDir, fmt.Sprintf("%s.tfstate", stateName))
	_, err := os.Stat(stepStateFile)
	return !os.IsNotExist(err)
}
