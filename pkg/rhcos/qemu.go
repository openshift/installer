package rhcos

import (
	"context"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
)

// QEMU fetches the URL of the Red Hat Enterprise Linux CoreOS release.
func QEMU(ctx context.Context, arch types.Architecture) (string, error) {
	meta, err := fetchRHCOSBuild(ctx, arch)
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

	// Attach sha256 checksum to the URL. Use compressed SHA256 for known decompressors
	if strings.HasSuffix(baseURL, ".gz") || strings.HasSuffix(baseURL, ".xz") {
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
