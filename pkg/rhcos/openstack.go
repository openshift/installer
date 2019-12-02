package rhcos

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
)

// OpenStack fetches the URL of the Red Hat Enterprise Linux CoreOS release,
// for the openstack platform
func OpenStack(ctx context.Context) (string, error) {
	meta, err := fetchRHCOSBuild(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	base, err := url.Parse(meta.BaseURI)
	if err != nil {
		return "", err
	}

	relOpenStack, err := url.Parse(meta.Images.OpenStack.Path)
	if err != nil {
		return "", err
	}

	// Attach uncompressed sha256 checksum to the URL
	checkSuffix := "?sha256=" + meta.Images.OpenStack.UncompressedSHA256
	baseURL := base.ResolveReference(relOpenStack).String() + checkSuffix

	// Check that we have generated a valid URL
	_, err = url.ParseRequestURI(baseURL)
	if err != nil {
		return "", err
	}

	return baseURL, nil
}

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
