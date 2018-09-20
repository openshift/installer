package terraform

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/openshift/installer/pkg/types/config"
)

const (
	// AssetsStep is the name of the step that generates the assets.
	// This is deprecated.
	// TODO(yifan) Remove this when removing the asset step.
	AssetsStep = "assets"

	// BootstrapStep is the name of the step that runs the bootstrap.
	BootstrapStep = "bootstrap"
	// InfraStep is the name of the step that sets up infra.
	InfraStep = "infra"

	stepsBaseDir = "steps"
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
func Apply(clusterDir string, state string, templateDir string, extraArgs ...string) (string, error) {
	stateFileName := fmt.Sprintf("%s.tfstate", state)
	defaultArgs := []string{
		"apply",
		"-auto-approve",
		fmt.Sprintf("-state=%s", stateFileName),
	}
	extraArgs = append(extraArgs, templateDir)
	args := append(defaultArgs, extraArgs...)

	return path.Join(clusterDir, stateFileName), terraformExec(clusterDir, args...)
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

// FindStepTemplates returns the directory containing templates for a given step.
// If platform is specified, it looks for a subdirectory with platform first,
// falling back if there are no platform-specific templates for that step.
func FindStepTemplates(dir, stepName string, platform config.Platform) (string, error) {
	stepDir := filepath.Join(dir, stepsBaseDir, stepName)
	for _, path := range []string{
		filepath.Join(stepDir, string(platform)),
		stepDir} {

		stat, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", fmt.Errorf("invalid path for %q templates: %s", path, err)
		}
		if !stat.IsDir() {
			return "", fmt.Errorf("invalid path for %q templates", path)
		}
		return path, nil
	}
	return "", os.ErrNotExist
}
