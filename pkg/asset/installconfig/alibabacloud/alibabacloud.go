package alibabacloud

import (
	"fmt"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/openshift/installer/pkg/types/alibabacloud"
	"github.com/pkg/errors"
)

const (
	DefaultRegion         = "cn-hangzhou"
	DefaultAcceptLanguage = "en-US"
)

const (
	noResourceGroup = "<none>"
)

func Platform() (*alibabacloud.Platform, error) {
	client, err := NewClient(DefaultRegion)
	if err != nil {
		return nil, err
	}

	region, err := selectRegion(client)
	if err != nil {
		return nil, err
	}

	resourceGroup, err := selectResourceGroup(client)
	if err != nil {
		return nil, err
	}

	return &alibabacloud.Platform{
		Region:            region,
		ResourceGroupName: resourceGroup,
	}, nil
}

func selectRegion(client *Client) (string, error) {
	regions_resp, err := client.DescribeRegions("cn-hangzhou")
	if err != nil {
		return "", err
	}
	regions := regions_resp.Regions.Region

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
				Default: DefaultRegion,
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
	groups_resp, err := client.ListResourceGroups()
	if err != nil {
		return "", errors.Wrap(err, "failed to list resource groups")
	}

	groups := groups_resp.ResourceGroups.ResourceGroup

	var options []string
	names := make(map[string]string)

	for _, group := range groups {
		option := fmt.Sprintf("%s (%s)", group.Name, group.Id)
		names[option] = group.Name
		options = append(options, option)
	}
	sort.Strings(options)

	options = append(options, noResourceGroup)
	names[noResourceGroup] = ""

	var selectedResourceGroup string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Resource Group",
				Help:    "The resource group where the cluster will be provisioned.",
				Options: options,
				Default: noResourceGroup,
			},
		},
	}, &selectedResourceGroup)
	return names[selectedResourceGroup], err
}
