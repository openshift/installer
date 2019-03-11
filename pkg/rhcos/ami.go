package rhcos

import (
	"context"

	"github.com/pkg/errors"
)

// AMI fetches the HVM AMI ID of the latest Red Hat Enterprise Linux CoreOS release.
func AMI(ctx context.Context, region string) (string, error) {
	meta, err := fetchLatestMetadata(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	for _, ami := range meta.AMIs {
		if ami.Name == region {
			return ami.HVM, nil
		}
	}

	return "", errors.Errorf("no RHCOS AMIs found in %s", region)
}
