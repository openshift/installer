package rhcos

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
)

// GCP fetches the URL of the public RHCOS image
func GCP(ctx context.Context, arch types.Architecture) (string, error) {
	meta, err := fetchRHCOSBuild(ctx, arch)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	return fmt.Sprintf("projects/%s/global/images/%s", meta.GCP.Project, meta.GCP.Image), nil
}
