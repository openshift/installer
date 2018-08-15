package openstack

import (
	"context"
	"os"

	"github.com/openshift/installer/pkg/installerassets"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func getNetwork(ctx context.Context) (data []byte, err error) {
	value := os.Getenv("OPENSHIFT_INSTALL_OPENSTACK_EXTERNAL_NETWORK")
	if value != "" {
		//FIXME(shardy) add some validation here
		return []byte(value), nil
	}

	question := &survey.Question{
		Prompt: &survey.Input{
			Message: "ExternalNetwork",
			Help:    "The OpenStack external network to be used for installation.",
		},
		Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
			//FIXME(shadower) add some validation here
			return nil
		}),
	}

	var response string
	err = survey.Ask([]*survey.Question{question}, &response)
	if err != nil {
		return nil, errors.Wrap(err, "ask")
	}

	return []byte(response), nil
}

func init() {
	installerassets.Defaults["openstack/external-network"] = getNetwork
}
