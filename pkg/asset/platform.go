package asset

import (
	"fmt"
	"strings"
)

const (
	// AWSPlatformType is used to install on AWS.
	AWSPlatformType = "aws"
	// LibvirtPlatformType is used to install of libvirt.
	LibvirtPlatformType = "libvirt"
)

// Platform generates, via user input, the platform upon which to install OpenShift.
type Platform struct{}

var _ Asset = (*Platform)(nil)

// Dependencies returns no dependencies.
func (a *Platform) Dependencies() []Asset {
	return []Asset{}
}

// Generate queries the user for the platform.
func (a *Platform) Generate(map[Asset]*State) (*State, error) {
	platform := queryUserForPlatform()
	return &State{
		Contents: []Content{
			{Data: []byte(platform)},
		},
	}, nil
}

func queryUserForPlatform() string {
	for {
		validPlatforms := []string{AWSPlatformType, LibvirtPlatformType}
		input := queryUser(fmt.Sprintf("Platform (%s): ", strings.Join(validPlatforms, ", ")))
		input = strings.ToLower(input)
		for _, p := range validPlatforms {
			if input == p {
				return p
			}
		}
		fmt.Println("Invalid platform")
	}
}
