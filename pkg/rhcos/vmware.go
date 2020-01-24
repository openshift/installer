package rhcos

import (
	"context"
	"net/url"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
)

// VMware fetches the URL of the Red Hat Enterprise Linux CoreOS release.
func VMware(ctx context.Context, arch types.Architecture) (string, error) {
	meta, err := fetchRHCOSBuild(ctx, arch)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	base, err := url.Parse(meta.BaseURI)
	if err != nil {
		return "", err
	}

	image, err := url.Parse(meta.Images.VMware.Path)
	if err != nil {
		return "", err
	}

	baseURL := base.ResolveReference(image).String()

	// TODO: Get Uncompressed SHA256 into rhcos.json
	// Attach sha256 checksum to the URL.  Always provide the
	// uncompressed SHA256; the cache will take care of
	// uncompressing before checksumming.
	//baseURL += "?sha256=" + meta.Images.VMware.UncompressedSHA256

	// Check that we have generated a valid URL
	_, err = url.ParseRequestURI(baseURL)
	if err != nil {
		return "", err
	}

	return baseURL, nil
}

// GenerateVSphereImageName returns Glance image name for instances.
func GenerateVSphereImageName(rhcosImage, infraID string) (string, error) {

	//TODO: Determine how to use Base Path as imageName
	/*
		rhcosImageURL, err := url.ParseRequestURI(rhcosImage)
		if err != nil {
			return "", errors.Errorf("expected RHCOS image to be a URL")
		}
		fileName := path.Base(rhcosImageURL.Path)
		fileNameNoExt := strings.TrimSuffix(fileName, path.Ext(fileName))
		imageName := fmt.Sprintf("%s-%s", infraID, fileNameNoExt)

		return imageName, nil
	*/

	return infraID + "-rhcos", nil
}
