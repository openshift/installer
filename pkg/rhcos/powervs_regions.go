package rhcos

// TODO(cklokman): This entire file should be programatically generated similar to aws

// PowerVSZones holds the zones cooresponding to each region found in PowerVS Regions
// TODO(cklokman): These are fictional zones for testing
var PowerVSZones = map[string][]string{
	"us-south": []string{
		"us-south-zone1",
		"us-south-zone2",
	},
}

// PowerVSRegions holds the regions for IBM Power VS, and descriptions used during the survey
// TODO(cklokman): These are fictional regtions for testing
var PowerVSRegions = []map[string]string{
	map[string]string{
		"name":        "us-south",
		"description": "This is the us-south test region.",
	},
}
