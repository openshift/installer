package ibmcloud

import (
	"context"
	"fmt"
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/ibmcloud/validation"
	"github.com/pkg/errors"
)

const (
	createNew = "<create new>"
)

// Platform collects IBM Cloud-specific configuration.
func Platform() (*ibmcloud.Platform, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	resourceGroup, err := selectResourceGroup(context.TODO(), client)
	if err != nil {
		return nil, err
	}

	region, err := selectRegion()
	if err != nil {
		return nil, err
	}

	return &ibmcloud.Platform{
		ResourceGroupName: resourceGroup,
		Region:            region,
	}, nil
}

func selectResourceGroup(ctx context.Context, client *Client) (string, error) {
	groups, err := client.GetResourceGroups(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to list resource groups")
	}

	var options []string
	names := make(map[string]string)

	for _, group := range groups {
		option := fmt.Sprintf("%s (%s)", *group.Name, *group.ID)
		names[option] = *group.Name
		options = append(options, option)
	}
	sort.Strings(options)

	options = append(options, createNew)
	names[createNew] = ""

	var selectedResourceGroup string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Resource Group",
				Help:    "The resource group where the cluster will be provisioned.",
				Options: options,
				Default: createNew,
			},
		},
	}, &selectedResourceGroup)
	return names[selectedResourceGroup], err
}

func selectRegion() (string, error) {
	longRegions := make([]string, 0, len(validation.Regions))
	shortRegions := make([]string, 0, len(validation.Regions))
	for id, location := range validation.Regions {
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
				Default: fmt.Sprintf("%s (%s)", defaultRegion, validation.Regions[defaultRegion]),
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
