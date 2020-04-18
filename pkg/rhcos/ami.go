//go:generate go run ami_regions_generate.go rhcos ../../data/data/rhcos-amd64.json ami_regions.go

package rhcos

import (
	"context"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
)

// AMI fetches the HVM AMI ID of the Red Hat Enterprise Linux CoreOS release.
func AMI(ctx context.Context, arch types.Architecture, region string) (string, error) {
	meta, err := fetchRHCOSBuild(ctx, arch)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	ami, ok := meta.AMIs[region]
	if !ok {
		return "", errors.Errorf("no RHCOS AMIs found in %s", region)
	}

	return ami.HVM, nil
}
