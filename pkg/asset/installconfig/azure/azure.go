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

	azsub "github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/subscriptions"
)


// Platform collects azure-specific configuration.
func Platform() (*azure.Platform, error) {
	regions, err := getRegions()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get list of regions")
	}
	longRegions := make([]string, 0, len(regions))
	shortRegions := make([]string, 0, len(regions))
	for id, location := range regions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
		shortRegions = append(shortRegions, id)
	}
	regionTransform := survey.TransformString(func(s string) string {
		return strings.SplitN(s, " ", 2)[0]
	})

	sort.Strings(longRegions)
	sort.Strings(shortRegions)

	var region string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The azure region to be used for installation.",				
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
	session, err := GetSession()
	if err != nil {
		return nil, err
	}
	client := azsub.NewClientWithBaseURI(session.Credentials.ResourceManagerEndpoint)
	client.Authorizer = session.Authorizer
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	locations, err := client.ListLocations(ctx, session.Credentials.SubscriptionID)
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
