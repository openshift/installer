package installconfig

import (
	"bufio"
	"fmt"
	"strings"

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
	platformPrompt = fmt.Sprintf("Platform (%s):", strings.Join(validPlatforms, ", "))
)

// Platform is an asset that queries the user for the platform on which to install
// the cluster.
//
// Contents[0] is the type of the platform.
//
// * AWS
// Contents[1] is the region.
// Contents[2] is the key pair name.
//
// * Libvirt
// Contents[1] is the URI.
// Contents[2] is the SSH key.
type Platform struct {
	InputReader *bufio.Reader
}

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

func (a *Platform) queryUserForPlatform() string {
	for {
		input := asset.QueryUser(a.InputReader, platformPrompt)
		input = strings.ToLower(input)
		for _, p := range validPlatforms {
			if input == p {
				return p
			}
		}
		fmt.Println("Invalid platform")
	}
}

func (a *Platform) awsPlatform() (*asset.State, error) {
	return assetStateForStringContents(
		AWSPlatformType,
		asset.QueryUser(a.InputReader, "Region:"),
		asset.QueryUser(a.InputReader, "Key Pair Name:"),
	), nil
}

func (a *Platform) libvirtPlatform() (*asset.State, error) {
	return assetStateForStringContents(
		LibvirtPlatformType,
		asset.QueryUser(a.InputReader, "URI:"),
		asset.QueryUser(a.InputReader, "SSH Key:"),
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
