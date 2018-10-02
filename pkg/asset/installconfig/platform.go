package installconfig

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
)

const (
	// AWSPlatformType is used to install on AWS.
	AWSPlatformType = "aws"
	// OpenStackPlatformType is used to install on OpenStack.
	OpenStackPlatformType = "openstack"
	// LibvirtPlatformType is used to install of libvirt.
	LibvirtPlatformType = "libvirt"
)

var (
	validPlatforms = []string{AWSPlatformType, OpenStackPlatformType, LibvirtPlatformType}

	validAWSRegions = map[string]string{
		"ap-northeast-1": "Tokyo",
		"ap-northeast-2": "Seoul",
		"ap-northeast-3": "Osaka-Local",
		"ap-south-1":     "Mumbai",
		"ap-southeast-1": "Singapore",
		"ap-southeast-2": "Sydney",
		"ca-central-1":   "Central",
		"cn-north-1":     "Beijing",
		"cn-northwest-1": "Ningxia",
		"eu-central-1":   "Frankfurt",
		"eu-west-1":      "Ireland",
		"eu-west-2":      "London",
		"eu-west-3":      "Paris",
		"sa-east-1":      "São Paulo",
		"us-east-1":      "N. Virginia",
		"us-east-2":      "Ohio",
		"us-west-1":      "N. California",
		"us-west-2":      "Oregon",
	}

	defaultVPCCIDR = "10.0.0.0/16"

	defaultLibvirtNetworkIfName  = "tt0"
	defaultLibvirtNetworkIPRange = "192.168.124.0/24"
	defaultLibvirtImageURL       = "http://aos-ostree.rhev-ci-vms.eng.rdu2.redhat.com/rhcos/images/cloud/latest/rhcos-qemu.qcow2.gz"
)

// Platform is an asset that queries the user for the platform on which to install
// the cluster.
//
// Contents[0] is the type of the platform.
//
// * AWS
// Contents[1] is the region.
//
// * Libvirt
// Contents[1] is the URI.
type Platform struct{}

var _ asset.Asset = (*Platform)(nil)

// Dependencies returns no dependencies.
func (a *Platform) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate queries for input from the user.
func (a *Platform) Generate(map[asset.Asset]*asset.State) (*asset.State, error) {
	platform, err := a.queryUserForPlatform()
	if err != nil {
		return nil, err
	}

	switch platform {
	case AWSPlatformType:
		return a.awsPlatform()
	case OpenStackPlatformType:
		return a.openstackPlatform()
	case LibvirtPlatformType:
		return a.libvirtPlatform()
	default:
		return nil, fmt.Errorf("unknown platform type %q", platform)
	}
}

// Name returns the human-friendly name of the asset.
func (a *Platform) Name() string {
	return "Platform"
}

func (a *Platform) queryUserForPlatform() (string, error) {
	sort.Strings(validPlatforms)
	prompt := asset.UserProvided{
		Question: &survey.Question{
			Prompt: &survey.Select{
				Message: "Platform",
				Options: validPlatforms,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := ans.(string)
				i := sort.SearchStrings(validPlatforms, choice)
				if i == len(validPlatforms) || validPlatforms[i] != choice {
					return fmt.Errorf("invalid platform %q", choice)
				}
				return nil
			}),
		},
		EnvVarName: "OPENSHIFT_INSTALL_PLATFORM",
	}

	platform, err := prompt.Generate(nil)
	if err != nil {
		return "", err
	}

	return string(platform.Contents[0].Data), nil
}

func (a *Platform) awsPlatform() (*asset.State, error) {
	platform := &types.AWSPlatform{
		VPCCIDRBlock: defaultVPCCIDR,
	}

	longRegions := make([]string, 0, len(validAWSRegions))
	shortRegions := make([]string, 0, len(validAWSRegions))
	for id, location := range validAWSRegions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
		shortRegions = append(shortRegions, id)
	}
	regionTransform := survey.TransformString(func(s string) string {
		return strings.SplitN(s, " ", 2)[0]
	})
	sort.Strings(longRegions)
	sort.Strings(shortRegions)
	prompt := asset.UserProvided{
		Question: &survey.Question{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The AWS region to be used for installation.",
				Default: "us-east-1 (N. Virginia)",
				Options: longRegions,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := regionTransform(ans).(string)
				i := sort.SearchStrings(shortRegions, choice)
				if i == len(shortRegions) || shortRegions[i] != choice {
					return fmt.Errorf("invalid region %q", choice)
				}
				return nil
			}),
			Transform: regionTransform,
		},
		EnvVarName: "OPENSHIFT_INSTALL_AWS_REGION",
	}
	region, err := prompt.Generate(nil)
	if err != nil {
		return nil, err
	}
	platform.Region = string(region.Contents[0].Data)

	if value, ok := os.LookupEnv("_CI_ONLY_STAY_AWAY_OPENSHIFT_INSTALL_AWS_USER_TAGS"); ok {
		if err := json.Unmarshal([]byte(value), &platform.UserTags); err != nil {
			return nil, fmt.Errorf("_CI_ONLY_STAY_AWAY_OPENSHIFT_INSTALL_AWS_USER_TAGS contains invalid JSON: %s (%v)", value, err)
		}
	}

	data, err := json.Marshal(platform)
	if err != nil {
		return nil, err
	}

	return &asset.State{
		Contents: []asset.Content{
			{
				Name: "platform",
				Data: []byte(AWSPlatformType),
			},
			{
				Name: "platform.json",
				Data: data,
			},
		},
	}, nil
}

