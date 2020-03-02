package rhcos

import (
	"context"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
)

// VHD fetches the URL of the public Azure storage bucket containing the RHCOS image
func VHD(ctx context.Context, arch types.Architecture, isOKD bool) (string, error) {
	meta, err := fetchRHCOSBuild(ctx, arch, isOKD)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	url := meta.Azure.URL
	if url == "" {
		return "", errors.New("no RHCOS Azure URL found")
	}

	return url, nil
}
