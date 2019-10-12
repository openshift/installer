package aws

import (
	"fmt"
	"sort"
	"strings"

	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/aws/validation"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Platform collects AWS-specific configuration.
func Platform() (*aws.Platform, error) {
	longRegions := make([]string, 0, len(validation.Regions))
	shortRegions := make([]string, 0, len(validation.Regions))
	for id, location := range validation.Regions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
		shortRegions = append(shortRegions, id)
	}
	regionTransform := survey.TransformString(func(s string) string {
		return strings.SplitN(s, " ", 2)[0]
	})

	defaultRegion := "us-east-1"
	_, ok := validation.Regions[defaultRegion]
	if !ok {
		panic(fmt.Sprintf("installer bug: invalid default AWS region %q", defaultRegion))
	}

	ssn, err := GetSession()
	if err != nil {
		return nil, err
	}

	defaultRegionPointer := ssn.Config.Region
	if defaultRegionPointer != nil && *defaultRegionPointer != "" {
		_, ok := validation.Regions[*defaultRegionPointer]
		if ok {
			defaultRegion = *defaultRegionPointer
		} else {
			logrus.Warnf("Unrecognized AWS region %q, defaulting to %s", *defaultRegionPointer, defaultRegion)
		}
	}

	sort.Strings(longRegions)
	sort.Strings(shortRegions)

	var region string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The AWS region to be used for installation.",
				Default: fmt.Sprintf("%s (%s)", defaultRegion, validation.Regions[defaultRegion]),
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
		},
	}, &region)
	if err != nil {
		return nil, err
	}

	return &aws.Platform{
		Region: region,
	}, nil
}