func (a *Platform) openstackPlatform() (*asset.State, error) {
	platform := &types.OpenStackPlatform{
		NetworkCIDRBlock: defaultVPCCIDR,
	}
	prompt := asset.UserProvided{
		Question: &survey.Question{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The OpenStack region to be used for installation.",
				Default: "regionOne",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				//value := ans.(string)
				//FIXME(shardy) add some validation here
				return nil
			}),
		},
		EnvVarName: "OPENSHIFT_INSTALL_OPENSTACK_REGION",
	}
	region, err := prompt.Generate(nil)
	if err != nil {
		return nil, err
	}
	platform.Region = string(region.Contents[0].Data)
	prompt2 := asset.UserProvided{
		Question: &survey.Question{
			Prompt: &survey.Select{
				Message: "Image",
				Help:    "The OpenStack image to be used for installation.",
				Default: "rhcos",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				//value := ans.(string)
				//FIXME(shardy) add some validation here
				return nil
			}),
		},
		EnvVarName: "OPENSHIFT_INSTALL_OPENSTACK_IMAGE",
	}
	image, err := prompt2.Generate(nil)
	if err != nil {
		return nil, err
	}
	platform.BaseImage = string(image.Contents[0].Data)
	prompt3 := asset.UserProvided{
		Question: &survey.Question{
			Prompt: &survey.Select{
				Message: "Cloud",
				Help:    "The OpenStack cloud name from clouds.yaml.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				//value := ans.(string)
				//FIXME(russellb) add some validation here
				return nil
			}),
		},
		EnvVarName: "OPENSHIFT_INSTALL_OPENSTACK_CLOUD",
	}
	cloud, err := prompt3.Generate(nil)
	if err != nil {
		return nil, err
	}
	platform.Cloud = string(cloud.Contents[0].Data)
	prompt4 := asset.UserProvided{
		Question: &survey.Question{
			Prompt: &survey.Select{
				Message: "ExternalNetwork",
				Help:    "The OpenStack external network to be used for installation.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				//value := ans.(string)
				//FIXME(shadower) add some validation here
				return nil
			}),
		},
		EnvVarName: "OPENSHIFT_INSTALL_OPENSTACK_EXTERNAL_NETWORK",
	}
	extNet, err := prompt4.Generate(nil)
	if err != nil {
		return nil, err
	}
	platform.ExternalNetwork = string(extNet.Contents[0].Data)

	data, err := json.Marshal(platform)
	if err != nil {
		return nil, err
	}

	return &asset.State{
		Contents: []asset.Content{
			{
				Name: "platform",
				Data: []byte(OpenStackPlatformType),
			},
			{
				Name: "platform.json",
				Data: data,
			},
		},
	}, nil

}

func (a *Platform) libvirtPlatform() (*asset.State, error) {
	platform := &types.LibvirtPlatform{
		Network: types.LibvirtNetwork{
			IfName:  defaultLibvirtNetworkIfName,
			IPRange: defaultLibvirtNetworkIPRange,
		},
		DefaultMachinePlatform: &types.LibvirtMachinePoolPlatform{
			Image: defaultLibvirtImageURL,
		},
	}

	uriPrompt := asset.UserProvided{
		Question: &survey.Question{
			Prompt: &survey.Input{
				Message: "Libvirt Connection URI",
				Help:    "The libvirt connection URI to be used. This must be accessible from the running cluster.",
				Default: "qemu+tcp://192.168.122.1/system",
			},
			Validate: survey.ComposeValidators(survey.Required, uriValidator),
		},
		EnvVarName: "OPENSHIFT_INSTALL_LIBVIRT_URI",
	}
	uri, err := uriPrompt.Generate(nil)
	if err != nil {
		return nil, err
	}

	imagePrompt := asset.UserProvided{
		Question: &survey.Question{
			Prompt: &survey.Input{
				Message: "Image",
				Help:    "URI of the OS image.",
				Default: platform.DefaultMachinePlatform.Image,
			},
			Validate: survey.ComposeValidators(survey.Required, uriValidator),
		},
		EnvVarName: "OPENSHIFT_INSTALL_LIBVIRT_IMAGE",
	}
	image, err := imagePrompt.Generate(nil)
	if err != nil {
		return nil, err
	}
	platform.URI = string(uri.Contents[0].Data)
	platform.DefaultMachinePlatform.Image = string(image.Contents[0].Data)

	data, err := json.Marshal(platform)
	if err != nil {
		return nil, err
	}

	return &asset.State{
		Contents: []asset.Content{
			{
				Name: "platform-type",
				Data: []byte(LibvirtPlatformType),
			},
			{
				Name: "platform.json",
				Data: data,
			},
		},
	}, nil
}

// uriValidator validates if the answer provided in prompt is a valid
// url and has non-empty scheme.
func uriValidator(ans interface{}) error {
	value := ans.(string)
	uri, err := url.Parse(value)
	if err != nil {
		return err
	}
	if uri.Scheme == "" {
		return fmt.Errorf("invalid URI %q (no scheme)", value)
	}
	return nil
}
