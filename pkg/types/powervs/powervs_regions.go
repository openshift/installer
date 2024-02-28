package powervs

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
)

// Since there is no API to query these, we have to hard-code them here.

// Region describes resources associated with a region in Power VS.
// We're using a few items from the IBM Cloud VPC offering. The region names
// for VPC are different so another function of this is to correlate those.
type Region struct {
	Description string
	VPCRegion   string
	COSRegion   string
	Zones       []string
	SysTypes    []string
}

// Regions holds the regions for IBM Power VS, and descriptions used during the survey.
var Regions = map[string]Region{
	"dal": {
		Description: "Dallas, USA",
		VPCRegion:   "us-south",
		COSRegion:   "us-south",
		Zones:       []string{"dal10", "dal12"},
		SysTypes:    []string{"s922", "e980"},
	},
	"eu-de": {
		Description: "Frankfurt, Germany",
		VPCRegion:   "eu-de",
		COSRegion:   "eu-de",
		Zones:       []string{"eu-de-1", "eu-de-2"},
		SysTypes:    []string{"s922", "e980"},
	},
	"mad": {
		Description: "Madrid, Spain",
		VPCRegion:   "eu-es",
		COSRegion:   "eu-de", // @HACK - PowerVS says COS not supported in this region
		Zones:       []string{"mad02", "mad04"},
		SysTypes:    []string{"s1022"},
	},
	"sao": {
		Description: "São Paulo, Brazil",
		VPCRegion:   "br-sao",
		COSRegion:   "br-sao",
		Zones:       []string{"sao04"},
		SysTypes:    []string{"s922", "e980"},
	},
	"wdc": {
		Description: "Washington DC, USA",
		VPCRegion:   "us-east",
		COSRegion:   "us-east",
		Zones:       []string{"wdc06", "wdc07"},
		SysTypes:    []string{"s922", "e980"},
	},
}

// VPCRegionForPowerVSRegion returns the VPC region for the specified PowerVS region.
func VPCRegionForPowerVSRegion(region string) (string, error) {
	if r, ok := Regions[region]; ok {
		return r.VPCRegion, nil
	}

	return "", fmt.Errorf("VPC region corresponding to a PowerVS region %s not found ", region)
}

// RegionShortNames returns the list of region names
func RegionShortNames() []string {
	keys := make([]string, len(Regions))
	i := 0
	for r := range Regions {
		keys[i] = r
		i++
	}
	return keys
}

// ValidateVPCRegion validates that given VPC region is known/tested.
func ValidateVPCRegion(region string) bool {
	found := false
	for r := range Regions {
		if region == Regions[r].VPCRegion {
			found = true
			break
		}
	}
	return found
}

// ValidateZone validates that the given zone is known/tested.
func ValidateZone(zone string) bool {
	for r := range Regions {
		for z := range Regions[r].Zones {
			if zone == Regions[r].Zones[z] {
				return true
			}
		}
	}
	return false
}

// ZoneNames returns the list of zone names.
func ZoneNames() []string {
	zones := []string{}
	for r := range Regions {
		for z := range Regions[r].Zones {
			zones = append(zones, Regions[r].Zones[z])
		}
	}
	return zones
}

// RegionFromZone returns the region name for a given zone name.
func RegionFromZone(zone string) string {
	for r := range Regions {
		for z := range Regions[r].Zones {
			if zone == Regions[r].Zones[z] {
				return r
			}
		}
	}
	return ""
}

// AvailableSysTypes returns the default system type for the zone.
func AvailableSysTypes(region string) ([]string, error) {
	knownRegion, ok := Regions[region]
	if !ok {
		return nil, fmt.Errorf("unknown region name provided")
	}
	return knownRegion.SysTypes, nil
}

// AllKnownSysTypes returns aggregated known system types from all regions.
func AllKnownSysTypes() sets.Set[string] {
	sysTypes := sets.New[string]()
	for _, region := range Regions {
		sysTypes.Insert(region.SysTypes...)
	}
	return sysTypes
}

// COSRegionForVPCRegion returns the corresponding COS region for the given VPC region.
func COSRegionForVPCRegion(vpcRegion string) (string, error) {
	for r := range Regions {
		if vpcRegion == Regions[r].VPCRegion {
			return Regions[r].COSRegion, nil
		}
	}

	return "", fmt.Errorf("COS region corresponding to a VPC region %s not found ", vpcRegion)
}

// COSRegionForPowerVSRegion returns the IBM COS region for the specified PowerVS region.
func COSRegionForPowerVSRegion(region string) (string, error) {
	if r, ok := Regions[region]; ok {
		return r.COSRegion, nil
	}

	return "", fmt.Errorf("COS region corresponding to a PowerVS region %s not found ", region)
}

// ValidateCOSRegion validates that given COS region is known/tested.
func ValidateCOSRegion(region string) bool {
	for r := range Regions {
		if region == Regions[r].COSRegion {
			return true
		}
	}
	return false
}
