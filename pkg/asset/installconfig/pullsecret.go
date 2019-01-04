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
	return survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "Pull Secret",
				Help:    "The container registry pull secret for this cluster, as a single line of JSON (e.g. {\"auths\": {...}}).\n\nYou can get this secret from https://cloud.openshift.com/clusters/install#pull-secret",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return validate.ImagePullSecret(ans.(string))
			}),
		},
	}, &a.PullSecret)
}

// Name returns the human-friendly name of the asset.
func (a *pullSecret) Name() string {
	return "Pull Secret"
}
