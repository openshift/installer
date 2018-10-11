package installconfig

import (
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/validate"
)

type pullSecret struct {
	PullSecret string
}

var _ asset.Asset = (*pullSecret)(nil)

// Dependencies returns no dependencies.
func (a *pullSecret) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate queries for the pull secret from the user.
func (a *pullSecret) Generate(asset.Parents) error {
	s, err := asset.GenerateUserProvidedAssetForPath(
		a.Name(),
		&survey.Question{
			Prompt: &survey.Input{
				Message: "Pull Secret",
				Help:    "The container registry pull secret for this cluster.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return validate.JSON([]byte(ans.(string)))
			}),
		},
		"OPENSHIFT_INSTALL_PULL_SECRET",
		"OPENSHIFT_INSTALL_PULL_SECRET_PATH",
	)
	a.PullSecret = s
	return err
}

// Name returns the human-friendly name of the asset.
func (a *pullSecret) Name() string {
	return "Pull Secret"
}
