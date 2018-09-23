package config

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/openshift/installer/pkg/rhcos"
)

// ParseConfig parses a yaml string and returns, if successful, a Cluster.
func ParseConfig(data []byte) (*Cluster, error) {
	cluster := defaultCluster

	if err := yaml.Unmarshal(data, &cluster); err != nil {
		return nil, err
	}

	// Deprecated: remove after openshift/release is ported to pullSecret
	if cluster.PullSecretPath != "" {
		if cluster.PullSecret != "" {
			return nil, errors.New("pullSecretPath is deprecated; just set pullSecret")
		}

		data, err := ioutil.ReadFile(cluster.PullSecretPath)
		if err != nil {
			return nil, err
		}
		cluster.PullSecret = string(data)
		cluster.PullSecretPath = ""
	}

	if cluster.Platform == PlatformAWS && cluster.EC2AMIOverride == "" {
		ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancel()

		ami, err := rhcos.AMI(ctx, rhcos.DefaultChannel, cluster.AWS.Region)
		if err != nil {
			return nil, fmt.Errorf("failed to determine default AMI: %v", err)
		}
		cluster.EC2AMIOverride = ami
	}

	return &cluster, nil
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
