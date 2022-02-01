package aws

import (
	"github.com/aws/aws-sdk-go/aws/endpoints"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
)

// knownRegions is a list of AWS regions that the installer recognizes.
// This is subset of AWS regions and the regions where RHEL CoreOS images are published.
// The result is a map of region identifier to region description
func knownRegions(architecture types.Architecture) map[string]string {
	required := rhcos.AMIRegions(architecture)

	regions := make(map[string]string)
	for _, region := range endpoints.AwsPartition().Regions() {
		if required.Has(region.ID()) {
			regions[region.ID()] = region.Description()
		}
	}
	return regions
}

// IsKnownRegion return true is a specified region is Known to the installer.
// A known region is subset of AWS regions and the regions where RHEL CoreOS images are published.
func IsKnownRegion(region string, architecture types.Architecture) bool {
	if _, ok := knownRegions(architecture)[region]; ok {
		return true
	}
	return false
}
