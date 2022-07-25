package aws

import (
	"github.com/aws/aws-sdk-go/aws/endpoints"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
)

// knownPublicRegions is the subset of public AWS regions where RHEL CoreOS images are published.
// This subset does not include supported regions which are found in other partitions, such as us-gov-east-1.
// Returns: a map of region identifier to region description.
func knownPublicRegions(architecture types.Architecture) map[string]string {
	required := rhcos.AMIRegions(architecture)

	regions := make(map[string]string)
	for _, region := range endpoints.AwsPartition().Regions() {
		if required.Has(region.ID()) {
			regions[region.ID()] = region.Description()
		}
	}
	return regions
}

// IsKnownPublicRegion returns true if a specified region is Known to the installer.
// A known region is the subset of public AWS regions where RHEL CoreOS images are published.
func IsKnownPublicRegion(region string, architecture types.Architecture) bool {
	if _, ok := knownPublicRegions(architecture)[region]; ok {
		return true
	}
	return false
}
