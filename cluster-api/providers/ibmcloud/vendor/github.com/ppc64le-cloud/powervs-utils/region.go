package utils

import (
	"fmt"
	"strings"
)

// SysTypes denotes the available system types against individual zones
// The system types against each zone can be fetched using either of the following:
// 1. ibmcloud cli command
//    `ibmcloud pi datacenter ls --json | jq -r '.datacenters[] | "\(.location.region),\(.capabilitiesDetails.supportedSystems.general[])"' | sort -u`
// 2. https://cloud.ibm.com/apidocs/power-cloud#v1-datacenters-getall (the above command uses this api underneath).

type SysTypes []string

func GetRegion(zone string) (region string, err error) {
	err = nil
	switch {
	case strings.HasPrefix(zone, "us-south"):
		region = "us-south"
	case strings.HasPrefix(zone, "dal"):
		region = "dal"
	case strings.HasPrefix(zone, "sao"):
		region = "sao"
	case strings.HasPrefix(zone, "us-east"):
		region = "us-east"
	case strings.HasPrefix(zone, "eu-de-"):
		region = "eu-de"
	case strings.HasPrefix(zone, "lon"):
		region = "lon"
	case strings.HasPrefix(zone, "syd"):
		region = "syd"
	case strings.HasPrefix(zone, "tok"):
		region = "tok"
	case strings.HasPrefix(zone, "osa"):
		region = "osa"
	case strings.HasPrefix(zone, "mon"):
		region = "mon"
	case strings.HasPrefix(zone, "mad"):
		region = "mad"
	case strings.HasPrefix(zone, "wdc"):
		region = "wdc"
	case strings.HasPrefix(zone, "tor"):
		region = "tor"
	case strings.HasPrefix(zone, "che"):
		region = "che"
	default:
		return "", fmt.Errorf("region not found for the zone, talk to the developer to add the support into the tool: %s", zone)
	}
	return
}

// Region describes respective IBM Cloud COS region, VPC region and Zones associated with a region in Power VS.
type Region struct {
	Description string
	VPCRegion   string
	COSRegion   string
	Zones       map[string]SysTypes
	VPCZones    []string
}

// regions provides a mapping between Power VS and IBM Cloud VPC and IBM COS regions.
var regions = map[string]Region{
	"dal": {
		Description: "Dallas, USA",
		VPCRegion:   "us-south",
		COSRegion:   "us-south",
		Zones: map[string]SysTypes{
			"dal10": {"e1080", "e980", "s1022", "s922"},
			"dal12": {"e980", "s922"},
			"dal14": {"e1050", "e1080", "s1022"},
		},
		VPCZones: []string{"us-south-1", "us-south-2", "us-south-3"},
	},
	"eu-de": {
		Description: "Frankfurt, Germany",
		VPCRegion:   "eu-de",
		COSRegion:   "eu-de",
		Zones: map[string]SysTypes{
			"eu-de-1": {"e1080", "e1050", "e980", "s1022", "s922"},
			"eu-de-2": {"e980", "s1022", "s922"},
		},
		VPCZones: []string{"eu-de-1", "eu-de-2", "eu-de-3"},
	},
	"lon": {
		Description: "London, UK.",
		VPCRegion:   "eu-gb",
		COSRegion:   "eu-gb",
		Zones: map[string]SysTypes{
			"lon04": {"e980", "s922"},
			"lon06": {"e980", "s922"},
		},
		VPCZones: []string{"eu-gb-1", "eu-gb-2", "eu-gb-3"},
	},
	"mad": {
		Description: "Madrid, Spain",
		VPCRegion:   "eu-es",
		COSRegion:   "eu-de", // @HACK - PowerVS says COS not supported in this region
		Zones: map[string]SysTypes{
			"mad02": {"e1050", "e1080", "e980", "s1022", "s922"},
			"mad04": {"e980", "s1022"},
		},
		VPCZones: []string{"eu-es-1", "eu-es-2", "eu-es-3"},
	},
	"mon": {
		Description: "Montreal, Canada",
		VPCRegion:   "",
		COSRegion:   "ca-tor",
		Zones: map[string]SysTypes{
			"mon01": {"e980", "s922"},
		},
		VPCZones: []string{},
	},
	"osa": {
		Description: "Osaka, Japan",
		VPCRegion:   "jp-osa",
		COSRegion:   "jp-osa",
		Zones: map[string]SysTypes{
			"osa21": {"e980", "s1022", "s922"},
		},
		VPCZones: []string{"jp-osa-1", "jp-osa-2", "jp-osa-3"},
	},
	"sao": {
		Description: "São Paulo, Brazil",
		VPCRegion:   "br-sao",
		COSRegion:   "br-sao",
		Zones: map[string]SysTypes{
			"sao01": {"e980", "s1022", "s922"},
			"sao04": {"e980", "s1022", "s922"},
		},
		VPCZones: []string{"br-sao-1", "br-sao-2", "br-sao-3"},
	},
	"syd": {
		Description: "Sydney, Australia",
		VPCRegion:   "au-syd",
		COSRegion:   "au-syd",
		Zones: map[string]SysTypes{
			"syd04": {"e980", "s922"},
			"syd05": {"e980", "s922"},
		},
		VPCZones: []string{"au-syd-1", "au-syd-2", "au-syd-3"},
	},
	"tok": {
		Description: "Tokyo, Japan",
		VPCRegion:   "jp-tok",
		COSRegion:   "jp-tok",
		Zones: map[string]SysTypes{
			"tok04": {"e980", "s1022", "s922"},
		},
		VPCZones: []string{"jp-tok-1", "jp-tok-2", "jp-tok-3"},
	},
	"tor": {
		Description: "Toronto, Canada",
		VPCRegion:   "ca-tor",
		COSRegion:   "ca-tor",
		Zones: map[string]SysTypes{
			"tor01": {"e980", "s922"},
		},
		VPCZones: []string{"ca-tor-1", "ca-tor-2", "ca-tor-3"},
	}, // Keeping us-east and us-south zones as individual entries to easily map the respective VPC and COS regions by matching the prefix of the zone like in GetRegion()
	"us-east": {
		Description: "Washington DC, USA",
		VPCRegion:   "us-east",
		COSRegion:   "us-east",
		Zones: map[string]SysTypes{
			"us-east": {"e980", "s922"},
		},
		VPCZones: []string{"us-east-1", "us-east-2", "us-east-3"},
	},
	"us-south": {
		Description: "Dallas, USA",
		VPCRegion:   "us-south",
		COSRegion:   "us-south",
		Zones: map[string]SysTypes{
			"us-south": {"e980", "s922"},
		},
		VPCZones: []string{"us-south-1", "us-south-2", "us-south-3"},
	},
	"wdc": {
		Description: "Washington DC, USA",
		VPCRegion:   "us-east",
		COSRegion:   "us-east",
		Zones: map[string]SysTypes{
			"wdc06": {"e980", "s1022", "s922"},
			"wdc07": {"e1050", "e1080", "e980", "s1022", "s922"},
		},
		VPCZones: []string{"us-east-1", "us-east-2", "us-east-3"},
	},
	"che": {
		Description: "Chennai, India",
		VPCRegion:   "",
		COSRegion:   "",
		Zones: map[string]SysTypes{
			"che01": {"e980", "s922"},
		},
		VPCZones: []string{},
	},
}

