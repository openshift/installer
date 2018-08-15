package installerassets

import (
	"context"
	"os"

	"github.com/openshift/installer/pkg/validate"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func getClusterName(ctx context.Context) ([]byte, error) {
	value := os.Getenv("OPENSHIFT_INSTALL_CLUSTER_NAME")
	if value != "" {
		err := validate.DomainName(value)
		if err != nil {
			return nil, err
		}
		return []byte(value), nil
	}

	question := &survey.Question{
		Prompt: &survey.Input{
			Message: "Cluster Name",
			Help:    "The name of the cluster. This will be used when generating sub-domains.",
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
	Defaults["cluster-name"] = getClusterName
}
