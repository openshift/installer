package rhcos

// TODO(cklokman): This entire file should be programatically generated similar to aws

// PowerVSZones holds the zones cooresponding to each region found in PowerVS Regions
// TODO(cklokman): These are fictional zones for testing
var PowerVSZones = map[string][]string{
	"eu-de": []string{
		"eu-de-1",
		"eu-de-2",
	},
	"dal": []string{"dal12"},
	"lon": []string{
		"lon04",
		"lon06",
	},
	"mon":     []string{"mon01"},
	"osa":     []string{"osa21"},
	"syd":     []string{"syd04"},
	"sao":     []string{"sao01"},
	"tor":     []string{"tor01"},
	"tok":     []string{"tok04"},
	"us-east": []string{"us-east"},
}

// PowerVSRegions holds the regions for IBM Power VS, and descriptions used during the survey
// TODO(cklokman): These are fictional regtions for testing
var PowerVSRegions = []map[string]string{
	map[string]string{
		"name":        "dal",
		"description": "Dallas, USA",
	},
	map[string]string{
		"name":        "eu-de",
		"description": "Frankfurt, Germany",
	},
	map[string]string{
		"name":        "lon",
		"description": "London, UK.",
	},
	map[string]string{
		"name":        "mon",
		"description": "Montreal, Canada",
	},

	map[string]string{
		"name":        "osa",
		"description": "Osaka, Japan",
	},
	map[string]string{
		"name":        "syd",
		"description": "Sydney, Australia",
	},
	map[string]string{
		"name":        "sao",
		"description": "São Paulo, Brazil",
	},
	map[string]string{
		"name":        "tor",
		"description": "Toronto, Canada",
	},
	map[string]string{
		"name":        "tok",
		"description": "Tokyo, Japan",
	},
	map[string]string{
		"name":        "us-east",
		"description": "Washington DC, USA",
	},
}
