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

// Partition identifiers.
const (
	AwsPartitionID      = "aws"        // AWS Standard partition.
	AwsCnPartitionID    = "aws-cn"     // AWS China partition.
	AwsUsGovPartitionID = "aws-us-gov" // AWS GovCloud (US) partition.
	AwsIsoPartitionID   = "aws-iso"    // AWS ISO (US) partition.
	AwsIsoBPartitionID  = "aws-iso-b"  // AWS ISOB (US) partition.
)

// AWS Standard partition's regions.
const (
	AfSouth1RegionID     = "af-south-1"     // Africa (Cape Town).
	ApEast1RegionID      = "ap-east-1"      // Asia Pacific (Hong Kong).
	ApNortheast1RegionID = "ap-northeast-1" // Asia Pacific (Tokyo).
	ApNortheast2RegionID = "ap-northeast-2" // Asia Pacific (Seoul).
	ApNortheast3RegionID = "ap-northeast-3" // Asia Pacific (Osaka).
	ApSouth1RegionID     = "ap-south-1"     // Asia Pacific (Mumbai).
	ApSouth2RegionID     = "ap-south-2"     // Asia Pacific (Hyderabad).
	ApSoutheast1RegionID = "ap-southeast-1" // Asia Pacific (Singapore).
	ApSoutheast2RegionID = "ap-southeast-2" // Asia Pacific (Sydney).
	ApSoutheast3RegionID = "ap-southeast-3" // Asia Pacific (Jakarta).
	ApSoutheast4RegionID = "ap-southeast-4" // Asia Pacific (Melbourne).
	CaCentral1RegionID   = "ca-central-1"   // Canada (Central).
	CaWest1RegionID      = "ca-west-1"      // Canada West (Calgary).
	EuCentral1RegionID   = "eu-central-1"   // Europe (Frankfurt).
	EuCentral2RegionID   = "eu-central-2"   // Europe (Zurich).
	EuNorth1RegionID     = "eu-north-1"     // Europe (Stockholm).
	EuSouth1RegionID     = "eu-south-1"     // Europe (Milan).
	EuSouth2RegionID     = "eu-south-2"     // Europe (Spain).
	EuWest1RegionID      = "eu-west-1"      // Europe (Ireland).
	EuWest2RegionID      = "eu-west-2"      // Europe (London).
	EuWest3RegionID      = "eu-west-3"      // Europe (Paris).
	IlCentral1RegionID   = "il-central-1"   // Israel (Tel Aviv).
	MeCentral1RegionID   = "me-central-1"   // Middle East (UAE).
	MeSouth1RegionID     = "me-south-1"     // Middle East (Bahrain).
	SaEast1RegionID      = "sa-east-1"      // South America (Sao Paulo).
	UsEast1RegionID      = "us-east-1"      // US East (N. Virginia).
	UsEast2RegionID      = "us-east-2"      // US East (Ohio).
	UsWest1RegionID      = "us-west-1"      // US West (N. California).
	UsWest2RegionID      = "us-west-2"      // US West (Oregon).
)

// AWS China partition's regions.
const (
	CnNorth1RegionID     = "cn-north-1"     // China (Beijing).
	CnNorthwest1RegionID = "cn-northwest-1" // China (Ningxia).
)

// AWS GovCloud (US) partition's regions.
const (
	UsGovEast1RegionID = "us-gov-east-1" // AWS GovCloud (US-East).
	UsGovWest1RegionID = "us-gov-west-1" // AWS GovCloud (US-West).
)
