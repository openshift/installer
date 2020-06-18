// Package packet collects packet-specific configuration.
package packet

import (
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/types/packet"
	packetdefaults "github.com/openshift/installer/pkg/types/packet/defaults"
	"github.com/openshift/installer/pkg/validate"
)

// Platform collects packet-specific configuration.
func Platform() (*packet.Platform, error) {
	var uri string
	err := survey.Ask([]*survey.Question{
		// TODO(displague) ask the right questions
		{
			Prompt: &survey.Input{
				Message: "Packet Connection URI",
				Help:    "The packet connection URI to be used. This must be accessible from the running cluster.",
				Default: packetdefaults.DefaultURI,
			},
			Validate: survey.ComposeValidators(survey.Required, uriValidator),
		},
	}, &uri)
	if err != nil {
		return nil, err
	}

	return &packet.Platform{}, nil
	// TODO(displague) fill in the params
}

// uriValidator validates if the answer provided in prompt is a valid
// url and has non-empty scheme.
func uriValidator(ans interface{}) error {
	return validate.URI(ans.(string))
}
