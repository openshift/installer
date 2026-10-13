// Package openstack collects OpenStack-specific configuration.
package openstack

import (
	"fmt"
	"sync"

	"github.com/gophercloud/utils/v2/openstack/clientconfig"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"

	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

var onceLoggers = map[string]*sync.Once{}

// Session is an object representing session for OpenStack.
type Session struct {
	CloudConfig *clientconfig.Cloud
	ClientOpts  *clientconfig.ClientOpts
}

// GetSession returns an OpenStack session for a given cloud name in clouds.yaml.
func GetSession(cloudName string) (*Session, error) {
	opts := openstackdefaults.DefaultClientOpts(cloudName)
	opts.YAMLOpts = new(yamlLoadOpts)

	cloudConfig, err := clientconfig.GetCloudFromYAML(opts)
	if err != nil {
		return nil, err
	}
	return &Session{
		CloudConfig: cloudConfig,
		ClientOpts:  opts,
	}, nil
}

type yamlLoadOpts struct{}

func (opts yamlLoadOpts) LoadCloudsYAML() (map[string]clientconfig.Cloud, error) {
	var clouds clientconfig.Clouds
	content, err := loadAndLog(clientconfig.FindAndReadCloudsYAML)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, &clouds)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return clouds.Clouds, nil
}

func (opts yamlLoadOpts) LoadSecureCloudsYAML() (map[string]clientconfig.Cloud, error) {
	var clouds clientconfig.Clouds
	content, err := loadAndLog(clientconfig.FindAndReadSecureCloudsYAML)
	if err != nil {
		if err.Error() == "no secure.yaml file found" {
			// secure.yaml is optional so just ignore read error
			return clouds.Clouds, nil
		}
	}
	err = yaml.Unmarshal(content, &clouds)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	return clouds.Clouds, err
}

func (opts yamlLoadOpts) LoadPublicCloudsYAML() (map[string]clientconfig.Cloud, error) {
	var publicClouds clientconfig.PublicClouds
	content, err := loadAndLog(clientconfig.FindAndReadPublicCloudsYAML)
	if err != nil {
		if err.Error() == "no clouds-public.yaml file found" {
			// clouds-public.yaml is optional so just ignore read error
			return publicClouds.Clouds, nil
		}
	}
	err = yaml.Unmarshal(content, &publicClouds)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	return publicClouds.Clouds, err
}

func loadAndLog(fn func() (string, []byte, error)) ([]byte, error) {
	filename, content, err := fn()
	if err != nil {
		return nil, err
	}

	if _, has := onceLoggers[filename]; !has {
		onceLoggers[filename] = new(sync.Once)
	}
	onceLoggers[filename].Do(func() {
		logrus.Infof("Credentials loaded from file %q", filename)
	})

	return content, nil
}
