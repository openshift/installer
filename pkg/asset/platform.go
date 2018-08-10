package asset

import (
	"fmt"
	"strings"
)

const (
	AWSPlatformType     = "aws"
	LibvirtPlatformType = "libvirt"
)

type Platform struct{}

var _ Asset = (*Platform)(nil)

func (a *Platform) Dependencies() []Asset {
	return []Asset{}
}

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
