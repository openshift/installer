package rhcos

import (
	"context"
	"net/url"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
)

// QEMU fetches the URL of the Red Hat Enterprise Linux CoreOS release.
func QEMU(ctx context.Context, arch types.Architecture, isOKD bool) (string, error) {
	meta, err := fetchRHCOSBuild(ctx, arch, isOKD)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	base, err := url.Parse(meta.BaseURI)
	if err != nil {
		return "", err
	}

	relQEMU, err := url.Parse(meta.Images.QEMU.Path)
	if err != nil {
		return "", err
	}

	baseURL := base.ResolveReference(relQEMU).String()

	// Attach sha256 checksum to the URL.  Always provide the
	// uncompressed SHA256; the cache will take care of
	// uncompressing before checksumming.
	baseURL += "?sha256=" + meta.Images.QEMU.UncompressedSHA256

	// Check that we have generated a valid URL
	_, err = url.ParseRequestURI(baseURL)
	if err != nil {
		return "", err
	}

	return baseURL, nil
}
