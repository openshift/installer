package installerassets

import (
	"context"
	"os"

	"github.com/openshift/installer/pkg/validate"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func getBaseDomain(ctx context.Context) ([]byte, error) {
	value := os.Getenv("OPENSHIFT_INSTALL_BASE_DOMAIN")
	if value != "" {
		err := validate.DomainName(value)
		if err != nil {
			return nil, err
		}
		return []byte(value), nil
	}

	question := &survey.Question{
		Prompt: &survey.Input{
			Message: "Base Domain",
			Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base.\n\nFor AWS, this must be a previously-existing public Route 53 zone.  You can check for any already in your account with:\n\n  $ aws route53 list-hosted-zones --query 'HostedZones[? !(Config.PrivateZone)].Name' --output text",
		},
		Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
			return validate.DomainName(ans.(string))
		}),
	}

	var response string
	err := survey.Ask([]*survey.Question{question}, &response)
	if err != nil {
		return nil, errors.Wrap(err, "ask")
	}

	return []byte(response), nil
}

func init() {
	Defaults["base-domain"] = getBaseDomain
}
