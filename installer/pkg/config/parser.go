package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ParseConfig parses a yaml string and returns, if successful, a Cluster.
func ParseConfig(data []byte) (*Cluster, error) {
	cluster := &Cluster{}

	if err := yaml.Unmarshal(data, cluster); err != nil {
		return nil, err
	}

	return cluster, nil
}

// ParseConfigFile parses a yaml file and returns, if successful, a Cluster.
func ParseConfigFile(path string) (*Cluster, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseConfig(data)
}

// ParseInternal parses a yaml string and returns, if successful, an internal.
func ParseInternal(data []byte) (*Internal, error) {
	internal := &Internal{}

	if err := yaml.Unmarshal(data, internal); err != nil {
		return nil, err
	}

	return internal, nil
}

// ParseInternalFile parses a yaml file and returns, if successful, an internal.
func ParseInternalFile(path string) (*Internal, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseInternal(data)
}
