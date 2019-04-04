package rhcos

import (
	"context"
	"net/url"

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

	return base.ResolveReference(relQEMU).String(), nil
}
