package aws

// AWS SDK Go v2 does not expose region and partition constants; thus
// we need to define them in the installer code.
var (
	// RegionLookupMap is a static map containing the known AWS regions and the
	// descriptive location information including the Continent and City/Area.
	RegionLookupMap = map[string]string{
		AfSouth1RegionID:     "Africa (Cape Town)",
		ApEast1RegionID:      "Asia Pacific (Hong Kong)",
		ApSouth2RegionID:     "Asia Pacific (Hyderabad)",
		ApSoutheast3RegionID: "Asia Pacific (Jakarta)",
		ApSoutheast5RegionID: "Asia Pacific (Malaysia)",
		ApSoutheast4RegionID: "Asia Pacific (Melbourne)",
		ApSouth1RegionID:     "Asia Pacific (Mumbai)",
		ApNortheast3RegionID: "Asia Pacific (Osaka)",
		ApNortheast2RegionID: "Asia Pacific (Seoul)",
		ApSoutheast1RegionID: "Asia Pacific (Singapore)",
		ApSoutheast2RegionID: "Asia Pacific (Sydney)",
		ApSoutheast6RegionID: "Asia Pacific (New Zealand)",
		ApEast2RegionID:      "Asia Pacific (Taipei)",
		ApSoutheast7RegionID: "Asia Pacific (Thailand)",
		ApNortheast1RegionID: "Asia Pacific (Tokyo)",
		UsGovEast1RegionID:   "AWS GovCloud (US-East)",
		UsGovWest1RegionID:   "AWS GovCloud (US-West)",
		CaCentral1RegionID:   "Canada (Central)",
		CaWest1RegionID:      "Canada West (Calgary)",
		EuCentral1RegionID:   "Europe (Frankfurt)",
		EuWest1RegionID:      "Europe (Ireland)",
		EuWest2RegionID:      "Europe (London)",
		EuSouth1RegionID:     "Europe (Milan)",
		EuWest3RegionID:      "Europe (Paris)",
		EuSouth2RegionID:     "Europe (Spain)",
		EuNorth1RegionID:     "Europe (Stockholm)",
		EuCentral2RegionID:   "Europe (Zurich)",
		IlCentral1RegionID:   "Israel (Tel Aviv)",
		MxCentral1RegionID:   "Mexico (Central)",
		MeSouth1RegionID:     "Middle East (Bahrain)",
		MeCentral1RegionID:   "Middle East (UAE)",
		SaEast1RegionID:      "South America (São Paulo)",
		UsEast1RegionID:      "US East (N. Virginia)",
		UsEast2RegionID:      "US East (Ohio)",
		UsWest1RegionID:      "US West (N. California)",
		UsWest2RegionID:      "US West (Oregon)",
	}

	// HostedZoneIDPerRegionNLBMap maps HostedZoneIDs from known regions.
	// See https://docs.aws.amazon.com/general/latest/gr/elb.html#elb_region
	HostedZoneIDPerRegionNLBMap = map[string]string{
		AfSouth1RegionID:     "Z203XCE67M25HM",
		ApEast1RegionID:      "Z12Y7K3UBGUAD1",
		ApNortheast1RegionID: "Z31USIVHYNEOWT",
		ApNortheast2RegionID: "ZIBE1TIR4HY56",
		ApNortheast3RegionID: "Z1GWIQ4HH19I5X",
		ApSouth1RegionID:     "ZVDDRBQ08TROA",
		ApSouth2RegionID:     "Z0711778386UTO08407HT",
		ApSoutheast1RegionID: "ZKVM4W9LS7TM",
		ApSoutheast2RegionID: "ZCT6FZBF4DROD",
		ApSoutheast3RegionID: "Z01971771FYVNCOVWJU1G",
		ApSoutheast4RegionID: "Z01156963G8MIIL7X90IV",
		CaCentral1RegionID:   "Z2EPGBW3API2WT",
		CnNorth1RegionID:     "Z3QFB96KMJ7ED6",
		CnNorthwest1RegionID: "ZQEIKTCZ8352D",
		EuCentral1RegionID:   "Z3F0SRJ5LGBH90",
		EuCentral2RegionID:   "Z02239872DOALSIDCX66S",
		EuNorth1RegionID:     "Z1UDT6IFJ4EJM",
		EuSouth1RegionID:     "Z23146JA1KNAFP",
		EuSouth2RegionID:     "Z1011216NVTVYADP1SSV",
		EuWest1RegionID:      "Z2IFOLAFXWLO4F",
		EuWest2RegionID:      "ZD4D7Y8KGAS4G",
		EuWest3RegionID:      "Z1CMS0P5QUZ6D5",
		MeCentral1RegionID:   "Z00282643NTTLPANJJG2P",
		MeSouth1RegionID:     "Z3QSRYVP46NYYV",
		SaEast1RegionID:      "ZTK26PT1VY4CU",
		UsEast1RegionID:      "Z26RNL4JYFTOTI",
		UsEast2RegionID:      "ZLMOA37VPKANP",
		UsGovEast1RegionID:   "Z1ZSMQQ6Q24QQ8",
		UsGovWest1RegionID:   "ZMG1MZ2THAWF1",
		UsWest1RegionID:      "Z24FKFUX50B4VW",
		UsWest2RegionID:      "Z18D5FSROUN65G",
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
	ApEast2RegionID      = "ap-east-2"      // Asia Pacific (Taipei).
	ApNortheast1RegionID = "ap-northeast-1" // Asia Pacific (Tokyo).
	ApNortheast2RegionID = "ap-northeast-2" // Asia Pacific (Seoul).
	ApNortheast3RegionID = "ap-northeast-3" // Asia Pacific (Osaka).
	ApSouth1RegionID     = "ap-south-1"     // Asia Pacific (Mumbai).
	ApSouth2RegionID     = "ap-south-2"     // Asia Pacific (Hyderabad).
	ApSoutheast1RegionID = "ap-southeast-1" // Asia Pacific (Singapore).
	ApSoutheast2RegionID = "ap-southeast-2" // Asia Pacific (Sydney).
	ApSoutheast3RegionID = "ap-southeast-3" // Asia Pacific (Jakarta).
	ApSoutheast4RegionID = "ap-southeast-4" // Asia Pacific (Melbourne).
	ApSoutheast5RegionID = "ap-southeast-5" // Asia Pacific (Malaysia).
	ApSoutheast6RegionID = "ap-southeast-6" // Asia Pacific (New Zealand).
	ApSoutheast7RegionID = "ap-southeast-7" // Asia Pacific (Thailand).
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
	MxCentral1RegionID   = "mx-central-1"   // Mexico (Central).
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
