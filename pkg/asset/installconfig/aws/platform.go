package aws

import (
	"fmt"
	"sort"

	survey "github.com/AlecAivazis/survey/v2"
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

	sort.Strings(regions)

	var region string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The AWS region to be used for installation.",
				Default: defaultRegion,
				Options: regions,
			},
		},
	}, &region)
	if err != nil {
		return nil, err
	}

	return &aws.Platform{
		Region: region,
	}, nil
}
