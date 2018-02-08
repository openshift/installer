package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Parse parses a yaml string and returns, if successful, a Config.
func Parse(data string) (*Config, error) {
	config := &Config{}

	err := yaml.Unmarshal([]byte(data), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// ParseFile parses a yaml file and returns, if successful, a Config.
func ParseFile(path string) (*Config, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Parse(string(dat))
}
