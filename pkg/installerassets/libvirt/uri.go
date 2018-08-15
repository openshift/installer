package libvirt

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/openshift/installer/pkg/installerassets"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func getURI(ctx context.Context) ([]byte, error) {
	value := os.Getenv("OPENSHIFT_INSTALL_LIBVIRT_URI")
	if value != "" {
		err := validURI(value)
		if err != nil {
			return nil, errors.Wrap(err, "resolve OPENSHIFT_INSTALL_LIBVIRT_URI")
		}
		return []byte(value), nil
	}

	question := &survey.Question{
		Prompt: &survey.Input{
			Message: "Libvirt Connection URI",
			Help:    "The libvirt connection URI to be used. This must be accessible from the running cluster.",
			Default: "qemu+tcp://192.168.122.1/system",
		},
		Validate: survey.ComposeValidators(survey.Required, uriValidator),
	}

	var response string
	err := survey.Ask([]*survey.Question{question}, &response)
	if err != nil {
		return nil, errors.Wrap(err, "ask")
	}

	return []byte(response), nil
}

// uriValidator validates if the answer provided in prompt is a valid
// url and has non-empty scheme.
func uriValidator(ans interface{}) error {
	return validURI(ans.(string))
}

// validURI validates if the URI is a valid URI with a non-empty scheme.
func validURI(uri string) error {
	parsed, err := url.Parse(uri)
	if err != nil {
		return err
	}
	if parsed.Scheme == "" {
		return fmt.Errorf("invalid URI %q (no scheme)", uri)
	}
	return nil
}

func init() {
	installerassets.Defaults["libvirt/uri"] = getURI
}
