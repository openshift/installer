package alibabacloud

import (
	"fmt"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/types/alibabacloud"
)

const (
	defaultRegion         = "cn-hangzhou"
	defaultAcceptLanguage = "en-US"
)

var supportedRegions = sets.NewString(
	"cn-qingdao",
	"cn-beijing",
	"cn-zhangjiakou",
	"cn-huhehaote",
	"cn-wulanchabu",
	"cn-hangzhou",
	"cn-shanghai",
	"cn-shenzhen",
	"cn-heyuan",
	"cn-guangzhou",
	"cn-chengdu",
	"cn-hongkong",
	"ap-northeast-1",
	"ap-southeast-1",
	"ap-southeast-2",
	"ap-southeast-3",
	"ap-southeast-6",
	"ap-southeast-5",
	"ap-south-1",
	"us-east-1",
	"us-west-1",
	"eu-west-1",
	"eu-central-1",

	// These regions are local cloud
	// "cn-nanjing",
	// "cn-fuzhou",

	// Dubai does not support private zone service
	// "me-east-1",

	// Some resources are not supported
	// "me-central-1",
	// "ap-northeast-2",
	// "ap-southeast-7",
)

// Platform collects AlibabaCloud-specific configuration.
func Platform() (*alibabacloud.Platform, error) {
	client, err := NewClient(defaultRegion)
	if err != nil {
		return nil, err
	}

	err = bypassDeprecation()
	if err != nil {
		return nil, err
	}

	region, err := selectRegion(client)
	if err != nil {
		return nil, err
	}

	return &alibabacloud.Platform{
		Region: region,
	}, nil
}

func selectRegion(client *Client) (string, error) {
	regionsResponse, err := client.DescribeRegions()
	if err != nil {
		return "", err
	}
	regions := regionsResponse.Regions.Region

	var defaultLongRegion string
	longRegions := []string{}
	shortRegions := []string{}
	for _, location := range regions {
		if !supportedRegions.Has(location.RegionId) {
			continue
		}
		longRegion := fmt.Sprintf("%s (%s)", location.RegionId, location.LocalName)
		longRegions = append(longRegions, longRegion)
		shortRegions = append(shortRegions, location.RegionId)
		if location.RegionId == defaultRegion {
			defaultLongRegion = longRegion
		}
	}
	if defaultLongRegion == "" {
		return "", errors.Errorf("installer bug: invalid default alibabacloud region %q", defaultRegion)
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
				Default: defaultLongRegion,
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

func bypassDeprecation() error {
	confirmationMsg := "DEPRECATED. Alibaba Cloud is deprecated and will be " +
		"removed in a future OpenShift version. Would you still like to continue?"

	shouldContinue := false
	prompt := &survey.Confirm{
		Message: confirmationMsg,
	}
	err := survey.AskOne(prompt, &shouldContinue)
	if err != nil {
		return err
	}

	if !shouldContinue {
		return errors.Errorf("deprecated platform")
	}

	return nil
}
