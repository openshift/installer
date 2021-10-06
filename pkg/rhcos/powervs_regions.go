package rhcos

// Since there is no API to query these, we have to hard-code them here.

// PowerVSRegion describes resources associated with a region in Power VS.
// We're using a few items from the IBM Cloud VPC offering. The region names
// for VPC are different so another function of this is to correlate those.
type PowerVSRegion struct {
	Name        string
	Description string
	VPCRegion   string
	Zones       []string
}

// PowerVSRegions holds the regions for IBM Power VS, and descriptions used during the survey
var PowerVSRegions = map[string]PowerVSRegion{
	"dal": {
		Name:        "dal",
		Description: "Dallas, USA",
		VPCRegion:   "us-south",
		Zones:       []string{"dal12"},
	},
	"eu-de": {
		Name:        "eu-de",
		Description: "Frankfurt, Germany",
		VPCRegion:   "eu-de",
		Zones: []string{
			"eu-de-1",
			"eu-de-2",
		},
	},
	"lon": {
		Name:        "lon",
		Description: "London, UK.",
		VPCRegion:   "eu-gb",
		Zones: []string{
			"lon04",
			"lon06",
		},
	},
	"osa": {
		Name:        "osa",
		Description: "Osaka, Japan",
		VPCRegion:   "jp-osa",
		Zones:       []string{"osa21"},
	},
	"syd": {
		Name:        "syd",
		Description: "Sydney, Australia",
		VPCRegion:   "au-syd",
		Zones:       []string{"syd04"},
	},
	"sao": {
		Name:        "sao",
		Description: "São Paulo, Brazil",
		VPCRegion:   "br-sao",
		Zones:       []string{"sao01"},
	},
	"tor": {
		Name:        "tor",
		Description: "Toronto, Canada",
		VPCRegion:   "ca-tor",
		Zones:       []string{"tor01"},
	},
	"tok": {
		Name:        "tok",
		Description: "Tokyo, Japan",
		VPCRegion:   "jp-tok",
		Zones:       []string{"tok04"},
	},
	"us-east": {
		Name:        "us-east",
		Description: "Washington DC, USA",
		VPCRegion:   "us-east",
		Zones:       []string{"us-east"},
	},
}

// PowerVSZones retrieves a slice of all zones in Power VS
func PowerVSZones() []string {
	var zones []string
	for _, r := range PowerVSRegions {
		zones = append(zones, r.Zones...)
	}
	return zones
}
