package rhcos

// TODO(cklokman): This entire file should be programatically generated similar to aws

// PowerVSRegion describes resources associated with a region in powervs.
type PowerVSRegion struct {
	Name        string
	Description string
	VPCRegion   string
	Zones       []string
}

// PowerVSRegions holds the regions for IBM Power VS, and descriptions used during the survey
var PowerVSRegions = map[string]PowerVSRegion{
	"dal": PowerVSRegion{
		Name:        "dal",
		Description: "Dallas, USA",
		VPCRegion:   "us-south",
		Zones:       []string{"dal12"},
	},
	"eu-de": PowerVSRegion{
		Name:        "eu-de",
		Description: "Frankfurt, Germany",
		VPCRegion:   "eu-de",
		Zones: []string{
			"eu-de-1",
			"eu-de-2",
		},
	},
	"lon": PowerVSRegion{
		Name:        "lon",
		Description: "London, UK.",
		VPCRegion:   "eu-gb",
		Zones: []string{
			"lon04",
			"lon06",
		},
	},
	"osa": PowerVSRegion{
		Name:        "osa",
		Description: "Osaka, Japan",
		VPCRegion:   "jp-osa",
		Zones:       []string{"osa21"},
	},
	"syd": PowerVSRegion{
		Name:        "syd",
		Description: "Sydney, Australia",
		VPCRegion:   "au-syd",
		Zones:       []string{"syd04"},
	},
	"sao": PowerVSRegion{
		Name:        "sao",
		Description: "São Paulo, Brazil",
		VPCRegion:   "br-sao",
		Zones:       []string{"sao01"},
	},
	"tor": PowerVSRegion{
		Name:        "tor",
		Description: "Toronto, Canada",
		VPCRegion:   "ca-tor",
		Zones:       []string{"tor01"},
	},
	"tok": PowerVSRegion{
		Name:        "tok",
		Description: "Tokyo, Japan",
		VPCRegion:   "jp-tok",
		Zones:       []string{"tok04"},
	},
	"us-east": PowerVSRegion{
		Name:        "us-east",
		Description: "Washington DC, USA",
		VPCRegion:   "us-east",
		Zones:       []string{"us-east"},
	},
}
