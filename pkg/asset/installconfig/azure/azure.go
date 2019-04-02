package azure

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/Azure/go-autorest/autorest/to"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/azure/validation"

	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"

	azsub "github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/subscriptions"
)

const (
	defaultRegion string = "eastus"
)

// Platform collects azure-specific configuration.
func Platform() (*azure.Platform, error) {
	regions, err := getRegions()
	logrus.Debugf("##### Retrieved regions for given subscription. count : %v", len(regions))
	longRegions := make([]string, 0, len(validation.Regions))
	shortRegions := make([]string, 0, len(validation.Regions))
	for id, location := range regions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
		shortRegions = append(shortRegions, id)
	}
	regionTransform := survey.TransformString(func(s string) string {
		return strings.SplitN(s, " ", 2)[0]
	})

	_, ok := validation.Regions[defaultRegion]
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

	return &azure.Platform{
		Region: region,
	}, nil
}

func getRegions() (map[string]string, error) {
	ctx := context.TODO()
	session, err := GetSession()
	if err != nil {
		return nil, err
	}
	client := azsub.NewClient()
	client.Authorizer = session.Authorizer
	locations, err := client.ListLocations(ctx, session.SubscriptionID)
	if err != nil {
		return nil, err
	}

	locationsValue := *locations.Value

	allLocations := map[string]string{}
	for _, location := range locationsValue {
		allLocations[to.String(location.Name)] = to.String(location.DisplayName)
	}
	return allLocations, nil
}
