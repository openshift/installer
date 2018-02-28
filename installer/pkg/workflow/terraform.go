package workflow

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func terraformExec(clusterDir string, args ...string) error {
	tf := exec.Command("terraform", args...)
	tf.Dir = clusterDir
	tf.Stdin = os.Stdin
	tf.Stdout = os.Stdout
	tf.Stderr = os.Stderr

	return tf.Run()
}

func tfApply(clusterDir, state, templateDir string) error {
	return terraformExec(clusterDir, "apply", "-auto-approve", fmt.Sprintf("-state=%s.tfstate", state), templateDir)
}

func tfDestroy(clusterDir, state, templateDir string) error {
	return terraformExec(clusterDir, "destroy", "-force", fmt.Sprintf("-state=%s.tfstate", state), templateDir)
}

func tfInit(clusterDir, templateDir string) error {
	return terraformExec(clusterDir, "init", templateDir)
}

func hasStateFile(stateDir string, stateName string) bool {
	stepStateFile := filepath.Join(stateDir, fmt.Sprintf("%s.tfstate", stateName))
	_, err := os.Stat(stepStateFile)
	return !os.IsNotExist(err)
}
