package asset

import (
	"bufio"
	"fmt"
	"strings"
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

func newPlatform(inputReader *bufio.Reader) *userProvided {
	return &userProvided{
		inputReader: inputReader,
		prompt:      fmt.Sprintf("Platform (%s):", strings.Join(validPlatforms, ", ")),
		validation:  validatePlatformInput,
	}
}

func validatePlatformInput(input string) (string, bool) {
	input = strings.ToLower(input)
	for _, p := range validPlatforms {
		if input == p {
			return p, true
		}
	}
	fmt.Println("Invalid platform")
	return input, false
}
