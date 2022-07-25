package installconfig

import (
	survey "github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	alibabacloudconfig "github.com/openshift/installer/pkg/asset/installconfig/alibabacloud"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	azureconfig "github.com/openshift/installer/pkg/asset/installconfig/azure"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	ibmcloudconfig "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	powervsconfig "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/types/alibabacloud"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/powervs"
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

	var err error
	switch platform.CurrentName() {
	case alibabacloud.Name:
		a.BaseDomain, err = alibabacloudconfig.GetBaseDomain()
		if err != nil {
			return err
		}
		return nil
	case aws.Name:
		a.BaseDomain, err = awsconfig.GetBaseDomain()
		cause := errors.Cause(err)
		if !(awsconfig.IsForbidden(cause) || request.IsErrorThrottle(cause)) {
			return err
		}
	case azure.Name:
		// Create client using public cloud because install config has not been generated yet.
		ssn, err := azureconfig.GetSession(azure.PublicCloud, "")
		if err != nil {
			return err
		}
		azureDNS := azureconfig.NewDNSConfig(ssn)
		zone, err := azureDNS.GetDNSZone()
		if err != nil {
			return err
		}
		a.BaseDomain = zone.Name
		return platform.Azure.SetBaseDomain(zone.ID)
	case gcp.Name:
		a.BaseDomain, err = gcpconfig.GetBaseDomain(platform.GCP.ProjectID)

		// We are done if success (err == nil) or an err besides forbidden/throttling
		if !(gcpconfig.IsForbidden(err) || gcpconfig.IsThrottled(err)) {
			return err
		}
	case ibmcloud.Name:
		zone, err := ibmcloudconfig.GetDNSZone()
		if err != nil {
			return err
		}
		a.BaseDomain = zone.Name
		return nil
	case powervs.Name:
		zone, err := powervsconfig.GetDNSZone()
		if err != nil {
			return err
		}
		a.BaseDomain = zone.Name
		return nil
	default:
		//Do nothing
	}

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Base Domain",
				Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base and will also include the cluster name.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return validate.DomainName(ans.(string), true)
			}),
		},
	}, &a.BaseDomain); err != nil {
		return errors.Wrap(err, "failed UserInput")
	}
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *baseDomain) Name() string {
	return "Base Domain"
}
