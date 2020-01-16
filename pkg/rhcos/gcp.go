package rhcos

import (
	"context"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
)

// GCP fetches the URL of the public GCP storage bucket containing the RHCOS image
func GCP(ctx context.Context, arch types.Architecture) (string, error) {
	meta, err := fetchRHCOSBuild(ctx, arch)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	url := meta.GCP.URL
	if url == "" {
		return "", errors.New("no RHCOS GCP URL found")
	}

	return url, nil
}
