package rhcos

import (
	"context"

	"github.com/pkg/errors"
)

// AMI fetches the HVM AMI ID of the Red Hat Enterprise Linux CoreOS release.
func AMI(ctx context.Context, region string) (string, error) {
	meta, err := fetchRHCOSBuild(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	ami, ok := meta.AMIs[region]
	if !ok {
		return "", errors.Errorf("no RHCOS AMIs found in %s", region)
	}

	return ami.HVM, nil
}
