package rhcos

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
)

// GCP fetches the URL of the public RHCOS image
func GCP(ctx context.Context, arch types.Architecture, isOKD bool) (string, error) {
	meta, err := fetchRHCOSBuild(ctx, arch, isOKD)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	return fmt.Sprintf("projects/%s/global/images/%s", meta.GCP.Project, meta.GCP.Image), nil
}

// GCPRaw fetches the URL of the public GCP storage bucket containing the RHCOS image
func GCPRaw(ctx context.Context, arch types.Architecture, isOKD bool) (string, error) {
	meta, err := fetchRHCOSBuild(ctx, arch, isOKD)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	return meta.GCP.URL, nil
}
