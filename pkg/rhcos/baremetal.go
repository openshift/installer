package rhcos

import (
	"context"
	"net/url"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
)

// Baremetal fetches the URL of the Red Hat Enterprise Linux CoreOS release
// for the baremetal platform
func Baremetal(ctx context.Context, arch types.Architecture) (string, error) {
	meta, err := fetchRHCOSBuild(ctx, arch)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	base, err := url.Parse(meta.BaseURI)
	if err != nil {
		return "", err
	}

	relBaremetal, err := url.Parse(meta.Images.Baremetal.Path)
	if err != nil {
		return "", err
	}

	baseURL := base.ResolveReference(relBaremetal).String()

	// Attach sha256 checksum to the URL. Always provide the
	// uncompressed SHA256; the cache will take care of
	// uncompressing before checksumming.
	baseURL += "?sha256=" + meta.Images.Baremetal.UncompressedSHA256

	// Check that we have generated a valid URL
	_, err = url.ParseRequestURI(baseURL)
	if err != nil {
		return "", err
	}

	return baseURL, nil
}
