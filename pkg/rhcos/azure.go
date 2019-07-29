package rhcos

import (
	"context"

	"github.com/pkg/errors"
)

// VHD fetches the URL of the public Azure storage bucket containing the RHCOS image
func VHD(ctx context.Context) (string, error) {
	meta, err := fetchRHCOSBuild(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	url := meta.Azure.URL
	if url == "" {
		return "", errors.New("no RHCOS Azure URL found")
	}

	return url, nil
}
