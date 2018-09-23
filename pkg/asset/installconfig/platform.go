package installconfig

import (
	"fmt"
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
	platform := a.queryUserForPlatform()
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

func (a *Platform) queryUserForPlatform() string {
	var platform string
	survey.AskOne(&survey.Select{
		Message: "Platform",
		Options: validPlatforms,
	}, &platform, nil)

	return platform
}

func (a *Platform) awsPlatform() (*asset.State, error) {
	var region string
	survey.AskOne(&survey.Select{
		Message: "Region",
		Help:    "The AWS region to be used for installation.",
		Default: "us-east-1 (N. Virginia)",
		Options: []string{
			"us-east-2 (Ohio)",
			"us-east-1 (N. Virginia)",
			"us-west-1 (N. California)",
			"us-west-2 (Oregon)",
			"ap-south-1 (Mumbai)",
			"ap-northeast-2 (Seoul)",
			"ap-northeast-3 (Osaka-Local)",
			"ap-southeast-1 (Singapore)",
			"ap-southeast-2 (Sydney)",
			"ap-northeast-1 (Tokyo)",
			"ca-central-1 (Central)",
			"cn-north-1 (Beijing)",
			"cn-northwest-1 (Ningxia)",
			"eu-central-1 (Frankfurt)",
			"eu-west-1 (Ireland)",
			"eu-west-2 (London)",
			"eu-west-3 (Paris)",
			"sa-east-1 (São Paulo)",
		},
	}, &region, nil)

	return assetStateForStringContents(
		AWSPlatformType,
		strings.Split(region, " ")[0],
	), nil
}

func (a *Platform) libvirtPlatform() (*asset.State, error) {
	var uri, image string
	survey.AskOne(&survey.Input{
		Message: "URI",
		Help:    "The libvirt connection URI to be used. This must be accessible from the running cluster.",
		Default: "qemu+tcp://192.168.122.1/system",
	}, &uri, nil)

	survey.AskOne(&survey.Input{
		Message: "Image",
		Help:    "The URL to the OS image.",
		Default: "",
	}, &image, nil)

	return assetStateForStringContents(
		LibvirtPlatformType,
		uri,
		image,
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
