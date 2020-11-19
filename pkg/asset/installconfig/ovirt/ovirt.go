package ovirt

import (
	ovirtsdk4 "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/types/ovirt"
)

const platformValidationMaxTries = 3

// Platform collects ovirt-specific configuration.
func Platform() (*ovirt.Platform, error) {
	p := ovirt.Platform{}

	var c *ovirtsdk4.Connection

	// Fetch config from file
	ovirtConfig, err := NewConfig()
	for tries := 0; tries < platformValidationMaxTries; tries++ {
		if err == nil {
			// If no error happened previously (loading file or configuration), validate the connection.
			c, err = ovirtConfig.getValidatedConnection()
			if err != nil {
				// If validation failed log and drop into reconfig below.
				logrus.Error(errors.Wrap(err, "failed to validate oVirt configuration"))
			} else {
				// If connection is valid, break
				break
			}
		}

		if err != nil {
			// If a previous error happened (validation or loading from file), rerun the setup.
			ovirtConfig, err = engineSetup()
			if err != nil {
				// If validation failed log and loop back to engineSetup()
				logrus.Error(errors.Wrap(err, "oVirt configuration failed"))
			}
		}
	}
	if err != nil {
		// Last error is not nil, we don't have a valid config.
		return nil, errors.Wrap(err, "maximum retries for configuration exhausted")
	}
	defer c.Close()
	if err = ovirtConfig.Save(); err != nil {
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
