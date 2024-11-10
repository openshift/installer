package ibmcloud

import (
	"fmt"
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/ibmcloud/validation"
)

const (
	createNew = "<create new>"

	// BootstrapSGNameSuffix is the suffix value to append for the bootstrap VPC Security Group name.
	BootstrapSGNameSuffix = "security-group-bootstrap"

	// KubernetesAPIPort is the Kubernetes API port.
	KubernetesAPIPort = 6443

	// KubernetesAPIPrivateSuffix is the name suffix for Kubernetes API Private LB resources.
	KubernetesAPIPrivateSuffix = "kubernetes-api-private"

	// KubernetesAPIPublicSuffix is the name suffix for Kubernetes API Public LB resources.
	KubernetesAPIPublicSuffix = "kubernetes-api-public"

	// MachineConfigServerPort is the Machine Config Server port.
	MachineConfigServerPort = 22623

	// MachineConfigSuffix is the name suffix for Machine Config Server LB resources.
	MachineConfigSuffix = "machine-config"
)

// Platform collects IBM Cloud-specific configuration.
func Platform() (*ibmcloud.Platform, error) {
	region, err := selectRegion()
	if err != nil {
		return nil, err
	}

	return &ibmcloud.Platform{
		Region: region,
	}, nil
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

// COSInstanceName creates a COS Instance name based on provided InfraID.
func COSInstanceName(infraID string) string {
	return fmt.Sprintf("%s-cos", infraID)
}

// VSIImageCOSBucketName creates a COS Bucket name for the VSI Image, based on provided InfraID.
func VSIImageCOSBucketName(infraID string) string {
	return fmt.Sprintf("%s-vsi-image", infraID)
}

// VSIImageName creates a VPC VSI Image name, based on provided InfraID.
func VSIImageName(infraID string) string {
	return fmt.Sprintf("%s-rhcos", infraID)
}
