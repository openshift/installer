package installconfig

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	azureconfig "github.com/openshift/installer/pkg/asset/installconfig/azure"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/validate"
)

type baseDomain struct {
	BaseDomain string
}

var _ asset.Asset = (*baseDomain)(nil)

// Dependencies returns no dependencies.
func (a *baseDomain) Dependencies() []asset.Asset {
	return []asset.Asset{
		&platform{},
	}
}

// Generate queries for the base domain from the user.
func (a *baseDomain) Generate(parents asset.Parents) error {
	platform := &platform{}
	parents.Get(platform)

	switch platform.CurrentName() {
	case aws.Name:
		var err error
		zone, err := awsconfig.GetBaseDomain()
		cause := errors.Cause(err)
		a.BaseDomain = zone.Name
		platform.AWS.SetBaseDomain(zone.ID)
		if !(awsconfig.IsForbidden(cause) || request.IsErrorThrottle(cause)) {
			return err
		}
	case azure.Name:
		var err error
		azureDNS, _ := azureconfig.NewDNSConfig()
		zone, err := azureDNS.GetDNSZone()
		if err != nil {
			return err
		}
		a.BaseDomain = zone.Name
		return platform.Azure.SetBaseDomain(zone.ID)
	case gcp.Name:
		var err error
		a.BaseDomain, err = gcpconfig.GetBaseDomain(platform.GCP.ProjectID)

		// We are done if success (err == nil) or an err besides forbidden/throttling
		if !(gcpconfig.IsForbidden(err) || gcpconfig.IsThrottled(err)) {
			return err
		}
	default:
		//Do nothing
	}

	return survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Base Domain",
				Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base and will also include the cluster name.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return validate.DomainName(ans.(string), true)
			}),
		},
	}, &a.BaseDomain)
}

// Name returns the human-friendly name of the asset.
func (a *baseDomain) Name() string {
	return "Base Domain"
}
