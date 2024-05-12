package utils

import (
	"fmt"
	"os"

	"github.com/openshift-kni/lifecycle-agent/api/ibiconfig"
)

func ReadIBIConfigFile(configFile string) (*ibiconfig.IBIPrepareConfig, error) {
	var config ibiconfig.IBIPrepareConfig
	if configFile == "" {
		return nil, fmt.Errorf("configuration file is required")
	}
	if _, err := os.Stat(configFile); err != nil {
		return nil, fmt.Errorf("configuration file %s does not exist", configFile)
	}

	if err := ReadYamlOrJSONFile(configFile, &config); err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configFile, err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config file validation failed: %w", err)
	}

	config.SetDefaultValues()

	return &config, nil
}
