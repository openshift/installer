package powervs

import (
	"fmt"
)

// Since there is no API to query these, we have to hard-code them here.

// VPCRegion describes zones associated with a region for IBM Cloud VPC.
type VPCRegion struct {
	Description string
	Zones       []string
}

// VPCRegions holds the region names for IBM Cloud VPC
var VPCRegions = map[string]VPCRegion{
	"au-syd": {
		Description: "Sydney, Australia",
		Zones: []string{
			"au-syd-1",
			"au-syd-2",
			"au-syd-3",
		},
	},
	"br-sao": {
		Description: "São Paulo, Brazil",
		Zones: []string{
			"br-sao-1",
			"br-sao-2",
			"br-sao-3",
		},
	},
	"ca-tor": {
		Description: "Toronto, Canada",
		Zones: []string{
			"ca-tor-1",
			"ca-tor-2",
			"ca-tor-3",
		},
	},
	"eu-de": {
		Description: "Frankfurt, Germany",
		Zones: []string{
			"eu-de-1",
			"eu-de-2",
			"eu-de-3",
		},
	},
	"eu-gb": {
		Description: "London, UK.",
		Zones: []string{
			"eu-gb-1",
			"eu-gb-2",
			"eu-gb-3",
		},
	},
	"jp-osa": {
		Description: "Osaka, Japan",
		Zones: []string{
			"jp-osa-1",
			"jp-osa-2",
			"jp-osa-3",
		},
	},
	"jp-tok": {
		Description: "Tokyo, Japan",
		Zones: []string{
			"jp-tok-1",
			"jp-tok-2",
			"jp-tok-3",
		},
	},
	"us-east": {
		Description: "Washington DC, USA",
		Zones: []string{
			"us-east-1",
			"us-east-2",
			"us-east-3",
		},
	},
	"us-south": {
		Description: "Dallas, USA",
		Zones: []string{
			"us-south-1",
			"us-south-2",
			"us-south-3",
		},
	},
}

// VPCRegionForVPCZone returns the VPC region for the specified VPC zone.
func VPCRegionForVPCZone(vpczone string) (string, error) {
	for rk := range VPCRegions {
		for zk := range VPCRegions[rk].Zones {
			if VPCRegions[rk].Zones[zk] == vpczone {
				return rk, nil
			}
		}
	}
	return "", fmt.Errorf("VPC region corresponding to given zone %s not found ", vpczone)
}
