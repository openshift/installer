package ibmcloud

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/pkg/errors"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types/ibmcloud"
)

const (
	createNew = "<create new>"
)

// Platform collects IBM Cloud-specific configuration.
func Platform() (*ibmcloud.Platform, error) {
	// Since installconfig has not be created yet, attempt to setup a new IBM Cloud Client.
	// Default will rely on Public IBM Cloud Service Endpoints, but for Region collection, we accept an environment variable to override the IBM Cloud Global Catalog Service endpoint, as the list of endpoint overrides would be specified within the installconfig.
	endpoints := make([]configv1.IBMCloudServiceEndpoint, 0)
	if gcEndpoint := os.Getenv(ibmcloud.IBMCloudServiceGlobalCatalogVar); gcEndpoint != "" {
		endpoints = append(endpoints, configv1.IBMCloudServiceEndpoint{
			Name: configv1.IBMCloudServiceGlobalCatalog,
			URL:  gcEndpoint,
		})
	}

	client, err := NewClient(endpoints)
	if err != nil {
		return nil, fmt.Errorf("failed creating ibmcloud client for region retrieval: %w", err)
	}
	regions, err := client.GetIBMCloudRegions(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve ibmcloud regions: %w", err)
	}

	region, err := selectRegion(regions)
	if err != nil {
		return nil, fmt.Errorf("failed to survey desired ibmcloud region: %w", err)
	}

	return &ibmcloud.Platform{
		Region: region,
	}, nil
}

func selectRegion(regions map[string]string) (string, error) {
	longRegions := make([]string, 0, len(regions))
	shortRegions := make([]string, 0, len(regions))
	for id, location := range regions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
		shortRegions = append(shortRegions, id)
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

	sort.Strings(longRegions)
	sort.Strings(shortRegions)

	defaultRegion := longRegions[0]

	var selectedRegion string
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The IBM Cloud region to be used for installation.",
				Default: fmt.Sprintf("%s (%s)", defaultRegion, regions[defaultRegion]),
				Options: longRegions,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := regionTransform(ans).(core.OptionAnswer).Value
				i := sort.SearchStrings(shortRegions, choice)
				if i == len(shortRegions) || shortRegions[i] != choice {
					return errors.Errorf("invalid region %q", choice)
				}
				return nil
			}),
			Transform: regionTransform,
		},
	}, &selectedRegion)
	if err != nil {
		return "", err
	}
	return selectedRegion, nil
}
