package alibabacloud

import (
	"fmt"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types/alibabacloud"
)

const (
	defaultRegion         = "cn-hangzhou"
	defaultAcceptLanguage = "en-US"
)

// Platform collects AlibabaCloud-specific configuration.
func Platform() (*alibabacloud.Platform, error) {
	client, err := NewClient(defaultRegion)
	if err != nil {
		return nil, err
	}

	region, err := selectRegion(client)
	if err != nil {
		return nil, err
	}

	client, err = NewClient(region)
	if err != nil {
		return nil, err
	}

	resourceGroup, err := selectResourceGroup(client)
	if err != nil {
		return nil, err
	}

	return &alibabacloud.Platform{
		Region:          region,
		ResourceGroupID: resourceGroup,
	}, nil
}

func selectRegion(client *Client) (string, error) {
	regionsResponse, err := client.DescribeRegions()
	if err != nil {
		return "", err
	}
	regions := regionsResponse.Regions.Region

	longRegions := make([]string, 0, len(regions))
	shortRegions := make([]string, 0, len(regions))
	for _, location := range regions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", location.RegionId, location.LocalName))
		shortRegions = append(shortRegions, location.RegionId)
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

	var selectedRegion string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The Alibaba Cloud region to be used for installation.",
				Default: defaultRegion,
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

func selectResourceGroup(client *Client) (string, error) {
	groupsResponse, err := client.ListResourceGroups()
	if err != nil {
		return "", errors.Wrap(err, "failed to list resource groups")
	}

	groups := groupsResponse.ResourceGroups.ResourceGroup

	if len(groups) == 0 {
		return "", errors.Wrap(err, "resource group not found")
	}

	var options []string
	names := make(map[string]string)

	for _, group := range groups {
		option := fmt.Sprintf("%s (%s)", group.Name, group.Id)
		names[option] = group.Id
		options = append(options, option)
	}
	sort.Strings(options)

	var selectedResourceGroup string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Resource Group",
				Help:    "The resource group where the cluster will be provisioned.",
				Options: options,
			},
		},
	}, &selectedResourceGroup)
	return names[selectedResourceGroup], err
}
