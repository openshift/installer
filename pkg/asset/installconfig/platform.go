package installconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/rhcos"
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
	defaultLibvirtNetworkIPRange = "192.168.126.0/24"
)

// Platform is an asset that queries the user for the platform on which to install
// the cluster.
type platform struct {
	AWS       *types.AWSPlatform
	Openstack *types.OpenStackPlatform
	Libvirt   *types.LibvirtPlatform
}

var _ asset.Asset = (*platform)(nil)

// Dependencies returns no dependencies.
func (a *platform) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate queries for input from the user.
func (a *platform) Generate(asset.Parents) error {
	platform, err := a.queryUserForPlatform()
	if err != nil {
		return err
	}

	switch platform {
	case AWSPlatformType:
		aws, err := a.awsPlatform()
		if err != nil {
			return err
		}
		a.AWS = aws
	case OpenStackPlatformType:
		openstack, err := a.openstackPlatform()
		if err != nil {
			return err
		}
		a.Openstack = openstack
	case LibvirtPlatformType:
		libvirt, err := a.libvirtPlatform()
		if err != nil {
			return err
		}
		a.Libvirt = libvirt
	default:
		return fmt.Errorf("unknown platform type %q", platform)
	}

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *platform) Name() string {
	return "Platform"
}

func (a *platform) queryUserForPlatform() (string, error) {
	sort.Strings(validPlatforms)
	return asset.GenerateUserProvidedAsset(
		"Platform",
		&survey.Question{
			Prompt: &survey.Select{
				Message: "Platform",
				Options: validPlatforms,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := ans.(string)
				i := sort.SearchStrings(validPlatforms, choice)
				if i == len(validPlatforms) || validPlatforms[i] != choice {
					return errors.Errorf("invalid platform %q", choice)
				}
				return nil
			}),
		},
		"OPENSHIFT_INSTALL_PLATFORM",
	)
}

func (a *platform) awsPlatform() (*types.AWSPlatform, error) {
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
	region, err := asset.GenerateUserProvidedAsset(
		"AWS Region",
		&survey.Question{
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
					return errors.Errorf("invalid region %q", choice)
				}
				return nil
			}),
			Transform: regionTransform,
		},
		"OPENSHIFT_INSTALL_AWS_REGION",
	)
	if err != nil {
		return nil, err
	}

	userTags := map[string]string{}
	if value, ok := os.LookupEnv("_CI_ONLY_STAY_AWAY_OPENSHIFT_INSTALL_AWS_USER_TAGS"); ok {
		if err := json.Unmarshal([]byte(value), &userTags); err != nil {
			return nil, errors.Wrapf(err, "_CI_ONLY_STAY_AWAY_OPENSHIFT_INSTALL_AWS_USER_TAGS contains invalid JSON: %s", value)
		}
	}

	return &types.AWSPlatform{
		VPCCIDRBlock: defaultVPCCIDR,
		Region:       region,
		UserTags:     userTags,
	}, nil
}

func (a *platform) openstackPlatform() (*types.OpenStackPlatform, error) {
	region, err := asset.GenerateUserProvidedAsset(
		"OpenStack Region",
		&survey.Question{
			Prompt: &survey.Input{
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
		"OPENSHIFT_INSTALL_OPENSTACK_REGION",
	)
	if err != nil {
		return nil, err
	}

	image, err := asset.GenerateUserProvidedAsset(
		"OpenStack Image",
		&survey.Question{
			Prompt: &survey.Input{
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
		"OPENSHIFT_INSTALL_OPENSTACK_IMAGE",
	)
	if err != nil {
		return nil, err
	}

	cloud, err := asset.GenerateUserProvidedAsset(
		"OpenStack Cloud",
		&survey.Question{
			//TODO(russellb) - We could open clouds.yaml here and read the list of defined clouds
			//and then use survey.Select to let the user choose one.
			Prompt: &survey.Input{
				Message: "Cloud",
				Help:    "The OpenStack cloud name from clouds.yaml.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				//value := ans.(string)
				//FIXME(russellb) add some validation here
				return nil
			}),
		},
		"OPENSHIFT_INSTALL_OPENSTACK_CLOUD",
	)
	if err != nil {
		return nil, err
	}

	extNet, err := asset.GenerateUserProvidedAsset(
		"OpenStack External Network",
		&survey.Question{
			Prompt: &survey.Input{
				Message: "ExternalNetwork",
				Help:    "The OpenStack external network to be used for installation.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				//value := ans.(string)
				//FIXME(shadower) add some validation here
				return nil
			}),
		},
		"OPENSHIFT_INSTALL_OPENSTACK_EXTERNAL_NETWORK",
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to Marshal %s platform", OpenStackPlatformType)
	}

	return &types.OpenStackPlatform{
		NetworkCIDRBlock: defaultVPCCIDR,
		Region:           region,
		BaseImage:        image,
		Cloud:            cloud,
		ExternalNetwork:  extNet,
	}, nil
}

func (a *platform) libvirtPlatform() (*types.LibvirtPlatform, error) {
	uri, err := asset.GenerateUserProvidedAsset(
		"Libvirt Connection URI",
		&survey.Question{
			Prompt: &survey.Input{
				Message: "Libvirt Connection URI",
				Help:    "The libvirt connection URI to be used. This must be accessible from the running cluster.",
				Default: "qemu+tcp://192.168.122.1/system",
			},
			Validate: survey.ComposeValidators(survey.Required, uriValidator),
		},
		"OPENSHIFT_INSTALL_LIBVIRT_URI",
	)
	if err != nil {
		return nil, err
	}

	// TODO: Ideally, this would live inside of a closure which is passed to
	//       asset.GenerateUserProvidedAsset and only called if the environment
	//       variable isn't present. As this exists, it ruins the abstraction.
	var qcowImage string
	if _, ok := os.LookupEnv("OPENSHIFT_INSTALL_LIBVIRT_IMAGE"); !ok {
		qcowImage, err = rhcos.QEMU(context.TODO(), rhcos.DefaultChannel)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch QEMU image URL")
		}
	}

	image, err := asset.GenerateUserProvidedAsset(
		"Libvirt Image",
		&survey.Question{
			Prompt: &survey.Input{
				Message: "Image",
				Help:    "URI of the OS image.",
				Default: qcowImage,
			},
			Validate: survey.ComposeValidators(survey.Required, uriValidator),
		},
		"OPENSHIFT_INSTALL_LIBVIRT_IMAGE",
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to Marshal %s platform", LibvirtPlatformType)
	}

	return &types.LibvirtPlatform{
		Network: types.LibvirtNetwork{
			IfName:  defaultLibvirtNetworkIfName,
			IPRange: defaultLibvirtNetworkIPRange,
		},
		DefaultMachinePlatform: &types.LibvirtMachinePoolPlatform{
			Image: image,
		},
		URI: uri,
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
