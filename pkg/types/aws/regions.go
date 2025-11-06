package aws

var (
	// RegionLookupMap is a static map containing the known AWS regions and the
	// descriptive location information including the Continent and City/Area.
	RegionLookupMap = map[string]string{
		"af-south-1":     "Africa (Cape Town)",
		"ap-east-1":      "Asia Pacific (Hong Kong)",
		"ap-south-2":     "Asia Pacific (Hyderabad)",
		"ap-southeast-3": "Asia Pacific (Jakarta)",
		"ap-southeast-5": "Asia Pacific (Malaysia)",
		"ap-southeast-4": "Asia Pacific (Melbourne)",
		"ap-south-1":     "Asia Pacific (Mumbai)",
		"ap-northeast-3": "Asia Pacific (Osaka)",
		"ap-northeast-2": "Asia Pacific (Seoul)",
		"ap-southeast-1": "Asia Pacific (Singapore)",
		"ap-southeast-2": "Asia Pacific (Sydney)",
		"ap-southeast-6": "Asia Pacific (New Zealand)",
		"ap-east-2":      "Asia Pacific (Taipei)",
		"ap-southeast-7": "Asia Pacific (Thailand)",
		"ap-northeast-1": "Asia Pacific (Tokyo)",
		"us-gov-east-1":  "AWS GovCloud (US-East)",
		"us-gov-west-1":  "AWS GovCloud (US-West)",
		"ca-central-1":   "Canada (Central)",
		"ca-west-1":      "Canada West (Calgary)",
		"eu-central-1":   "Europe (Frankfurt)",
		"eu-west-1":      "Europe (Ireland)",
		"eu-west-2":      "Europe (London)",
		"eu-south-1":     "Europe (Milan)",
		"eu-west-3":      "Europe (Paris)",
		"eu-south-2":     "Europe (Spain)",
		"eu-north-1":     "Europe (Stockholm)",
		"eu-central-2":   "Europe (Zurich)",
		"il-central-1":   "Israel (Tel Aviv)",
		"mx-central-1":   "Mexico (Central)",
		"me-south-1":     "Middle East (Bahrain)",
		"me-central-1":   "Middle East (UAE)",
		"sa-east-1":      "South America (São Paulo)",
		"us-east-1":      "US East (N. Virginia)",
		"us-east-2":      "US East (Ohio)",
		"us-west-1":      "US West (N. California)",
		"us-west-2":      "US West (Oregon)",
	}
)
