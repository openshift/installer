package alicloud

//import "github.com/denverdino/aliyungo/ecs"

type GroupRuleNicType string

const (
	GroupRuleInternet = GroupRuleNicType("internet")
	GroupRuleIntranet = GroupRuleNicType("intranet")
)

type Direction string

const (
	DirectionIngress = Direction("ingress")
	DirectionEgress  = Direction("egress")
)

type GroupRulePolicy string

const (
	GroupRulePolicyAccept = GroupRulePolicy("accept")
	GroupRulePolicyDrop   = GroupRulePolicy("drop")
)

type GroupInnerAccessPolicy string

const (
	GroupInnerAccept = GroupInnerAccessPolicy("Accept")
	GroupInnerDrop   = GroupInnerAccessPolicy("Drop")
)

type SpotStrategyType string

// Constants of SpotStrategyType
const (
	NoSpot             = SpotStrategyType("NoSpot")
	SpotWithPriceLimit = SpotStrategyType("SpotWithPriceLimit")
	SpotAsPriceGo      = SpotStrategyType("SpotAsPriceGo")
)

type DestinationResource string

const (
	ZoneResource         = DestinationResource("Zone")
	IoOptimizedResource  = DestinationResource("IoOptimized")
	InstanceTypeResource = DestinationResource("InstanceType")
	SystemDiskResource   = DestinationResource("SystemDisk")
	DataDiskResource     = DestinationResource("DataDisk")
	NetworkResource      = DestinationResource("Network")
)

const GenerationOne = "ecs-1"
const GenerationTwo = "ecs-2"
const GenerationThree = "ecs-3"
const GenerationFour = "ecs-4"

var NoneIoOptimizedFamily = map[string]string{"ecs.t1": "", "ecs.t2": "", "ecs.s1": ""}
var NoneIoOptimizedInstanceType = map[string]string{"ecs.s2.small": ""}
var HalfIoOptimizedFamily = map[string]string{"ecs.s2": "", "ecs.s3": "", "ecs.m1": "", "ecs.m2": "", "ecs.c1": "", "ecs.c2": ""}

var OutdatedDiskCategory = map[DiskCategory]DiskCategory{
	DiskCloud: DiskCloud}

var SupportedDiskCategory = map[DiskCategory]DiskCategory{
	DiskCloudSSD:        DiskCloudSSD,
	DiskCloudEfficiency: DiskCloudEfficiency,
	DiskEphemeralSSD:    DiskEphemeralSSD,
	DiskCloudESSD:       DiskCloudESSD,
	DiskCloud:           DiskCloud,
}

const AllPortRange = "-1/-1"

const (
	KubernetesMasterNumber = 3
)

type RenewalStatus string

const (
	RenewAutoRenewal = RenewalStatus("AutoRenewal")
	RenewNormal      = RenewalStatus("Normal")
	RenewNotRenewal  = RenewalStatus("NotRenewal")
)

type DiskType string

const (
	DiskTypeAll    = DiskType("all")
	DiskTypeSystem = DiskType("system")
	DiskTypeData   = DiskType("data")
)

type DiskCategory string

const (
	DiskAll             = DiskCategory("all") //Default
	DiskCloud           = DiskCategory("cloud")
	DiskEphemeralSSD    = DiskCategory("ephemeral_ssd")
	DiskCloudESSD       = DiskCategory("cloud_essd")
	DiskCloudEfficiency = DiskCategory("cloud_efficiency")
	DiskCloudSSD        = DiskCategory("cloud_ssd")
	DiskLocalDisk       = DiskCategory("local_disk")
)

type DiskResizeType string

const (
	DiskResizeTypeOffline = DiskResizeType("offline")
	DiskResizeTypeOnline  = DiskResizeType("online")
)

type ImageOwnerAlias string

const (
	ImageOwnerSystem      = ImageOwnerAlias("system")
	ImageOwnerSelf        = ImageOwnerAlias("self")
	ImageOwnerOthers      = ImageOwnerAlias("others")
	ImageOwnerMarketplace = ImageOwnerAlias("marketplace")
	ImageOwnerDefault     = ImageOwnerAlias("") //Return the values for system, self, and others
)

type SecurityEnhancementStrategy string

const (
	ActiveSecurityEnhancementStrategy   = SecurityEnhancementStrategy("Active")
	DeactiveSecurityEnhancementStrategy = SecurityEnhancementStrategy("Deactive")
)

type CreditSpecification string

const (
	CreditSpecificationStandard  = CreditSpecification("Standard")
	CreditSpecificationUnlimited = CreditSpecification("Unlimited")
)
