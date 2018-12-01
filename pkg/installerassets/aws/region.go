package aws

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/openshift/installer/pkg/installerassets"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

var (
	validAWSRegions = map[string]string{
		"ap-northeast-1": "Tokyo",
		"ap-northeast-2": "Seoul",
		"ap-northeast-3": "Osaka-Local",
		"ap-south-1":     "Mumbai",
		"ap-southeast-1": "Singapore",
		"ap-southeast-2": "Sydney",
		"ca-central-1":   "Central",
		"cn-north-1":     "Beijing",
		"cn-northwest-1": "Ningxia",
		"eu-central-1":   "Frankfurt",
		"eu-west-1":      "Ireland",
		"eu-west-2":      "London",
		"eu-west-3":      "Paris",
		"sa-east-1":      "São Paulo",
		"us-east-1":      "N. Virginia",
		"us-east-2":      "Ohio",
		"us-west-1":      "N. California",
		"us-west-2":      "Oregon",
	}
)

func getRegion(ctx context.Context) ([]byte, error) {
	longRegions := make([]string, 0, len(validAWSRegions))
	shortRegions := make([]string, 0, len(validAWSRegions))
	for id, location := range validAWSRegions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
		shortRegions = append(shortRegions, id)
	}
	regionTransform := survey.TransformString(func(s string) string {
		return strings.SplitN(s, " ", 2)[0]
	})
	sort.Strings(longRegions)
	sort.Strings(shortRegions)

	value := os.Getenv("OPENSHIFT_INSTALL_AWS_REGION")
	if value != "" {
		i := sort.SearchStrings(shortRegions, value)
		if i == len(shortRegions) || shortRegions[i] != value {
			return nil, errors.Errorf("invalid region %q", value)
		}
		return []byte(value), nil
	}

	question := &survey.Question{
		Prompt: &survey.Select{
			Message: "Region",
			Help:    "The AWS region to be used for installation.",
			Default: "us-east-1 (N. Virginia)",
			Options: longRegions,
		},
		Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
			choice := regionTransform(ans).(string)
			i := sort.SearchStrings(shortRegions, choice)
			if i == len(shortRegions) || shortRegions[i] != choice {
				return errors.Errorf("invalid region %q", choice)
			}
			return nil
		}),
		Transform: regionTransform,
	}

	var response string
	err := survey.Ask([]*survey.Question{question}, &response)
	if err != nil {
		return nil, errors.Wrap(err, "ask")
	}

	return []byte(response), nil
}

func init() {
	installerassets.Defaults["aws/region"] = getRegion
}
