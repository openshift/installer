package baremetal

import (
	"github.com/AlecAivazis/survey/v2"

	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/validate"
)

// Host prompts the user for hardware details about a baremetal host.
func Host() (*baremetal.Host, error) {
	var host baremetal.Host

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Name",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &host.Name); err != nil {
		return nil, err
	}

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "BMC Address",
				Help:    "The address for the BMC, e.g. ipmi://192.168.0.1",
			},
			Validate: survey.ComposeValidators(survey.Required, uriValidator),
		},
	}, &host.BMC.Address); err != nil {
		return nil, err
	}

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "BMC Username",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &host.BMC.Username); err != nil {
		return nil, err
	}

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "BMC Password",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &host.BMC.Password); err != nil {
		return nil, err
	}

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Boot MAC Address",
			},
			Validate: survey.ComposeValidators(survey.Required, macValidator),
		},
	}, &host.BootMACAddress); err != nil {
		return nil, err
	}

	return &host, nil
}

// uriValidator validates if the answer provided in prompt is a valid
// url and has non-empty scheme.
func uriValidator(ans interface{}) error {
	return validate.URI(ans.(string))
}

// macValidator validates if the answer provided is a valid mac address
func macValidator(ans interface{}) error {
	return validate.MAC(ans.(string))
}
