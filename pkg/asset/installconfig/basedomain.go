package installconfig

import (
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/validate"
)

type baseDomain struct {
	BaseDomain string
}

var _ asset.Asset = (*baseDomain)(nil)

// Dependencies returns no dependencies.
func (a *baseDomain) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate queries for the base domain from the user.
func (a *baseDomain) Generate(asset.Parents) error {
	bd, err := asset.GenerateUserProvidedAsset(
		a.Name(),
		&survey.Question{
			Prompt: &survey.Input{
				Message: "Base Domain",
				Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base.\n\nFor AWS, this must be a previously-existing public Route 53 zone.  You can check for any already in your account with:\n\n  $ aws route53 list-hosted-zones --query 'HostedZones[? !(Config.PrivateZone)].Name' --output text",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return validate.DomainName(ans.(string))
			}),
		},
		"OPENSHIFT_INSTALL_BASE_DOMAIN",
	)
	a.BaseDomain = bd
	return err
}

// Name returns the human-friendly name of the asset.
func (a *baseDomain) Name() string {
	return "Base Domain"
}
