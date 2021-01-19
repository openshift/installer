package rhcos

import (
	"net/url"
)

// GenerateOpenStackImageName returns Glance image name for instances.
func GenerateOpenStackImageName(rhcosImage, infraID string) (imageName string, isURL bool) {
	// Here we check whether rhcosImage is a URL or not. If this is the first case, it means that Glance image
	// should be created by the installer with the universal name "<infraID>-rhcos". Otherwise, it means
	// that we are given the name of the pre-created Glance image, which the installer should use for node
	// provisioning.
	_, err := url.ParseRequestURI(rhcosImage)
	if err != nil {
		return rhcosImage, false
	}

	return infraID + "-rhcos", true
}
