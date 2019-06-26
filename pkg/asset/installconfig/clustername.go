package installconfig

import (
	"github.com/sirupsen/logrus"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types/validation"
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
	}
}

// Generate queries for the cluster name from the user.
func (a *clusterName) Generate(log *logrus.Entry, parents asset.Parents) error {
	bd := &baseDomain{}
	parents.Get(bd)

	return survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Cluster Name",
				Help:    "The name of the cluster.  This will be used when generating sub-domains.\n\nFor libvirt, choose a name that is unique enough to be used as a prefix during cluster deletion.  For example, if you use 'demo' as your cluster name, `openshift-install destroy cluster` may destroy all domains, networks, pools, and volumes that begin with 'demo'.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return validate.DomainName(validation.ClusterDomain(bd.BaseDomain, ans.(string)), false)
			}),
		},
	}, &a.ClusterName)
}

// Name returns the human-friendly name of the asset.
func (a *clusterName) Name() string {
	return "Cluster Name"
}
