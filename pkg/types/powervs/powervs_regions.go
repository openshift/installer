package powervs

import (
	"fmt"
)

// Since there is no API to query these, we have to hard-code them here.

// Region describes resources associated with a region in Power VS.
// We're using a few items from the IBM Cloud VPC offering. The region names
// for VPC are different so another function of this is to correlate those.
type Region struct {
	Description string
	VPCRegion   string
	Zones       map[string]Zone
}

// Zone holds the sysTypes for a zone in a IBM Power VS region.
type Zone struct {
	SysTypes []string
}

// Regions holds the regions for IBM Power VS, and descriptions used during the survey.
var Regions = map[string]Region{
	"dal": {
		Description: "Dallas, USA",
		VPCRegion:   "us-south",
		Zones: map[string]Zone{
			"dal10": {
				SysTypes: []string{"s922", "s1022", "e980", "e1080"},
			},
			"dal12": {
				SysTypes: []string{"s922", "e980"},
			},
		},
	},
	"eu-de": {
		Description: "Frankfurt, Germany",
		VPCRegion:   "eu-de",
		Zones: map[string]Zone{
			"eu-de-1": {
				SysTypes: []string{"s922", "s1022", "e980"},
			},
			"eu-de-2": {
				SysTypes: []string{"s922", "e980"},
			},
		},
	},
	"lon": {
		Description: "London, UK",
		VPCRegion:   "eu-gb",
		Zones: map[string]Zone{
			"lon04": {
				SysTypes: []string{"s922", "e980"},
			},
			"lon06": {
				SysTypes: []string{"s922", "e980"},
			},
		},
	},
	"mon": {
		Description: "Montreal, Canada",
		VPCRegion:   "ca-tor",
		Zones: map[string]Zone{
			"mon01": {
				SysTypes: []string{"s922", "e980"},
			},
		},
	},
	"osa": {
		Description: "Osaka, Japan",
		VPCRegion:   "jp-osa",
		Zones: map[string]Zone{
			"osa21": {
				SysTypes: []string{"s922", "s1022", "e980"},
			},
		},
	},
	"sao": {
		Description: "São Paulo, Brazil",
		VPCRegion:   "br-sao",
		Zones: map[string]Zone{
			"sao01": {
				SysTypes: []string{"s922", "e980"},
			},
		},
	},
	"syd": {
		Description: "Sydney, Australia",
		VPCRegion:   "au-syd",
		Zones: map[string]Zone{
			"syd04": {
				SysTypes: []string{"s922", "e980"},
			},
			"syd05": {
				SysTypes: []string{"s922", "e980"},
			},
		},
	},
	"tor": {
		Description: "Toronto, Canada",
		VPCRegion:   "ca-tor",
		Zones: map[string]Zone{
			"tor01": {
				SysTypes: []string{"s922", "e980"},
			},
		},
	},
	"us-east": {
		Description: "Washington DC, USA",
		VPCRegion:   "us-east",
		Zones: map[string]Zone{
			"us-east": {
				SysTypes: []string{"s922", "e980"},
			},
		},
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
			if zone == z {
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
			zones = append(zones, z)
		}
	}
	return zones
}

// RegionFromZone returns the region name for a given zone name.
func RegionFromZone(zone string) string {
	for r := range Regions {
		for z := range Regions[r].Zones {
			if zone == z {
				return r
			}
		}
	}
	return ""
}
