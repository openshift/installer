package rhcos

import (
	"context"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// QEMU fetches the URL of the Red Hat Enterprise Linux CoreOS release.
func QEMU(ctx context.Context) (string, error) {
	meta, err := fetchRHCOSBuild(ctx)
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

	// Attach sha256 checksum to the URL. If the file has the ".gz" extension, then the
	// data is compressed and we use SHA256 value; otherwise we work with uncompressed
	// data and therefore need UncompressedSHA256.
	if strings.HasSuffix(baseURL, ".gz") {
		baseURL += "?sha256=" + meta.Images.QEMU.SHA256
	} else {
		baseURL += "?sha256=" + meta.Images.QEMU.UncompressedSHA256
	}

	// Check that we have generated a valid URL
	_, err = url.ParseRequestURI(baseURL)
	if err != nil {
		return "", err
	}

	return baseURL, nil
}
