package openstack

import (
	"context"
	"os"

	"github.com/openshift/installer/pkg/installerassets"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func getCloud(ctx context.Context) ([]byte, error) {
	value := os.Getenv("OPENSHIFT_INSTALL_OPENSTACK_CLOUD")
	if value != "" {
		//FIXME(russellb) add some validation here
		return []byte(value), nil
	}

	question := &survey.Question{
		//TODO(russellb) - We could open clouds.yaml here and read the list of defined clouds
		//and then use survey.Select to let the user choose one.
		Prompt: &survey.Input{
			Message: "Cloud",
			Help:    "The OpenStack cloud name from clouds.yaml.",
			Default: "cloudOne",
		},
		Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
			//FIXME(russellb) add some validation here
			return nil
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
	installerassets.Defaults["openstack/cloud"] = getCloud
}
