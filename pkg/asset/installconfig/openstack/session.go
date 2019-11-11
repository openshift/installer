// Package openstack collects OpenStack-specific configuration.
package openstack

import (
	"github.com/pkg/errors"

	"github.com/ghodss/yaml"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sirupsen/logrus"
)

// Session is an object representing session for OpenStack.
type Session struct {
	CloudConfig *clientconfig.Cloud
}

// GetSession returns an OpenStack session for a given cloud name in clouds.yaml.
func GetSession(cloudName string) (*Session, error) {
	opts := defaultClientOpts(cloudName)
	cloudConfig, err := clientconfig.GetCloudFromYAML(opts)
	if err != nil {
		return nil, err
	}
	return &Session{
		CloudConfig: cloudConfig,
	}, nil
}

func defaultClientOpts(cloudName string) *clientconfig.ClientOpts {
	opts := new(clientconfig.ClientOpts)
	opts.Cloud = cloudName
	opts.YAMLOpts = new(yamlLoadOpts)
	return opts
}

type yamlLoadOpts struct{}

func (opts yamlLoadOpts) LoadCloudsYAML() (map[string]clientconfig.Cloud, error) {
	return loadAndLog(clientconfig.FindAndReadCloudsYAML)
}

func (opts yamlLoadOpts) LoadSecureCloudsYAML() (map[string]clientconfig.Cloud, error) {
	clouds, err := loadAndLog(clientconfig.FindAndReadSecureCloudsYAML)
	if err != nil {
		if err.Error() == "no secure.yaml file found" {
			// secure.yaml is optional so just ignore read error
			return clouds, nil
		}
	}
	return clouds, err
}

func (opts yamlLoadOpts) LoadPublicCloudsYAML() (map[string]clientconfig.Cloud, error) {
	clouds, err := loadAndLog(clientconfig.FindAndReadPublicCloudsYAML)
	if err != nil {
		if err.Error() == "no clouds-public.yaml file found" {
			// clouds-public.yaml is optional so just ignore read error
			return clouds, nil
		}
	}
	return clouds, err
}

func loadAndLog(fn func() (string, []byte, error)) (map[string]clientconfig.Cloud, error) {
	filename, content, err := fn()
	if err != nil {
		return nil, err
	}

	var clouds clientconfig.Clouds
	err = yaml.Unmarshal(content, &clouds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal yaml")
	}

	logrus.Infof("Credentials loaded from file %q", filename)

	return clouds.Clouds, nil
}
