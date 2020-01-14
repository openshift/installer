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

	// Attach sha256 checksum to the URL.  Always provide the
	// uncompressed SHA256; the cache will take care of
	// uncompressing before checksumming.
	baseURL += "?sha256=" + meta.Images.VMware.SHA256

	// Check that we have generated a valid URL
	_, err = url.ParseRequestURI(baseURL)
	if err != nil {
		return "", err
	}

	return baseURL, nil
}
