package powervs

import (
	"github.com/openshift/installer/pkg/rhcos"
)

func knownRegions() map[string]string {

	regions := make(map[string]string)

	for _, region := range rhcos.PowerVSRegions {
		regions[region["name"]] = region["description"]
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

// Todo(cklokman): Need some form of error handing in this function...
func knownZones(region string) []string {
	return rhcos.PowerVSZones[region]
}

// IsKnownZone return true is a specified zone is Known to the installer.
func IsKnownZone(region string, zone string) bool {
	if _, ok := knownRegions()[region]; ok {
		zones := knownZones(region)
		for _, z := range zones {
			if z == zone {
				return true
			}
		}
		return false
	}
	return false
}
