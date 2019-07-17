package rhcos

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
)

// BaseURI fetches the BaseURI where images can be downloaded from
func BaseURI() (string, error) {
	meta, err := fetchRHCOSBuild(context.TODO())
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	base, err := url.Parse(meta.BaseURI)
	if err != nil {
		return "", err
	}

	return base.String(), nil
}
