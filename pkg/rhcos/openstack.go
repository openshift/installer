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

	return base.ResolveReference(relOpenStack).String(), nil
}
