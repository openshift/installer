package installconfig

import (
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
)

const (
	noPasswordHash = ""
)

type passwordHash struct {
	Hash string
}

var _ asset.Asset = (*passwordHash)(nil)

// Dependencies returns no dependencies.
func (a *passwordHash) Dependencies() []asset.Asset {
	return nil
}

// Generate generates the SSH public key asset.
func (a *passwordHash) Generate(asset.Parents) error {
	var hash string

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Password Hash",
				Help:    "The password hash is used to access the bootstrap node from a console. This is optional.",
				Default: noPasswordHash,
			},
		},
	}, &hash); err != nil {
		return errors.Wrap(err, "failed UserInput for Password Hash")
	}
	a.Hash = hash
	return nil
}

// Name returns the human-friendly name of the asset.
func (a passwordHash) Name() string {
	return "Password Hash"
}
