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

// ParseConfig parses a yaml string and returns, if successful, a Config.
func ParseConfig(data []byte) (*Config, error) {
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

// ParseConfigFile parses a yaml file and returns, if successful, a Config.
func ParseConfigFile(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseConfig(data)
}

// ParseInternal parses a yaml string and returns, if successful, an internal.
func ParseInternal(data []byte) (*internal, error) {
	internal := &internal{}

	if err := yaml.Unmarshal(data, internal); err != nil {
		return nil, err
	}

	return internal, nil
}

// ParseInternalFile parses a yaml file and returns, if successful, an internal.
func ParseInternalFile(path string) (*internal, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseInternal(data)
}
