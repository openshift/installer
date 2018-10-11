package installconfig

import (
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/validate"
)

type emailAddress struct {
	EmailAddress string
}

var _ asset.Asset = (*emailAddress)(nil)

// Dependencies returns no dependencies.
func (a *emailAddress) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate queries for the email address from the user.
func (a *emailAddress) Generate(asset.Parents) error {
	email, err := asset.GenerateUserProvidedAsset(
		a.Name(),
		&survey.Question{
			Prompt: &survey.Input{
				Message: "Email Address",
				Help:    "The email address of the cluster administrator. This will be used to log in to the console.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return validate.Email(ans.(string))
			}),
		},
		"OPENSHIFT_INSTALL_EMAIL_ADDRESS",
	)
	a.EmailAddress = email
	return err
}

// Name returns the human-friendly name of the asset.
func (a *emailAddress) Name() string {
	return "Email Address"
}
