package ovirt

import (
	ovirtsdk4 "github.com/ovirt/go-ovirt"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/types/ovirt"
)

// Platform collects ovirt-specific configuration.
func Platform() (*ovirt.Platform, error) {
	p := ovirt.Platform{}

	ovirtConfig, err := NewConfig()
	if err != nil {
		ovirtConfig, err = askCredentials()
		if err != nil {
			return nil, err
		}
		defer ovirtConfig.Save()
	}

	c, err := ovirtsdk4.NewConnectionBuilder().
		URL(ovirtConfig.URL).
		Username(ovirtConfig.Username).
		Password(ovirtConfig.Password).
		CAFile(ovirtConfig.CAFile).
		Insecure(ovirtConfig.Insecure).
		Build()

	if err != nil {
		return nil, err
	}
	defer c.Close()
	err = c.Test()
	if err != nil {
		return nil, err
	}

	clusterName, err := askCluster(c, &p)
	if err != nil {
		return &p, err
	}

	err = askStorage(c, &p, clusterName)
	if err != nil {
		return &p, err
	}

	err = askNetwork(c, &p)
	if err != nil {
		return &p, err
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Enter the internal API Virtual IP",
				Help:    "Make sure the IP address is not in use",
				Default: "",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &p.APIVIP)
	if err != nil {
		return nil, err
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Enter the internal DNS Virtual IP",
				Help:    "Make sure the IP address is not in use",
				Default: "",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &p.DNSVIP)
	if err != nil {
		return nil, err
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Enter the ingress IP ",
				Help:    "Make sure the IP address is not in use",
				Default: "",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &p.IngressVIP)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
