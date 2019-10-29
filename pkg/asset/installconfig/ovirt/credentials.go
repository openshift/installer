package ovirt

import (
	"gopkg.in/AlecAivazis/survey.v1"
)

func askCredentials() (Config, error) {
	c := Config{}
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Enter oVirt's api endpoint URL",
				Help:    "oVirt engine api url, for example https://ovirt-engine-fqdn/ovirt-engine/api",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &c.URL)
	if err != nil {
		return c, err
	}

	var ovirtCertTrusted bool
	err = survey.AskOne(
		&survey.Confirm{
			Message: "Is the installed oVirt certificate trusted?",
			Default: false,
			Help:    "",
		},
		&ovirtCertTrusted,
		nil)
	if err != nil {
		return c, err
	}
	c.Insecure = !ovirtCertTrusted

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Enter ovirt-engine username",
				Help:    "The user must have permissions to create VMs and disks on the Storage Domain with the same name as the OpenShift cluster",
				Default: "admin@internal",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &c.Username)
	if err != nil {
		return c, err
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "Enter password",
				Help:    "",
			},
			Validate: survey.ComposeValidators(survey.Required, authenticated(&c)),
		},
	}, &c.Password)
	if err != nil {
		return c, err
	}

	return c, nil
}
