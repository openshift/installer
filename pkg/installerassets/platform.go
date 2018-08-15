package installerassets

import (
	"context"
	"os"
	"sort"

	"github.com/openshift/installer/pkg/types"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func getPlatform(ctx context.Context) (data []byte, err error) {
	value := os.Getenv("OPENSHIFT_INSTALL_PLATFORM")
	if value != "" {
		i := sort.SearchStrings(types.PlatformNames, value)
		if i == len(types.PlatformNames) || types.PlatformNames[i] != value {
			return nil, errors.Errorf("invalid platform %q", value)
		}
		return []byte(value), nil
	}

	question := &survey.Question{
		Prompt: &survey.Select{
			Message: "Platform",
			Options: types.PlatformNames,
		},
		Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
			choice := ans.(string)
			i := sort.SearchStrings(types.PlatformNames, choice)
			if i == len(types.PlatformNames) || types.PlatformNames[i] != choice {
				return errors.Errorf("invalid platform %q", choice)
			}
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
	Defaults["platform"] = getPlatform
}
