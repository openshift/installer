package azure

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/to"

	"github.com/openshift/installer/pkg/types/azure"

	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

const (
	defaultRegion string = "eastus"
)

// Platform collects azure-specific configuration.
func Platform() (*azure.Platform, error) {
	regions, err := getRegions()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get list of regions")
	}

	resourceCapableRegions, err := getResourceCapableRegions()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get list of resources to check available regions")
	}

	longRegions := make([]string, 0, len(regions))
	shortRegions := make([]string, 0, len(regions))
	for id, location := range regions {
		for _, resourceCapableRegion := range resourceCapableRegions {
			// filter our regions not capable of having resources created (we check for resource groups)
			if resourceCapableRegion == location {
				longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
				shortRegions = append(shortRegions, id)
			}
		}
	}

	regionTransform := survey.TransformString(func(s string) string {
		return strings.SplitN(s, " ", 2)[0]
	})

	_, ok := regions[defaultRegion]
	if !ok {
		return nil, errors.Errorf("installer bug: invalid default azure region %q", defaultRegion)
	}

	sort.Strings(longRegions)
	sort.Strings(shortRegions)

	var region string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The azure region to be used for installation.",
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

	return &azure.Platform{
		Region: region,
	}, nil
}

func getRegions() (map[string]string, error) {
	client, err := NewClient(context.TODO())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()

	locations, err := client.ListLocations(ctx)
	if err != nil {
		return nil, err
	}

	allLocations := map[string]string{}
	for _, location := range *locations {
		allLocations[to.String(location.Name)] = to.String(location.DisplayName)
	}
	return allLocations, nil
}

func getResourceCapableRegions() ([]string, error) {
	client, err := NewClient(context.TODO())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()

	provider, err := client.GetResourcesProvider(ctx, "Microsoft.Resources")
	if err != nil {
		return nil, err
	}

	for _, resType := range *provider.ResourceTypes {
		if *resType.ResourceType == "resourceGroups" {
			return *resType.Locations, nil
		}
	}

	return []string{}, nil
}
