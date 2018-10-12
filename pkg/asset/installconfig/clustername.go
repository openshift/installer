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
				Help:    "The name of the cluster. This will be used when generating sub-domains.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return validate.DomainName(ans.(string))
			}),
		},
		"OPENSHIFT_INSTALL_CLUSTER_NAME",
	)
	a.ClusterName = n
	return err
}

// Name returns the human-friendly name of the asset.
func (a *clusterName) Name() string {
	return "Cluster Name"
}
