package preprovision

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// SetTerraformEnvironment injects the environment variables required
// authenticate against the OpenStack API.
func SetTerraformEnvironment() error {
	// Terraform runs in a different directory but we want to allow people to
	// use clouds.yaml files in their local directory. Emulate this by setting
	// the necessary environment variable to point to this file if (a) the user
	// hasn't already set this environment variable and (b) there is actually
	// a local file
	if path := os.Getenv("OS_CLIENT_CONFIG_FILE"); path != "" {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to determine working directory: %w", err)
	}

	cloudsYAML := filepath.Join(cwd, "clouds.yaml")
	if _, err = os.Stat(cloudsYAML); err == nil {
		os.Setenv("OS_CLIENT_CONFIG_FILE", cloudsYAML)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("unable to determine if clouds.yaml exists: %w", err)
	}

	return nil
}
