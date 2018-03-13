package config

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Error codes returned by failures to parse a config.
var (
	ErrMultipleClusters = errors.New("multiple cluster configurations are not currently supported")
	ErrNoClusters       = errors.New("no clusters were defined")
)

// Parse parses a yaml string and returns, if successful, a Config.
func Parse(data []byte) (*Config, error) {
	config := &Config{}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	if len(config.Clusters) == 0 {
		return config, ErrNoClusters
	}

	if len(config.Clusters) > 1 {
		return config, ErrMultipleClusters
	}

	return config, nil
}

// ParseFile parses a yaml file and returns, if successful, a Config.
func ParseFile(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Parse(data)
}
