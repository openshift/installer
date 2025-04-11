package installconfig

import (
	"context"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/validate"
)

type clusterName struct {
	ClusterName string
}

var _ asset.Asset = (*clusterName)(nil)

// Dependencies returns no dependencies.
func (a *clusterName) Dependencies() []asset.Asset {
	return []asset.Asset{
		&baseDomain{},
		&platform{},
	}
}

// Generate queries for the cluster name from the user.
func (a *clusterName) Generate(_ context.Context, parents asset.Parents) error {
	bd := &baseDomain{}
	platform := &platform{}
	parents.Get(bd, platform)

	validator := survey.Required

	if platform.GCP != nil || platform.Azure != nil {
		validator = survey.ComposeValidators(validator, func(ans interface{}) error {
			return validate.ClusterName1035(ans.(string))
		})
		if platform.GCP != nil {
			validator = survey.ComposeValidators(validator, func(ans interface{}) error {
				return validate.GCPClusterName(ans.(string))
			})
		}
	}

	if platform.Ovirt != nil {
		validator = survey.ComposeValidators(validator, func(ans interface{}) error {
			return validate.ClusterName(ans.(string))
		})
	}
	if platform.VSphere != nil || platform.BareMetal != nil || platform.Nutanix != nil {
		validator = survey.ComposeValidators(validator, func(ans interface{}) error {
			return validate.OnPremClusterName(ans.(string))
		})
	}
	validator = survey.ComposeValidators(validator, func(ans interface{}) error {
		installConfig := &types.InstallConfig{BaseDomain: bd.BaseDomain}
		installConfig.ObjectMeta.Name = ans.(string)
		return validate.DomainName(installConfig.ClusterDomain(), false)
	})

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Cluster Name",
				Help:    "The name of the cluster.  This will be used when generating sub-domains.\n\nFor libvirt, choose a name that is unique enough to be used as a prefix during cluster deletion.  For example, if you use 'demo' as your cluster name, `openshift-install destroy cluster` may destroy all domains, networks, pools, and volumes that begin with 'demo'.",
			},
			Validate: validator,
		},
	}, &a.ClusterName); err != nil {
		return errors.Wrap(err, "failed UserInput")
	}
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *clusterName) Name() string {
	return "Cluster Name"
}
