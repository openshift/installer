package installerassets

import (
	"context"
	"os"

	"github.com/openshift/installer/pkg/validate"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func getAdminEmail(ctx context.Context) ([]byte, error) {
	value := os.Getenv("OPENSHIFT_INSTALL_EMAIL_ADDRESS")
	if value != "" {
		err := validate.Email(value)
		if err != nil {
			return nil, err
		}
		return []byte(value), nil
	}

	question := &survey.Question{
		Prompt: &survey.Input{
			Message: "Email Address",
			Help:    "The email address of the cluster administrator. This will be used to log in to the console.",
		},
		Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
			return validate.Email(ans.(string))
		}),
	}

	var response string
	err := survey.Ask([]*survey.Question{question}, &response)
	if err != nil {
		return nil, errors.Wrap(err, "ask")
	}

	return []byte(response), nil
}

func getAdminPassword(ctx context.Context) ([]byte, error) {
	value := os.Getenv("OPENSHIFT_INSTALL_PASSWORD")
	if value != "" {
		return []byte(value), nil
	}

	question := &survey.Question{
		Prompt: &survey.Password{
			Message: "Password",
			Help:    "The password of the cluster administrator. This will be used to log in to the console.",
		},
	}

	var response string
	err := survey.Ask([]*survey.Question{question}, &response)
	if err != nil {
		return nil, errors.Wrap(err, "ask")
	}

	return []byte(response), nil
}

func init() {
	Defaults["admin/email"] = getAdminEmail
	Defaults["admin/password"] = getAdminPassword
}
