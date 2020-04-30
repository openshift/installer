package aws

import (
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/rhcos"
)

// knownRegions is a list of AWS regions that the installer recognizes.
// This is subset of AWS regions and the regions where RHEL CoreOS images are published.
// The result is a map of region identifier to region description
func knownRegions() map[string]string {
	required := sets.NewString(rhcos.AMIRegions...)

	regions := make(map[string]string)
	for _, partition := range endpoints.DefaultPartitions() {
		for _, partitionRegion := range partition.Regions() {
			partitionRegion := partitionRegion
			if required.Has(partitionRegion.ID()) {
				regions[partitionRegion.ID()] = partitionRegion.Description()
			}
		}
	}
	return regions
}

// IsKnownRegion return true is a specified region is Known to the installer.
// A known region is subset of AWS regions and the regions where RHEL CoreOS images are published.
func IsKnownRegion(region string) bool {
	if _, ok := knownRegions()[region]; ok {
		return true
	}
	return false
}
