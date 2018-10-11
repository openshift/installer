package installconfig

import (
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
)

type password struct {
	Password string
}

var _ asset.Asset = (*password)(nil)

// Dependencies returns no dependencies.
func (a *password) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate queries for the password from the user.
func (a *password) Generate(asset.Parents) error {
	p, err := asset.GenerateUserProvidedAsset(
		a.Name(),
		&survey.Question{
			Prompt: &survey.Password{
				Message: "Password",
				Help:    "The password of the cluster administrator. This will be used to log in to the console.",
			},
		},
		"OPENSHIFT_INSTALL_PASSWORD",
	)
	a.Password = p
	return err
}

// Name returns the human-friendly name of the asset.
func (a *password) Name() string {
	return "Password"
}
