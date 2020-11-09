package ovirt

import (
	ovirtsdk4 "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/types/ovirt"
)

// Platform collects ovirt-specific configuration.
func Platform() (*ovirt.Platform, error) {
	p := ovirt.Platform{}

	ovirtConfig, err := NewConfig()
	if err != nil {
		ovirtConfig, err = engineSetup()
		if err != nil {
			return nil, err
		}
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
	ovirtConfig.Save()

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

	err = askVNICProfileID(c, &p)
	if err != nil {
		return &p, err
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Internal API virtual IP",
				Help:    "This is the virtual IP address that will be used to address the OpenShift control plane. Make sure the IP address is not in use.",
				Default: "",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &p.APIVIP)
	if err != nil {
		return nil, errors.Wrap(err, "failed UserInput")
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Ingress virtual IP",
				Help:    "This is the virtual IP address that will be used to address the OpenShift ingress routers. Make sure the IP address is not in use.",
				Default: "",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &p.IngressVIP)
	if err != nil {
		return nil, errors.Wrap(err, "failed UserInput")
	}

	return &p, nil
}
