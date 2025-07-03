package aws

import (
	"fmt"
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/version"
)

// Platform collects AWS-specific configuration.
func Platform() (*aws.Platform, error) {
	architecture := version.DefaultArch()
	regions, err := knownPublicRegions(architecture)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS public regions: %w", err)
	}
	longRegions := make([]string, 0, len(aws.RegionLookupMap))
	shortRegions := make([]string, 0, len(aws.RegionLookupMap))
	for _, region := range regions {
		if longName, ok := aws.RegionLookupMap[region]; ok {
			longRegions = append(longRegions, fmt.Sprintf("%s (%s)", region, longName))
			shortRegions = append(shortRegions, region)
		}
	}

	var regionTransform survey.Transformer = func(ans interface{}) interface{} {
		switch v := ans.(type) {
		case core.OptionAnswer:
			return core.OptionAnswer{Value: strings.SplitN(v.Value, " ", 2)[0], Index: v.Index}
		case string:
			return strings.SplitN(v, " ", 2)[0]
		}
		return ""
	}

	defaultRegion := "us-east-1"
	if found, err := IsKnownPublicRegion(defaultRegion, architecture); !found || err != nil {
		panic(fmt.Sprintf("installer bug: invalid default AWS region %q", defaultRegion))
	}

	ssn, err := GetSession()
	if err != nil {
		return nil, err
	}

	defaultRegionPointer := ssn.Config.Region
	if defaultRegionPointer != nil && *defaultRegionPointer != "" {
		found, err := IsKnownPublicRegion(*defaultRegionPointer, architecture)
		if err != nil {
			return nil, fmt.Errorf("failed to determine if region is public: %w", err)
		}
		if found {
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
				Default: fmt.Sprintf("%s (%s)", defaultRegion, aws.RegionLookupMap[defaultRegion]),
				Options: longRegions,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := regionTransform(ans).(core.OptionAnswer).Value
				i := sort.SearchStrings(shortRegions, choice)
				if i == len(shortRegions) || shortRegions[i] != choice {
					return fmt.Errorf("invalid region %q", choice)
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
