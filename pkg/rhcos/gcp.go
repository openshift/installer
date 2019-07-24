package rhcos

import (
	"context"

	"github.com/pkg/errors"
)

// GCP fetches the URL of the public GCP storage bucket containing the RHCOS image
func GCP(ctx context.Context) (string, error) {
	meta, err := fetchRHCOSBuild(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	url := meta.GCP.URL
	if url == "" {
		return "", errors.New("no RHCOS GCP URL found")
	}

	return url, nil
}
