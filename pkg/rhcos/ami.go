package rhcos

import (
	"context"
	"sort"
	"strings"

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
		regions := make([]string, 0, len(meta.AMIs))
		for rgn := range meta.AMIs {
			regions = append(regions, rgn)
		}
		sort.Strings(regions)

		return "", errors.Errorf("no RHCOS AMIs found in %q (%s)", region, strings.Join(regions, ", "))
	}

	return ami.HVM, nil
}

// AMIRegions returns a set of AWS regions with HVM AMIs of the Red
// Hat Enterprise Linux CoreOS release.
func AMIRegions(ctx context.Context) (map[string]struct{}, error) {
	meta, err := fetchRHCOSBuild(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	exists := struct{}{}
	regions := make(map[string]struct{}, len(meta.AMIs))
	for region := range meta.AMIs {
		regions[region] = exists
	}

	return regions, nil
}
