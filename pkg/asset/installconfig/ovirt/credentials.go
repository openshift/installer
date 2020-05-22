package ovirt

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"
)

// Check if URL can be reached before we proceed with the installation
// Parms:
//	urlAddr - Full URL
func checkURLResponse(urlAddr string) {

	logrus.Debug("Checking URL Response: ", urlAddr)
	_, err := http.Get(urlAddr)
	if err != nil {
		logrus.Fatal(err)
	}
}

func askCredentials() (Config, error) {
	oVirtConfig := Config{}
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "oVirt FQDN[:PORT]",
				Help:    "The oVirt FQDN[:PORT] (ovirt.example.com:443)",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &oVirtConfig.FQDN)
	if err != nil {
		return oVirtConfig, err
	}


	// Set c.URL with the API endpoint
	oVirtConfig.URL = fmt.Sprintf(
		"https://%s/ovirt-engine/api",
		oVirtConfig.FQDN)

	checkURLResponse(oVirtConfig.URL)

	var ovirtCertTrusted bool
	err = survey.AskOne(
		&survey.Confirm{
			Message: "Is the oVirt CA trusted locally?",
			Default: true,
			Help:    "In order to securly communicate with the oVirt engine, the certificate authority must be trusted by the local system.",
		},
		&ovirtCertTrusted,
		nil)
	if err != nil {
		return oVirtConfig, err
	}
	oVirtConfig.Insecure = !ovirtCertTrusted

	if ovirtCertTrusted {
		if err != nil {
			// should have passed validation, this is unexpected
			return oVirtConfig, err
		}

		oVirtConfig.PemURL = fmt.Sprintf(
			"https://%s/ovirt-engine/services/pki-resource?resource=ca-certificate&format=X509-PEM-CA",
			oVirtConfig.FQDN)
		logrus.Debug("PEM URL: ", oVirtConfig.PemURL)

		err = survey.AskOne(&survey.Multiline{
			Message: "oVirt certificate bundle",
			Help:    fmt.Sprintf("The oVirt certificate bundle can be downloaded from %s", oVirtConfig.PemURL),
		},
			&oVirtConfig.CABundle,
			survey.ComposeValidators(survey.Required))
		if err != nil {
			return oVirtConfig, err
		}
	} else {
		logrus.Warning("Communication with the oVirt engine will be insecure.")
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
	}, &oVirtConfig.Username)
	if err != nil {
		return oVirtConfig, err
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "oVirt engine password",
				Help:    "",
			},
			Validate: survey.ComposeValidators(survey.Required, authenticated(&oVirtConfig)),
		},
	}, &oVirtConfig.Password)
	if err != nil {
		return oVirtConfig, err
	}

	return oVirtConfig, nil
}
