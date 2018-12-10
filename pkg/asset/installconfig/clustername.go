package installconfig

import (
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/validate"
)

type clusterName struct {
	ClusterName string
}

var _ asset.Asset = (*clusterName)(nil)

// Dependencies returns no dependencies.
func (a *clusterName) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate queries for the cluster name from the user.
func (a *clusterName) Generate(asset.Parents) error {
	n, err := asset.GenerateUserProvidedAsset(
		a.Name(),
		&survey.Question{
			Prompt: &survey.Input{
				Message: "Cluster Name",
				Help:    "The name of the cluster.  This will be used when generating sub-domains.\n\nFor libvirt, choose a name that is unique enough to be used as a prefix during cluster deletion.  For example, if you use 'demo' as your cluster name, `openshift-install destroy cluster` may destroy all domains, networks, pools, and volumes that begin with 'demo'.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return validate.DomainName(ans.(string))
			}),
		},
	)
	a.ClusterName = n
	return err
}

// Name returns the human-friendly name of the asset.
func (a *clusterName) Name() string {
	return "Cluster Name"
}