// COSRegionForVPCRegion returns the corresponding COS region for the given VPC region
func COSRegionForVPCRegion(vpcRegion string) (string, error) {
	for r := range regions {
		if vpcRegion == regions[r].VPCRegion {
			return regions[r].COSRegion, nil
		}
	}

	return "", fmt.Errorf("COS region corresponding to a VPC region %s not found ", vpcRegion)
}

// VPCRegionForPowerVSRegion returns the VPC region for the specified PowerVS region.
func VPCRegionForPowerVSRegion(region string) (string, error) {
	if r, ok := regions[region]; ok {
		return r.VPCRegion, nil
	}

	return "", fmt.Errorf("VPC region corresponding to a PowerVS region %s not found ", region)
}

// COSRegionForPowerVSRegion returns the IBM COS region for the specified PowerVS region.
func COSRegionForPowerVSRegion(region string) (string, error) {
	if r, ok := regions[region]; ok {
		return r.COSRegion, nil
	}

	return "", fmt.Errorf("COS region corresponding to a PowerVS region %s not found ", region)
}

// ValidateVPCRegion validates that given VPC region is known/tested.
func ValidateVPCRegion(region string) bool {
	for r := range regions {
		if region == regions[r].VPCRegion {
			return true
		}
	}
	return false
}

// ValidateCOSRegion validates that given COS region is known/tested.
func ValidateCOSRegion(region string) bool {
	for r := range regions {
		if region == regions[r].COSRegion {
			return true
		}
	}
	return false
}

// RegionShortNames returns the list of region names
func RegionShortNames() []string {
	var keys []string
	for r := range regions {
		keys = append(keys, r)
	}
	return keys
}

// ValidateZone validates that the given zone is known/tested.
func ValidateZone(zone string) bool {
	for r := range regions {
		_, exists := regions[r].Zones[zone]
		if exists {
			return exists
		}
	}
	return false
}

// ZoneNames returns the list of zone names.
func ZoneNames() []string {
	zones := []string{}
	for r := range regions {
		for zone := range regions[r].Zones {
			zones = append(zones, zone)
		}
	}
	return zones
}

// RegionFromZone returns the region name for a given zone name.
func RegionFromZone(zone string) string {
	for r := range regions {
		for z := range regions[r].Zones {
			if zone == z {
				return r
			}
		}
	}
	return ""
}

// AvailableSysTypes returns the supported system types for the zone.
func AvailableSysTypes(region, zone string) ([]string, error) {
	knownRegion, ok := regions[region]
	if !ok {
		return nil, fmt.Errorf("unknown region name provided")
	}
	knownZone, ok := knownRegion.Zones[zone]
	if !ok {
		return nil, fmt.Errorf("unknown zone name provided")
	}
	return knownZone, nil
}

// IsGlobalRoutingRequiredForTG returns true when powervs and vpc regions are different.
func IsGlobalRoutingRequiredForTG(powerVSRegion string, vpcRegion string) bool {
	if r, ok := regions[powerVSRegion]; ok && r.VPCRegion == vpcRegion {
		return false
	}
	return true
}

// VPCZonesForVPCRegion returns the VPC zones associated with the VPC region.
func VPCZonesForVPCRegion(region string) ([]string, error) {
	for _, regionDetails := range regions {
		if regionDetails.VPCRegion == region {
			return regionDetails.VPCZones, nil
		}
	}
	return nil, fmt.Errorf("VPC zones corresponding to the VPC region %s is not found", region)
}
