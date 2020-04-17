package ovirt

import (
	"fmt"
	"net/url"

	"gopkg.in/AlecAivazis/survey.v1"
)

func askCredentials() (Config, error) {
	c := Config{}
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "oVirt API endpoint URL",
				Help:    "The URL of the oVirt engine API. For example, https://ovirt-engine-fqdn/ovirt-engine/api.",
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
			Message: "Is the oVirt CA bundle trusted locally?",
			Default: true,
			Help:    "In order to securly communicate with the oVirt engine, the certificate authority bundle must either be trusted by the local system or explicitly provided.",
		},
		&ovirtCertTrusted,
		nil)
	if err != nil {
		return c, err
	}
	c.Insecure = !ovirtCertTrusted

	if ovirtCertTrusted {
		ovirtURL, err := url.Parse(c.URL)
		if err != nil {
			// should have passed validation, this is unexpected
			return c, err
		}
		pemURL := fmt.Sprintf(
			"%s://%s/ovirt-engine/services/pki-resource?resource=ca-certificate&format=X509-PEM-CA",
			ovirtURL.Scheme,
			ovirtURL.Host)

		err = survey.AskOne(&survey.Multiline{
			Message: "oVirt CA bundle",
			Help:    fmt.Sprintf("The oVirt CA bundle can be downloaded from %s.", pemURL),
		},
			&c.CABundle,
			survey.ComposeValidators(survey.Required))
		if err != nil {
			return c, err
		}
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "oVirt engine username",
				Help:    "The user must have permissions to create VMs and disks on the Storage Domain with the same name as the OpenShift cluster.",
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
				Message: "oVirt engine password",
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
