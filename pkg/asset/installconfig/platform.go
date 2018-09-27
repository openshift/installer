package installconfig

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey"

	"github.com/openshift/installer/pkg/asset"
)

const (
	// AWSPlatformType is used to install on AWS.
	AWSPlatformType = "aws"
	// LibvirtPlatformType is used to install of libvirt.
	LibvirtPlatformType = "libvirt"
)

var (
	validPlatforms = []string{AWSPlatformType, LibvirtPlatformType}

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

	return assetStateForStringContents(
		AWSPlatformType,
		string(region.Contents[0].Data),
	), nil
}

func (a *Platform) libvirtPlatform() (*asset.State, error) {
	prompt := asset.UserProvided{
		Question: &survey.Question{
			Prompt: &survey.Input{
				Message: "URI",
				Help:    "The libvirt connection URI to be used. This must be accessible from the running cluster.",
				Default: "qemu+tcp://192.168.122.1/system",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				value := ans.(string)
				uri, err := url.Parse(value)
				if err != nil {
					return err
				}
				if uri.Scheme == "" {
					return fmt.Errorf("invalid URI %q (no scheme)", value)
				}
				return nil
			}),
		},
		EnvVarName: "OPENSHIFT_INSTALL_LIBVIRT_URI",
	}
	uri, err := prompt.Generate(nil)
	if err != nil {
		return nil, err
	}

	return assetStateForStringContents(
		LibvirtPlatformType,
		string(uri.Contents[0].Data),
	), nil
}

func assetStateForStringContents(s ...string) *asset.State {
	c := make([]asset.Content, len(s))
	for i, d := range s {
		c[i] = asset.Content{
			Data: []byte(d),
		}
	}
	return &asset.State{
		Contents: c,
	}
}
