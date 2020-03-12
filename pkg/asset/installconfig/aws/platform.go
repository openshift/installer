package aws

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/types/aws"
)

// Platform collects AWS-specific configuration.
func Platform() (*aws.Platform, error) {
	regions := knownRegions()
	longRegions := make([]string, 0, len(regions))
	shortRegions := make([]string, 0, len(regions))
	for id, location := range regions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
		shortRegions = append(shortRegions, id)
	}
	regionTransform := survey.TransformString(func(s string) string {
		return strings.SplitN(s, " ", 2)[0]
	})

	defaultRegion := "us-east-1"
	if !IsKnownRegion(defaultRegion) {
		panic(fmt.Sprintf("installer bug: invalid default AWS region %q", defaultRegion))
	}

	ssn, err := GetSession()
	if err != nil {
		return nil, err
	}

	defaultRegionPointer := ssn.Config.Region
	if defaultRegionPointer != nil && *defaultRegionPointer != "" {
		if IsKnownRegion(*defaultRegionPointer) {
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
				Default: fmt.Sprintf("%s (%s)", defaultRegion, regions[defaultRegion]),
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
